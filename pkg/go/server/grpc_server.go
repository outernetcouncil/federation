// Copyright 2024 Outernet Council Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package server

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/outernetcouncil/federation/gen/go/federation/v1alpha"
	"github.com/outernetcouncil/federation/pkg/go/handler"
)

// GrpcServer implements the Server interface for gRPC services.
type GrpcServer struct {
	pb.UnimplementedFederationServiceServer
	port    int
	handler handler.FederationHandler
	logger  zerolog.Logger
	srv     *grpc.Server
	lis     net.Listener
}

// NewGrpcServer creates a new GrpcServer with the given port, handler, and logger.
func NewGrpcServer(port int, handler handler.FederationHandler, logger zerolog.Logger) *GrpcServer {
	return &GrpcServer{
		port:    port,
		handler: handler,
		logger:  logger,
	}
}

// Start begins serving gRPC requests and blocks until the server is stopped.
func (g *GrpcServer) Start(ctx context.Context) error {
	if g.srv != nil {
		return fmt.Errorf("server already running")
	}

	addr := fmt.Sprintf(":%d", g.port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		g.logger.Error().Err(err).Msgf("Failed to listen on port %d", g.port)
		return fmt.Errorf("failed to listen on port %d: %w", g.port, err)
	}
	g.lis = lis

	g.srv = grpc.NewServer()
	pb.RegisterFederationServiceServer(g.srv, g.handler)
	reflection.Register(g.srv)

	g.logger.Info().Msgf("Starting gRPC server on port %d", g.port)

	// Directly call Serve to block until the server is stopped
	if err := g.srv.Serve(lis); err != nil && err != grpc.ErrServerStopped {
		g.logger.Error().Err(err).Msg("gRPC server encountered an error")
		return err
	}

	return nil
}

// Shutdown gracefully stops the gRPC server.
func (g *GrpcServer) Shutdown(ctx context.Context) error {
	if g.srv == nil {
		g.logger.Warn().Msg("gRPC server is not running")
		return nil
	}

	g.logger.Info().Msg("Shutting down gRPC server")
	g.srv.GracefulStop()
	return nil
}
