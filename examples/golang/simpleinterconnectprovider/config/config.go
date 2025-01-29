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

// Package config provides utilities related to configuration reading
// and enactment.
package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"google.golang.org/protobuf/encoding/prototext"

	configpb "github.com/outernetcouncil/federation/gen/go/examples/golang/simpleinterconnectprovider/config"
)

func ReadParams(confPath string) (*configpb.ConnectorParams, error) {
	if confPath == "" {
		return nil, errors.New("no config (--config) provided")
	}

	confData, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	} else if len(confData) == 0 {
		return nil, errors.New("empty config (--config) provided")
	}

	conf := &configpb.ConnectorParams{}
	if err = prototext.Unmarshal(confData, conf); err != nil {
		return nil, fmt.Errorf("unmarshalling config proto: %w", err)
	}

	return conf, err
}

// logLevelFlag is a flag.Value implementation for the zerolog.Level type.
type LogLevelFlag zerolog.Level

func (l *LogLevelFlag) String() string {
	return fmt.Sprintf("%q", zerolog.Level(*l).String())
}

func (l *LogLevelFlag) Set(value string) error {
	level, err := zerolog.ParseLevel(value)
	if err != nil {
		return err
	}

	*l = LogLevelFlag(level)
	return nil
}
