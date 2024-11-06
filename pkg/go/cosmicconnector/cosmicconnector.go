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

// Package CosmicConector provides a Cosmic Connector, a Federation gRPC bridge to
// RESTFUL Providers.
package cosmicconnector

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/outernetcouncil/federation/pkg/go/server"
)

type CosmicConnector struct {
	servers []server.Server
	logger  zerolog.Logger
}

func NewCosmicConnector(logger zerolog.Logger, servers ...server.Server) *CosmicConnector {
	return &CosmicConnector{servers: servers, logger: logger}
}

func (o *CosmicConnector) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	// Add shutdown goroutine
	g.Go(func() error {
		<-ctx.Done()
		return o.shutdown(context.Background())
	})

	// Start servers
	for _, srv := range o.servers {
		s := srv
		g.Go(func() error {
			o.logger.Info().Msgf("Starting server: %T", s)
			return s.Start(ctx)
		})
	}

	return g.Wait()
}

func (o *CosmicConnector) shutdown(ctx context.Context) error {
	var errs []error
	for _, srv := range o.servers {
		if err := srv.Shutdown(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}
	return nil
}
