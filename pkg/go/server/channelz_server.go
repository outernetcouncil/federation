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
	"sync/atomic"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/channelz/service"
)

// ChannelzServer implements the Server interface for Channelz service.
type ChannelzServer struct {
	address string
	logger  zerolog.Logger
	srv     *grpc.Server
	lis     net.Listener
	running atomic.Bool
}

// NewChannelzServer creates a new ChannelzServer with the given address.
func NewChannelzServer(address string, logger zerolog.Logger) *ChannelzServer {
	return &ChannelzServer{
		address: address,
		logger:  logger,
		srv:     grpc.NewServer(),
	}
}

// Start begins serving the Channelz service and blocks until the server is stopped.
func (c *ChannelzServer) Start(ctx context.Context) error {
	if !c.running.CompareAndSwap(false, true) {
		return fmt.Errorf("server already running")
	}

	lis, err := net.Listen("tcp", c.address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", c.address, err)
	}
	c.lis = lis
	service.RegisterChannelzServiceToServer(c.srv)
	c.logger.Info().Msgf("Starting Channelz server on %s", c.address)

	// Directly call Serve to block until the server is stopped
	if err := c.srv.Serve(lis); err != nil && err != grpc.ErrServerStopped {
		c.logger.Error().Err(err).Msg("Channelz server encountered an error")
		return err
	}

	return nil
}

// Shutdown gracefully stops the Channelz server.
func (c *ChannelzServer) Shutdown(ctx context.Context) error {
	if !c.running.CompareAndSwap(true, false) {
		return nil
	}

	c.logger.Info().Msg("Shutting down Channelz server")
	c.srv.GracefulStop()

	return nil
}
