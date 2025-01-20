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
Package server provides a collection of server implementations for the Cosmic Connector
application. It includes three main server types: GrpcServer, ChannelzServer, and
PprofServer, all implementing the common Server interface.

Server Types:

  - GrpcServer: Handles the main gRPC communication for the Federation service
  - ChannelzServer: Provides gRPC channelz monitoring capabilities
  - PprofServer: Exposes Go runtime profiling data via HTTP endpoints

Each server implementation provides consistent Start and Shutdown methods for
lifecycle management. The servers can be used independently or together as part
of a larger application.

Example usage:

	import (
	 "context"
	 "github.com/rs/zerolog"
	 "aalyria.com/spacetime/apps/interconnectprovider/server"
	 "aalyria.com/spacetime/apps/interconnectprovider/handler"
	)

	func main() {
	 logger := zerolog.New(os.Stdout)

	 // Create a new gRPC server
	 federationHandler := handler.NewFederationHandler()
	 grpcServer := server.NewGrpcServer(8080, federationHandler, logger)

	 // Start the server
	 ctx := context.Background()
	 if err := grpcServer.Start(ctx); err != nil {
	  logger.Fatal().Err(err).Msg("Failed to start gRPC server")
	 }

	 // Graceful shutdown
	 if err := grpcServer.Shutdown(ctx); err != nil {
	  logger.Error().Err(err).Msg("Error during server shutdown")
	 }
	}

The Server interface ensures consistent behavior across all server implementations:

	type Server interface {
	 Start(ctx context.Context) error
	 Shutdown(ctx context.Context) error
	}

All server implementations use structured logging via zerolog and support graceful
shutdown operations. They are designed to be concurrent-safe and suitable for
production deployments.
*/
package server
