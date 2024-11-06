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
package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"

	"github.com/rs/zerolog"
)

// PprofServer implements the Server interface for pprof HTTP endpoints.
type PprofServer struct {
	address string
	srv     *http.Server
	logger  zerolog.Logger
	lis     net.Listener
}

// NewPprofServer creates a new PprofServer with the given address and logger.
func NewPprofServer(address string, logger zerolog.Logger) *PprofServer {
	return &PprofServer{
		address: address,
		logger:  logger,
	}
}

// Start begins serving pprof endpoints and blocks until the server is stopped.
func (p *PprofServer) Start(ctx context.Context) error {
	if p.srv != nil {
		return fmt.Errorf("server already running")
	}

	if p.address == "" {
		p.logger.Warn().Msg("Pprof address not configured, skipping pprof server")
		return nil
	}

	// Create listener
	lis, err := net.Listen("tcp", p.address)
	if err != nil {
		p.logger.Error().Err(err).Msgf("Failed to listen on pprof address %s", p.address)
		return fmt.Errorf("failed to listen on pprof address %s: %w", p.address, err)
	}
	p.lis = lis

	// Create HTTP server
	p.srv = &http.Server{
		Handler: http.DefaultServeMux,
	}

	p.logger.Info().Msgf("Starting pprof server on %s", p.address)

	// Directly call Serve to block until the server is stopped
	if err := p.srv.Serve(lis); err != nil && err != http.ErrServerClosed {
		p.logger.Error().Err(err).Msg("pprof server encountered an error")
		return err
	}

	return nil
}

// Shutdown gracefully shuts down the pprof server.
func (p *PprofServer) Shutdown(ctx context.Context) error {
	if p.srv == nil {
		p.logger.Warn().Msg("pprof server is not running")
		return nil
	}

	p.logger.Info().Msg("Shutting down pprof server")
	if err := p.srv.Shutdown(ctx); err != nil {
		p.logger.Error().Err(err).Msg("Failed to shutdown pprof server gracefully")
		return fmt.Errorf("failed to shutdown pprof server: %w", err)
	}
	p.logger.Info().Msg("pprof server shutdown complete")
	return nil
}
