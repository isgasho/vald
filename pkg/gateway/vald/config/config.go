//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package setting stores all server application settings
package config

import (
	"github.com/vdaas/vald/internal/config"
)

// Config represent a application setting data content (config.yaml).
// In K8s environment, this configuration is stored in K8s ConfigMap.
type Data struct {
	// Version represent configuration file version.
	Version string `json:"version" yaml:"version"`

	// Server represent all server configurations
	Server *config.Servers `json:"server_config" yaml:"server_config"`

	// Discoverer represent agent discoverer service configuration
	Discoverer *config.Discoverer `json:"discoverer" yaml:"discoverer"`

	// MetaProxy represent metadata gateway configuration
	Meta *config.MetaProxy `json:"meta" yaml:"meta"`

	// IndexManager represent index manager configuration
	IndexManager *config.IndexManager `json:"meta" yaml:"meta"`
}

func NewConfig(path string) (cfg *Data, err error) {
	err = config.Read(path, &cfg)

	if err != nil {
		return nil, err
	}

	if cfg.Server != nil {
		cfg.Server = cfg.Server.Bind()
	}
	if cfg.Discoverer.Bind() != nil {
		cfg.Discoverer = cfg.Discoverer.Bind()
	}
	if cfg.Meta.Bind() != nil {
		cfg.Meta = cfg.Meta.Bind()
	}
	if cfg.IndexManager.Bind() != nil {
		cfg.IndexManager = cfg.IndexManager.Bind()
	}

	return cfg, nil
}

// func FakeData() {
// 	d := Data{
// 		Version: "v0.0.1",
// 		Server: &config.Servers{
// 			Servers: []*config.Server{
// 				{
// 					Name:              "agent-rest",
// 					Host:              "127.0.0.1",
// 					Port:              8080,
// 					Mode:              "REST",
// 					ProbeWaitTime:     "3s",
// 					ShutdownDuration:  "5s",
// 					HandlerTimeout:    "5s",
// 					IdleTimeout:       "2s",
// 					ReadHeaderTimeout: "1s",
// 					ReadTimeout:       "1s",
// 					WriteTimeout:      "1s",
// 				},
// 				{
// 					Name: "agent-grpc",
// 					Host: "127.0.0.1",
// 					Port: 8082,
// 					Mode: "GRPC",
// 				},
// 			},
// 			MetricsServers: []*config.Server{
// 				{
// 					Name:              "pprof",
// 					Host:              "127.0.0.1",
// 					Port:              6060,
// 					Mode:              "REST",
// 					ProbeWaitTime:     "3s",
// 					ShutdownDuration:  "5s",
// 					HandlerTimeout:    "5s",
// 					IdleTimeout:       "2s",
// 					ReadHeaderTimeout: "1s",
// 					ReadTimeout:       "1s",
// 					WriteTimeout:      "1s",
// 				},
// 			},
// 			HealthCheckServers: []*config.Server{
// 				{
// 					Name: "livenesss",
// 					Host: "127.0.0.1",
// 					Port: 3000,
// 				},
// 				{
// 					Name: "readiness",
// 					Host: "127.0.0.1",
// 					Port: 3001,
// 				},
// 			},
// 			StartUpStrategy: []string{
// 				"livenesss",
// 				"pprof",
// 				"agent-grpc",
// 				"agent-rest",
// 				"readiness",
// 			},
// 			ShutdownStrategy: []string{
// 				"readiness",
// 				"agent-rest",
// 				"agent-grpc",
// 				"pprof",
// 				"livenesss",
// 			},
// 			FullShutdownDuration: "30s",
// 			TLS: &config.TLS{
// 				Enabled: false,
// 				Cert:    "/path/to/cert",
// 				Key:     "/path/to/key",
// 				CA:      "/path/to/ca",
// 			},
// 		},
// 		NGT: &config.NGT{
// 			IndexPath:           "/path/to/index",
// 			Dimension:           4096,
// 			BulkInsertChunkSize: 10,
// 			DistanceType:        "l2",
// 			ObjectType:          "float",
// 			CreationEdgeSize:    20,
// 			SearchEdgeSize:      10,
// 		},
// 	}
// 	fmt.Println(config.ToRawYaml(d))
// }