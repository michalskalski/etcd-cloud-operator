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

// Package main implements basic logic to start the etcd-cloud-operator.
package main

import (
	"flag"
	"os"

	"go.uber.org/zap"

	"github.com/quentin-m/etcd-cloud-operator/pkg/logger"
	"github.com/quentin-m/etcd-cloud-operator/pkg/operator"

	// Register providers.
	_ "github.com/quentin-m/etcd-cloud-operator/pkg/providers/asg/aws"
	_ "github.com/quentin-m/etcd-cloud-operator/pkg/providers/asg/docker"
	_ "github.com/quentin-m/etcd-cloud-operator/pkg/providers/asg/sts"
	_ "github.com/quentin-m/etcd-cloud-operator/pkg/providers/snapshot/file"
	_ "github.com/quentin-m/etcd-cloud-operator/pkg/providers/snapshot/s3"
)

func main() {
	// Parse command-line arguments.
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagConfigPath := flag.String("config", "", "Load configuration from the specified file.")
	flagLogLevel := flag.String("log-level", "info", "Define the logging level.")
	flagWeServerPort := flag.Int("web-server-port", 2378, "Define listen port for status service.")
	flag.Parse()

	// Initialize logging system.
	logger.Configure(*flagLogLevel)

	// Read configuration.
	config, err := loadConfig(*flagConfigPath)
	if err != nil {
		zap.S().With(zap.Error(err)).Fatal("failed to load configuration")
	}

	if config.ECO.WebServerPort == nil {
		port := int32(*flagWeServerPort)
		config.ECO.WebServerPort = &port
	}

	if config.ECO.Etcd.AdvertiseAddress == "" {
		config.ECO.Etcd.AdvertiseAddress = os.Getenv("POD_IP")
		logger.Logger.Warn("Set Advertise address to", zap.String("ip", config.ECO.Etcd.AdvertiseAddress))
	}

	// Run.
	operator.New(config.ECO).Run(*flagWeServerPort)
}
