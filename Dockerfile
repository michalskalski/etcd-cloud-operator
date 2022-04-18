FROM golang:1.18-bullseye AS build-env

WORKDIR /go/src/github.com/quentin-m/etcd-cloud-operator

# Install & cache dependencies
RUN apt-get install -y git curl

RUN curl -L https://github.com/etcd-io/etcd/releases/download/v3.5.3/etcd-v3.5.3-linux-amd64.tar.gz -o /tmp/etcd.tar.gz && \
    mkdir /etcd && \
    tar xzvf /tmp/etcd.tar.gz -C /etcd --strip-components=1 && \
    rm /tmp/etcd.tar.gz

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM build-env as builder
COPY . .
RUN cd cmd/operator && go build -o  /go/bin/operator
RUN cd cmd/tester && go build  -o  /go/bin/tester

# Copy ECO and etcdctl into an ubi container image.
FROM registry.access.redhat.com/ubi8/ubi-minimal:8.5-240.1648458092

COPY --from=builder /go/bin/operator /operator
COPY --from=builder /go/bin/tester /tester
COPY --from=builder /etcd/etcdctl /usr/local/bin/etcdctl

ENTRYPOINT ["/operator"]
CMD ["-config", "/etc/eco/eco.yaml"]
