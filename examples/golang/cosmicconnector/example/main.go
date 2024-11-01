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

// package main produces an example Cosmic Connector binary, offering a simple and incomplete example of a Federation service.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/outernetcouncil/federation/pkg/go/cosmicconnector"
	"github.com/outernetcouncil/federation/examples/golang/cosmicconnector/config"
	examplehandler "github.com/outernetcouncil/federation/examples/golang/cosmicconnector/example/handler"
	"github.com/outernetcouncil/federation/pkg/go/server"
)

const (
	appName = "cosmic_connector"
)

func baseContext(ctx context.Context) context.Context {
	var log zerolog.Logger
	if os.Getenv("TERM") != "" {
		log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02 03:04:05PM"})
	} else {
		log = zerolog.New(os.Stdout)
	}
	return log.With().Timestamp().Logger().WithContext(ctx)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = baseContext(ctx)
	logger := zerolog.Ctx(ctx)

	fs := flag.NewFlagSet(appName, flag.ContinueOnError)
	confPath := fs.String("config", "", "The path to a text protobuf representation of the connector's configuration (a ConnectorParams message).")
	dryRunOnly := fs.Bool("dry-run", false, "Just validate the config, don't start the agent. Exits with a non-zero return code if the config is invalid.")
	logLevel := config.LogLevelFlag(zerolog.InfoLevel)
	fs.Var(&logLevel, "log-level", "The log level (one of disabled, warn, panic, info, fatal, error, debug, or trace) to use.")
	fs.Usage = func() {
		w := fs.Output()
		fmt.Fprintf(w, "Usage: %s [options] [query|exec|dump]\n", appName)
		fmt.Fprint(w, "\nOptions:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(os.Args[1:]); err == flag.ErrHelp {
		fs.Usage()
		os.Exit(0)
	} else if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse flags")
	}

	cp, err := config.ReadParams(*confPath)
	if err != nil {
		logger.Fatal().Err(err).Msgf("failed config.ReadParams(%s)", *confPath)
	}

	if *dryRunOnly {
		logger.Info().Msg("Dry run only, terminating")
		return
	}

	// Initialize Servers based on configuration
	grpcServer := server.NewGrpcServer(int(cp.GetPort()), examplehandler.NewPrototypeHandler(), *logger)
	pprofServer := server.NewPprofServer(cp.GetObservabilityParams().GetPprofAddress(), *logger)
	channelzServer := server.NewChannelzServer(cp.GetObservabilityParams().GetChannelzAddress(), *logger)

	// Create CosmicConnector with initialized servers
	connector := cosmicconnector.NewCosmicConnector(*logger, grpcServer, pprofServer, channelzServer)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	errChan := make(chan error, 1)
	go func() {
		errChan <- connector.Run(ctx)
	}()

	// Wait for either an error from the connector or an interrupt signal
	select {
	case err := <-errChan:
		if err != nil {
			logger.Fatal().Err(err).Msg("cosmicconnector.Run(...) failed")
		}
	case <-sigChan:
		logger.Info().Msg("Received interrupt signal. Shutting down...")
		cancel()

		// Give some time for graceful shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		select {
		case <-shutdownCtx.Done():
			logger.Warn().Msg("Shutdown timed out")
		case err := <-errChan:
			if err != nil {
				logger.Error().Err(err).Msg("Error during shutdown")
			} else {
				logger.Info().Msg("Shutdown completed successfully")
			}
		}
	}
}
