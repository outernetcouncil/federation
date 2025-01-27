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
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	pb "github.com/outernetcouncil/federation/gen/go/federation/interconnect/v1alpha"
	"google.golang.org/genproto/googleapis/type/interval"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
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
			name:          "Fails creating a transceiver with wrong antenna type in receive signal chain",
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
		{
			name:          "Returns correct name, even if initial name is not set correctly",
			transceiverID: "correct",
			transceiver: &pb.Transceiver{
				Name: "wrong_name",
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
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.CreateTransceiverRequest{
				TransceiverId: tt.transceiverID,
				Transceiver:   tt.transceiver,
			}

			resp, err := h.CreateTransceiver(ctx, req)
			if !tt.wantError && err != nil {
				t.Fatalf("Creating a transceiver failed %s.", err)
			}
			if tt.wantError && err == nil {
				t.Fatalf("Creating a transceiver should have failed, but did not.")
			}
			if !tt.wantError && resp.Name != fmt.Sprintf("transceivers/%s", tt.transceiverID) {
				t.Fatalf("Name is incorrect")
			}
		})
	}
}

func TestPrototypeHandler_UpdateTransceiver(t *testing.T) {
	h, ctx := createExistingTransceiver(t)

	tests := []struct {
		name        string
		transceiver *pb.Transceiver
		wantError   bool
	}{
		{
			name: "Fails updating a transceiver with none-existing ID",
			transceiver: &pb.Transceiver{
				Name: "transceiver/non-existant",
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
			name: "Fails creating a transceiver with wrong antenna type in transmit signal chain",
			transceiver: &pb.Transceiver{
				Name: "transceiver/existing",
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
			name: "Fails creating a transceiver with wrong antenna type in receive signal chain",
			transceiver: &pb.Transceiver{
				Name: "transceiver/existing",
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
			name:        "Fails creating a transceiver with unspecified antenna type",
			transceiver: &pb.Transceiver{},
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.UpdateTransceiverRequest{
				Transceiver: tt.transceiver,
			}

			_, err := h.UpdateTransceiver(ctx, req)
			if !tt.wantError && err != nil {
				t.Fatalf("Updating a transceiver failed %s.", err)
			}
			if tt.wantError && err == nil {
				t.Fatalf("Updating a transceiver should have failed, but did not.")
			}
		})
	}
}

func TestPrototypeHandler_TransceiverCannotBeChangedIfBearerIsAttached(t *testing.T) {
	h, ctx := createExistingTransceiver(t)

	t.Run("gets information about transceiver", func(t *testing.T) {
		resp, err := h.GetTransceiver(ctx, &pb.GetTransceiverRequest{
			Name: "transceivers/existing",
		})

		if err != nil {
			t.Fatalf("no error expected, but was %v", err)
		}
		if diff := cmp.Diff(&pb.Transceiver{
			Name: "transceivers/existing",
			ReceiveSignalChain: &pb.ReceiveSignalChain{
				Antenna: &physical.Antenna{
					Type: physical.Antenna_OPTICAL,
				},
			},
			TransmitSignalChain: &pb.TransmitSignalChain{
				Antenna: &physical.Antenna{
					Type: physical.Antenna_OPTICAL,
				},
			},
		}, resp,
			protocmp.Transform()); diff != "" {
			t.Errorf("Bearer mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("returns error when getting information about transceiver does not exist", func(t *testing.T) {
		_, err := h.GetTransceiver(ctx, &pb.GetTransceiverRequest{
			Name: "transceivers/non-existing",
		})

		if err == nil {
			t.Fatal("error expected, but was none")
		}
		if !strings.HasPrefix(err.Error(), "rpc error: code = NotFound") {
			t.Fatalf("expected error did not match error %v", err)
		}
	})

	t.Run("can delete existing transceiver", func(t *testing.T) {
		_, err := h.DeleteTransceiver(ctx, &pb.DeleteTransceiverRequest{
			Name: "transceivers/existing",
		})

		if err != nil {
			t.Fatalf("no error expected, but was %v", err)
		}
	})

	t.Run("throws error if transceiver does not exist", func(t *testing.T) {
		_, err := h.DeleteTransceiver(ctx, &pb.DeleteTransceiverRequest{
			Name: "transceivers/existing",
		})

		if err == nil {
			t.Fatalf("error expected")
		}
	})

	h, ctx = createExistingTransceiver(t)

	_, err := h.CreateBearer(ctx, &pb.CreateBearerRequest{
		BearerId: "newBearer",
		Bearer: &pb.Bearer{
			Target:      TARGET_NAME,
			Transceiver: "transceivers/existing",
			Interval: &interval.Interval{
				StartTime: &timestamppb.Timestamp{
					Seconds: int64(time.Now().Unix()),
				},
				EndTime: &timestamppb.Timestamp{
					Seconds: int64(time.Now().Unix()) + 60*60, // let's just have a one hour
				},
			},
			RxCenterFrequencyHz: 16000000000,
			RxBandwidthHz:       30000000,
			TxCenterFrequencyHz: 16000000000,
			TxBandwidthHz:       30000000,
		},
	})
	if err != nil {
		t.Fatalf("Test setup failed: %v", err)
	}

	t.Run("cannot update transceiver if bearer is attached", func(t *testing.T) {
		_, err := h.UpdateTransceiver(ctx, &pb.UpdateTransceiverRequest{
			Transceiver: &pb.Transceiver{
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
			},
		})

		if err == nil {
			t.Fatal("should have failed updating")
		}
	})

	t.Run("cannot delete transceiver if bearer is attached", func(t *testing.T) {
		_, err := h.DeleteTransceiver(ctx, &pb.DeleteTransceiverRequest{
			Name: "transceivers/existing",
		})

		if err == nil {
			t.Fatal("should have failed updating")
		}
		if err.Error() != "rpc error: code = FailedPrecondition desc = transceiver has bearer attached and cannot be deleted" {
			t.Fatalf("Wrong error message, was %s", err.Error())
		}
	})
}

func createExistingTransceiver(t *testing.T) (*PrototypeHandler, context.Context) {
	t.Helper()

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

	return h, ctx
}

func TestPrototypeHandler_CreateBearer(t *testing.T) {
	h, ctx := createExistingTransceiver(t)

	_, err := h.CreateBearer(ctx, &pb.CreateBearerRequest{
		BearerId: "existing",
		Bearer: &pb.Bearer{
			Name:                "bearers/existing",
			Target:              TARGET_NAME,
			Transceiver:         "transceivers/existing",
			Interval:            createInterval(60*60*2, 60*60*3),
			RxCenterFrequencyHz: 16000000000,
			RxBandwidthHz:       30000000,
			TxCenterFrequencyHz: 16000000000,
			TxBandwidthHz:       30000000,
		},
	})
	if err != nil {
		t.Fatalf("Test setup failed: %v", err)
	}

	tests := []struct {
		name      string
		bearerID  string
		bearer    *pb.Bearer
		wantError bool
	}{
		{
			name:     "Creates a new bearer",
			bearerID: "unique",
			bearer: &pb.Bearer{
				Name:                "bearers/unique",
				Target:              TARGET_NAME,
				Transceiver:         "transceivers/existing",
				Interval:            createInterval(20, 60*60),
				RxCenterFrequencyHz: 16000000000,
				RxBandwidthHz:       30000000,
				TxCenterFrequencyHz: 16000000000,
				TxBandwidthHz:       30000000,
			},
			wantError: false,
		},
		{
			name:     "Fails creating a bearer with already existing ID",
			bearerID: "existing",
			bearer: &pb.Bearer{
				Name:                "bearers/existing",
				Target:              TARGET_NAME,
				Transceiver:         "transceivers/existing",
				Interval:            createInterval(60*60*5, 60*60*6),
				RxCenterFrequencyHz: 16000000000,
				RxBandwidthHz:       30000000,
				TxCenterFrequencyHz: 16000000000,
				TxBandwidthHz:       30000000,
			}, wantError: true,
		},
		{
			name:     "Fails creating a bearer with interval outside of any contact window",
			bearerID: "false",
			bearer: &pb.Bearer{
				Name:                "bearers/false",
				Target:              TARGET_NAME,
				Transceiver:         "transceivers/existing",
				Interval:            createInterval(60*60*40, 60*60*41),
				RxCenterFrequencyHz: 16000000000,
				RxBandwidthHz:       30000000,
				TxCenterFrequencyHz: 16000000000,
				TxBandwidthHz:       30000000,
			}, wantError: true,
		},
		{
			name:     "Fails creating a bearer with negative interval",
			bearerID: "negative",
			bearer: &pb.Bearer{
				Name:                "bearers/negative",
				Target:              TARGET_NAME,
				Transceiver:         "transceivers/existing",
				Interval:            createInterval(60*60*20, 60*60*18),
				RxCenterFrequencyHz: 16000000000,
				RxBandwidthHz:       30000000,
				TxCenterFrequencyHz: 16000000000,
				TxBandwidthHz:       30000000,
			}, wantError: true,
		},
		{
			name:     "Fails creating a bearer with overlapping interval towards the end",
			bearerID: "overlap",
			bearer: &pb.Bearer{
				Name:                "bearers/overlap",
				Target:              TARGET_NAME,
				Transceiver:         "transceivers/existing",
				Interval:            createInterval(60*60*2+20, 60*60*3+20),
				RxCenterFrequencyHz: 16000000000,
				RxBandwidthHz:       30000000,
				TxCenterFrequencyHz: 16000000000,
				TxBandwidthHz:       30000000,
			}, wantError: true,
		},
		{
			name:     "Fails creating a bearer with contained interval",
			bearerID: "overlap",
			bearer: &pb.Bearer{
				Name:                "bearers/overlap",
				Target:              TARGET_NAME,
				Transceiver:         "transceivers/existing",
				Interval:            createInterval(60*60*2-20, 60*60*3-20),
				RxCenterFrequencyHz: 16000000000,
				RxBandwidthHz:       30000000,
				TxCenterFrequencyHz: 16000000000,
				TxBandwidthHz:       30000000,
			}, wantError: true,
		},
		{
			name:     "Fails creating a bearer with contained interval",
			bearerID: "overlap",
			bearer: &pb.Bearer{
				Name:                "bearers/overlap",
				Target:              TARGET_NAME,
				Transceiver:         "transceivers/existing",
				Interval:            createInterval(60*60*2-20, 60*60*3+20),
				RxCenterFrequencyHz: 16000000000,
				RxBandwidthHz:       30000000,
				TxCenterFrequencyHz: 16000000000,
				TxBandwidthHz:       30000000,
			}, wantError: true,
		},
		{
			name:     "Fails creating a bearer with wrong receiver center frequency",
			bearerID: "overlap",
			bearer: &pb.Bearer{
				Name:                "bearers/overlap",
				Target:              TARGET_NAME,
				Transceiver:         "transceivers/existing",
				Interval:            createInterval(60*60*20, 60*60*21),
				RxCenterFrequencyHz: 16,
				RxBandwidthHz:       30000000,
				TxCenterFrequencyHz: 16000000000,
				TxBandwidthHz:       30000000,
			}, wantError: true,
		},
		{
			name:     "Fails creating a bearer with wrong transmitter center frequency",
			bearerID: "overlap",
			bearer: &pb.Bearer{
				Name:                "bearers/overlap",
				Target:              TARGET_NAME,
				Transceiver:         "transceivers/existing",
				Interval:            createInterval(60*60*20, 60*60*21),
				RxCenterFrequencyHz: 16000000000,
				RxBandwidthHz:       30000000,
				TxCenterFrequencyHz: 16000,
				TxBandwidthHz:       30000000,
			}, wantError: true,
		},
		{
			name:     "Fails creating a bearer with wrong receiver bandwidth",
			bearerID: "overlap",
			bearer: &pb.Bearer{
				Name:                "bearers/overlap",
				Target:              TARGET_NAME,
				Transceiver:         "transceivers/existing",
				Interval:            createInterval(60*60*20, 60*60*21),
				RxCenterFrequencyHz: 16000000000,
				RxBandwidthHz:       3,
				TxCenterFrequencyHz: 16000000000,
				TxBandwidthHz:       30000000,
			}, wantError: true,
		},
		{
			name:     "Fails creating a bearer with wrong transmitter bandwidth",
			bearerID: "overlap",
			bearer: &pb.Bearer{
				Name:                "bearers/overlap",
				Target:              TARGET_NAME,
				Transceiver:         "transceivers/existing",
				Interval:            createInterval(60*60*20, 60*60*21),
				RxCenterFrequencyHz: 16000000000,
				RxBandwidthHz:       30000000,
				TxCenterFrequencyHz: 16000000000,
				TxBandwidthHz:       3000,
			}, wantError: true,
		},
		{
			name:     "Creates bearer with overlapping interval but non-overlapping frequencies",
			bearerID: "overlapfrequencies",
			bearer: &pb.Bearer{
				Name:                "bearers/overlapfrequencies",
				Target:              TARGET_NAME,
				Transceiver:         "transceivers/existing",
				Interval:            createInterval(60*60*20, 60*60*21),
				RxCenterFrequencyHz: 17000000000,
				RxBandwidthHz:       30000000,
				TxCenterFrequencyHz: 17000000000,
				TxBandwidthHz:       30000000,
			}, wantError: false,
		},
		{
			name:     "Returns correct name, even if initial name is not set correctly",
			bearerID: "correct",
			bearer: &pb.Bearer{
				Name:                "incorrect",
				Target:              TARGET_NAME,
				Transceiver:         "transceivers/existing",
				Interval:            createInterval(60*60*20, 60*60*21),
				RxCenterFrequencyHz: 16000000000,
				RxBandwidthHz:       30000000,
				TxCenterFrequencyHz: 16000000000,
				TxBandwidthHz:       30000000,
			}, wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.CreateBearerRequest{
				BearerId: tt.bearerID,
				Bearer:   tt.bearer,
			}

			resp, err := h.CreateBearer(ctx, req)
			if !tt.wantError && err != nil {
				t.Fatalf("Creating bearer in test '%s' failed %s.", tt.name, err)
			}
			if tt.wantError && err == nil {
				t.Fatalf("Creating bearer should have failed in test '%s', but did not.", tt.name)
			}
			if !tt.wantError && resp.Name != fmt.Sprintf("bearers/%s", tt.bearerID) {
				t.Fatalf("Should have updated name, but did not.")
			}
		})
	}
}

func TestPrototypeHandler_Bearer(t *testing.T) {
	h, ctx := createExistingTransceiver(t)

	_, err := h.CreateBearer(ctx, &pb.CreateBearerRequest{
		BearerId: "existing",
		Bearer: &pb.Bearer{
			Name:                "bearers/existing",
			Target:              TARGET_NAME,
			Transceiver:         "transceivers/existing",
			Interval:            createInterval(0, 60*60*1),
			RxCenterFrequencyHz: 16000000000,
			RxBandwidthHz:       30000000,
			TxCenterFrequencyHz: 16000000000,
			TxBandwidthHz:       30000000,
		},
	})
	if err != nil {
		t.Fatalf("Test setup failed: %v", err)
	}

	t.Run("Can get existing bearer with its information", func(t *testing.T) {
		resp, err := h.GetBearer(ctx, &pb.GetBearerRequest{
			Name: "bearers/existing",
		})
		if err != nil {
			t.Fatalf("Expected no error but was %v", err)
		}
		if diff := cmp.Diff(&pb.Bearer{
			Name:                "bearers/existing",
			Target:              TARGET_NAME,
			Transceiver:         "transceivers/existing",
			Interval:            createInterval(0, 60*60*1),
			RxCenterFrequencyHz: 16000000000,
			RxBandwidthHz:       30000000,
			TxCenterFrequencyHz: 16000000000,
			TxBandwidthHz:       30000000,
		}, resp,
			protocmp.Transform()); diff != "" {
			t.Errorf("Bearer mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Throws 404 error when bearer does not exist", func(t *testing.T) {
		_, err := h.GetBearer(ctx, &pb.GetBearerRequest{
			Name: "bearers/non-existing",
		})
		if err == nil {
			t.Fatal("Expected error but was none")
		}
		if !strings.HasPrefix(err.Error(), "rpc error: code = NotFound") {
			t.Fatalf("Expected error mismatch, got %v", err)
		}
	})

	t.Run("Can delete existing bearer", func(t *testing.T) {
		_, err := h.DeleteBearer(ctx, &pb.DeleteBearerRequest{
			Name: "bearers/existing",
		})
		if err != nil {
			t.Fatalf("Expected no error but was %v", err)
		}
	})

	t.Run("Throws error when deleting non-existant bearer", func(t *testing.T) {
		_, err := h.DeleteBearer(ctx, &pb.DeleteBearerRequest{
			Name: "bearers/non-existing",
		})
		if err == nil {
			t.Fatal("Expected error but there was none")
		}
	})

	_, err = h.CreateBearer(ctx, &pb.CreateBearerRequest{
		BearerId: "other",
		Bearer: &pb.Bearer{
			Name:                "bearers/other",
			Target:              TARGET_NAME,
			Transceiver:         "transceivers/existing",
			Interval:            createInterval(60*60*1, 60*60*2),
			RxCenterFrequencyHz: 16000000000,
			RxBandwidthHz:       30000000,
			TxCenterFrequencyHz: 16000000000,
			TxBandwidthHz:       30000000,
		},
	})
	if err != nil {
		t.Fatalf("Test setup failed: %v", err)
	}

	_, err = h.CreateAttachmentCircuit(ctx, &pb.CreateAttachmentCircuitRequest{
		AttachmentCircuitId: "other",
		AttachmentCircuit: &pb.AttachmentCircuit{
			Name:     "attachmentCircuits/other",
			Interval: createInterval(60*60*1, 60*60*2),
			L2Connection: &pb.AttachmentCircuit_L2Connection{
				Bearer: "bearers/other",
			},
		},
	})
	if err != nil {
		t.Fatalf("Test setup failed: %v", err)
	}

	t.Run("Cannot delete bearer with attached circuit", func(t *testing.T) {
		_, err := h.DeleteBearer(ctx, &pb.DeleteBearerRequest{
			Name: "bearers/other",
		})
		if err == nil {
			t.Fatal("Expected error but there was none")
		}
		if !strings.HasPrefix(err.Error(), "rpc error: code = FailedPrecondition desc = bearer has attachment circuit attached and cannot be deleted") {
			t.Fatalf("Mismatch in error message, got: %v", err)
		}
	})
}

func TestPrototypeHandler_CreateAttachmentCircuit(t *testing.T) {
	h, ctx := createExistingTransceiver(t)

	_, err := h.CreateBearer(ctx, &pb.CreateBearerRequest{
		BearerId: "existing",
		Bearer: &pb.Bearer{
			Name:                "bearers/existing",
			Target:              TARGET_NAME,
			Transceiver:         "transceivers/existing",
			Interval:            createInterval(0, 60*60*10),
			RxCenterFrequencyHz: 16000000000,
			RxBandwidthHz:       30000000,
			TxCenterFrequencyHz: 16000000000,
			TxBandwidthHz:       30000000,
		},
	})
	if err != nil {
		t.Fatalf("Test setup failed: %v", err)
	}

	_, err = h.CreateAttachmentCircuit(ctx, &pb.CreateAttachmentCircuitRequest{
		AttachmentCircuitId: "existing",
		AttachmentCircuit: &pb.AttachmentCircuit{
			Name:     "attachmentCircuits/existing",
			Interval: createInterval(60*60*2, 60*60*3),
			L2Connection: &pb.AttachmentCircuit_L2Connection{
				Bearer: "bearers/existing",
			},
		},
	})
	if err != nil {
		t.Fatalf("Test setup failed: %v", err)
	}

	tests := []struct {
		name      string
		acID      string
		ac        *pb.AttachmentCircuit
		wantError bool
	}{
		{
			name: "Creates a new attachment circuit",
			acID: "unique",
			ac: &pb.AttachmentCircuit{
				Name:     "attachmentCircuits/unique",
				Interval: createInterval(20, 60*60),
				L2Connection: &pb.AttachmentCircuit_L2Connection{
					Bearer: "bearers/existing",
				},
			},
			wantError: false,
		},
		{
			name: "Creates no new attachment circuit if underlying bearer does not exist",
			acID: "bearermissing",
			ac: &pb.AttachmentCircuit{
				Name:     "attachmentCircuits/bearermissing",
				Interval: createInterval(20, 60*60),
				L2Connection: &pb.AttachmentCircuit_L2Connection{
					Bearer: "bearers/non-existant",
				},
			},
			wantError: true,
		},
		{
			name: "Creates no new attachment circuit if underlying bearer is in wrong time interval",
			acID: "wrongtime",
			ac: &pb.AttachmentCircuit{
				Name:     "attachmentCircuits/wrongtime",
				Interval: createInterval(20, 60*60*40),
				L2Connection: &pb.AttachmentCircuit_L2Connection{
					Bearer: "bearers/existing",
				},
			},
			wantError: true,
		},
		{
			name: "Correctly refills the name when creating a new ac",
			acID: "newAC",
			ac: &pb.AttachmentCircuit{
				Name:     "attachmentCircuits/blablabla",
				Interval: createInterval(60*60*4, 60*60*5),
				L2Connection: &pb.AttachmentCircuit_L2Connection{
					Bearer: "bearers/existing",
				},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.CreateAttachmentCircuitRequest{
				AttachmentCircuitId: tt.acID,
				AttachmentCircuit:   tt.ac,
			}

			resp, err := h.CreateAttachmentCircuit(ctx, req)
			if !tt.wantError && err != nil {
				t.Fatalf("Creating attachment circuit in test '%s' failed %s.", tt.name, err)
			}
			if tt.wantError && err == nil {
				t.Fatalf("Creating attachment circuit should have failed in test '%s', but did not.", tt.name)
			}
			if !tt.wantError && resp.Name != fmt.Sprintf("attachmentCircuits/%s", tt.acID) {
				t.Fatalf("Should have updated name, but did not. Name was %s", resp.Name)
			}
		})
	}
}

func TestPrototypeHandler_AttachmentCircuits(t *testing.T) {
	h, ctx := createExistingTransceiver(t)

	_, err := h.CreateBearer(ctx, &pb.CreateBearerRequest{
		BearerId: "existing",
		Bearer: &pb.Bearer{
			Name:                "bearers/existing",
			Target:              TARGET_NAME,
			Transceiver:         "transceivers/existing",
			Interval:            createInterval(0, 60*60*10),
			RxCenterFrequencyHz: 16000000000,
			RxBandwidthHz:       30000000,
			TxCenterFrequencyHz: 16000000000,
			TxBandwidthHz:       30000000,
		},
	})
	if err != nil {
		t.Fatalf("Test setup failed: %v", err)
	}

	_, err = h.CreateAttachmentCircuit(ctx, &pb.CreateAttachmentCircuitRequest{
		AttachmentCircuitId: "existing",
		AttachmentCircuit: &pb.AttachmentCircuit{
			Name:     "attachmentCircuits/existing",
			Interval: createInterval(60*60*2, 60*60*3),
			L2Connection: &pb.AttachmentCircuit_L2Connection{
				Bearer: "bearers/existing",
			},
		},
	})
	if err != nil {
		t.Fatalf("Test setup failed: %v", err)
	}

	t.Run("Can get existing attachment circuit with its information", func(t *testing.T) {
		resp, err := h.GetAttachmentCircuit(ctx, &pb.GetAttachmentCircuitRequest{
			Name: "attachmentCircuits/existing",
		})
		if err != nil {
			t.Fatalf("Expected no error but was %v", err)
		}
		if diff := cmp.Diff(&pb.AttachmentCircuit{
			Name:     "attachmentCircuits/existing",
			Interval: createInterval(60*60*2, 60*60*3),
			L2Connection: &pb.AttachmentCircuit_L2Connection{
				Bearer: "bearers/existing",
			},
		}, resp,
			protocmp.Transform()); diff != "" {
			t.Errorf("AttachmentCircuit mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Throws 404 error when circuit does not exist", func(t *testing.T) {
		_, err := h.GetAttachmentCircuit(ctx, &pb.GetAttachmentCircuitRequest{
			Name: "attachmentCircuits/non-existing",
		})
		if err == nil {
			t.Fatal("Expected error but was none")
		}
		if !strings.HasPrefix(err.Error(), "rpc error: code = NotFound") {
			t.Fatalf("Expected error mismatch, got %v", err)
		}
	})

	t.Run("Can delete existing attachment cirucit", func(t *testing.T) {
		_, err := h.DeleteAttachmentCircuit(ctx, &pb.DeleteAttachmentCircuitRequest{
			Name: "attachmentCircuits/existing",
		})
		if err != nil {
			t.Fatalf("Expected no error but was %v", err)
		}
	})

	t.Run("Throws error when deleting non-existant circuit", func(t *testing.T) {
		_, err := h.DeleteAttachmentCircuit(ctx, &pb.DeleteAttachmentCircuitRequest{
			Name: "attachmentCircuits/non-existing",
		})
		if err == nil {
			t.Fatal("Expected error but there was none")
		}
	})
}

func createInterval(startTimeOffset int, endTimeOffset int) *interval.Interval {
	return &interval.Interval{
		StartTime: &timestamppb.Timestamp{
			Seconds: int64(time.Now().Unix()) + int64(startTimeOffset),
		},
		EndTime: &timestamppb.Timestamp{
			Seconds: int64(time.Now().Unix()) + int64(endTimeOffset),
		},
	}
}
