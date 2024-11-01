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

The package supports various configuration aspects including:

  - Connection parameters for gRPC services
  - Transport security settings (TLS and insecure options)
  - Authentication strategies (JWT-based authentication)
  - Observability settings (channelz and pprof)

Key functionalities:

  - ReadParams: Reads and unmarshals configuration from a file path
  - GetDialOpts: Generates gRPC dial options based on connection parameters
  - Custom logging level management through LogLevelFlag

Configuration Structure:

The configuration is defined in a Protocol Buffer format (config.proto) with the following
main components:

  - ConnectorParams: Top-level configuration containing port and other settings
  - ConnectionParams: Network and security-related settings
  - ObservabilityParams: Debugging and monitoring configurations
  - AuthStrategy: Authentication configuration, supporting JWT

Example Usage:

    params, err := config.ReadParams(configPath)
    if err != nil {
        log.Fatal(err)
    }

    dialOpts, err := config.GetDialOpts(ctx, params.ConnectionParams, clockwork.NewRealClock())
    if err != nil {
        log.Fatal(err)
    }

The package is designed to be used as part of the Cosmic Connector service but can be
adapted for other gRPC-based services requiring similar configuration management.
*/
package config
