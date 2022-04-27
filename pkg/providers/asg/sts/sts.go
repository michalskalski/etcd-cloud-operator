// Copyright 2017 Quentin Machu & eco authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sts

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/quentin-m/etcd-cloud-operator/pkg/providers/asg"
)

func init() {
	asg.Register("sts", &sts{})
}

type sts struct {
	namespace, name, serviceName, dnsClusterSuffix string
	replicas                                       int

	self instance
}

type instance struct {
	id, name, address, bindAddress string
}

func (i *instance) Name() string {
	return i.name
}

func (i *instance) Address() string {
	return i.address
}

func (i *instance) BindAddress() string {
	return "0.0.0.0"
}

func (a *sts) Configure(providerConfig asg.Config) (err error) {
	a.namespace, err = envOrErr("STATEFULSET_NAMESPACE")
	if err != nil {
		return
	}

	a.name, err = envOrErr("STATEFULSET_NAME")
	if err != nil {
		return
	}

	a.self.address, err = envOrErr("POD_IP")
	if err != nil {
		return
	}

	a.self.name, err = envOrErr("POD_NAME")
	if err != nil {
		return
	}

	return
}

func (a *sts) AutoScalingGroupStatus() ([]asg.Instance, asg.Instance, int, error) {
	var instances []asg.Instance
	var instancesStr []string

	client := k8Client()
	pods, err := GetPodsForApp(client, context.Background(), a.namespace, a.name)
	if err != nil {
		return nil, &a.self, a.replicas, err
	}

	for _, pod := range pods {
		instance := instance{
			name:    pod.Name,
			address: pod.Status.PodIP,
		}
		instances = append(instances, &instance)
		instancesStr = append(instancesStr, instance.address)
	}
	replicas, err := GetSTSReplicas(client, context.Background(), a.namespace, a.name)
	if err != nil {
		return nil, &a.self, a.replicas, err
	}
	zap.S().Debugf("Discovered %d / %d replicas: %s", len(instances), replicas, strings.Join(instancesStr, ", "))
	return instances, &a.self, replicas, nil
}

func envOrErr(key string) (value string, err error) {
	value = os.Getenv(key)
	if value == "" {
		err = fmt.Errorf("Required environment variable: %s was not set", key)
	}

	return
}
