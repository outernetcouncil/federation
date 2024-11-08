// Copyright 2024 Outernet Council Foundation
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

/*
Package config provides configuration management utilities for the Cosmic Connector service.
It handles reading, parsing, and validating configuration parameters from protobuf-based
configuration files.

# Configuration Components

The package supports the following configuration aspects:

  - Server configuration (port, TLS settings)
  - Observability parameters (metrics, tracing, debugging)
  - Logging configuration

# Core Types

  - ConnectorParams: Primary configuration container
  - ObservabilityParams: Monitoring and debugging settings
  - LogLevelFlag: Custom flag.Value implementation for zerolog levels

# Usage Examples

Basic configuration reading:

	params, err := config.ReadParams(configPath)
	if err != nil {
	    log.Fatal(err)
	}

Configuration validation:

	if err := params.Validate(); err != nil {
	    log.Fatal("Invalid configuration:", err)
	}

# Configuration Format

Configuration files use Protocol Buffer text format. Example:

	port: 50052
	observability_params {
	    channelz_address: "0.0.0.0:50051"
	    pprof_address: "0.0.0.0:6060"
	}

# Integration Points

The package integrates with:

  - zerolog for logging
  - protobuf for configuration format
  - flag package for command-line parsing

# Monitoring and Debugging

The package supports:

  - Channelz for gRPC monitoring
  - pprof for performance profiling
  - Custom metric endpoints

For detailed technical specifications, refer to the config.proto file.
*/
package config
