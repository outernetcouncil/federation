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
package handler

import (
	"context"
	"testing"

	pb "github.com/outernetcouncil/federation/gen/go/federation/interconnect/v1alpha"
	"outernetcouncil.org/nmts/v1alpha/proto/ek/physical"
)

func TestPrototypeHandler_CreateTransceiver(t *testing.T) {
	h := NewPrototypeHandler()
	ctx := context.Background()

	defaultTransceiver := &pb.Transceiver{
		TransmitSignalChain: &pb.TransmitSignalChain{
			Antenna: &physical.Antenna{
				Type: physical.Antenna_OPTICAL,
			},
		},
		ReceiveSignalChain: &pb.ReceiveSignalChain{
			Antenna: &physical.Antenna{
				Type: physical.Antenna_OPTICAL,
			},
		},
	}

	_, err := h.CreateTransceiver(ctx, &pb.CreateTransceiverRequest{
		TransceiverId: "existing",
		Transceiver:   defaultTransceiver,
	})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name          string
		transceiverID string
		transceiver   *pb.Transceiver
		wantError     bool
	}{
		{
			name:          "Creates a new transceiver",
			transceiverID: "unique",
			transceiver:   defaultTransceiver,
			wantError:     false,
		},
		{
			name:          "Fails creating a transceiver with already existing ID",
			transceiverID: "existing",
			transceiver:   defaultTransceiver,
			wantError:     true,
		},
		{
			name:          "Fails creating a transceiver with wrong antenna type in transmit signal chain",
			transceiverID: "wrong_antenna",
			transceiver: &pb.Transceiver{
				TransmitSignalChain: &pb.TransmitSignalChain{
					Antenna: &physical.Antenna{
						Type: physical.Antenna_RF,
					},
				},
				ReceiveSignalChain: &pb.ReceiveSignalChain{
					Antenna: &physical.Antenna{
						Type: physical.Antenna_OPTICAL,
					},
				},
			},
			wantError: true,
		},
		{
			name:          "Fails creating a transceiver with wrong antenna type in transmit signal chain",
			transceiverID: "wrong_antenna",
			transceiver: &pb.Transceiver{
				TransmitSignalChain: &pb.TransmitSignalChain{
					Antenna: &physical.Antenna{
						Type: physical.Antenna_OPTICAL,
					},
				},
				ReceiveSignalChain: &pb.ReceiveSignalChain{
					Antenna: &physical.Antenna{
						Type: physical.Antenna_RF,
					},
				},
			},
			wantError: true,
		},
		{
			name:          "Fails creating a transceiver with unspecified antenna type",
			transceiverID: "wrong_antenna",
			transceiver:   &pb.Transceiver{},
			wantError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.CreateTransceiverRequest{
				TransceiverId: tt.transceiverID,
				Transceiver:   tt.transceiver,
			}

			_, err := h.CreateTransceiver(ctx, req)
			if !tt.wantError && err != nil {
				t.Fatalf("Creating a transceiver failed %s.", err)
			}
			if tt.wantError && err == nil {
				t.Fatalf("Creating a transceiver should have failed, but did not.")
			}
		})
	}
}

// TODO: Add test to get transceiver filters + transceiver lifecycle test, i.e. Get, Update, Delete tests.
