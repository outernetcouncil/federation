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
	"outernetcouncil.org/nmts/v1alpha/proto/types/geophys"
)

func TestPrototypeHandler_Targets(t *testing.T) {
	h := NewPrototypeHandler()
	ctx := context.Background()

	t.Run("GetTarget returns current known targets", func(t *testing.T) {
		resp, err := h.GetTarget(ctx, &pb.GetTargetRequest{
			Name: TARGET_NAME,
		})

		if err != nil {
			t.Fatalf("error not expected but was %s", err)
		}
		if diff := cmp.Diff(&pb.Target{
			Name:   TARGET_NAME,
			Motion: &geophys.Motion{},
		}, resp,
			protocmp.Transform()); diff != "" {
			t.Errorf("Target mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("GetTarget returns error if target is unknown", func(t *testing.T) {
		_, err := h.GetTarget(ctx, &pb.GetTargetRequest{
			Name: "unknown name",
		})

		if err == nil {
			t.Fatalf("error expected")
		}
		if !strings.HasPrefix(err.Error(), "rpc error: code = NotFound") {
			t.Fatalf("expected error did not match error %v", err)
		}
	})

	t.Run("ListTarget returns a list of all targets", func(t *testing.T) {
		resp, err := h.ListTargets(ctx, &pb.ListTargetsRequest{})

		if err != nil {
			t.Fatal("error expected")
		}
		if len(resp.Targets) != 1 {
			t.Fatal("unexpected number of targets")
		}
	})
}

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
		wantError   string
	}{
		{
			name: "Allows update to compatible receiver",
			transceiver: &pb.Transceiver{
				Name: "transceivers/existing",
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
			wantError: "",
		},
		{
			name: "Fails updating a transceiver with non-existing ID",
			transceiver: &pb.Transceiver{
				Name: "transceivers/non-existant",
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
			wantError: "NotFound",
		},
		{
			name: "Fails updating a transceiver with wrong antenna type in transmit signal chain",
			transceiver: &pb.Transceiver{
				Name: "transceivers/existing",
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
			wantError: "FailedPrecondition",
		},
		{
			name: "Fails updating a transceiver with wrong antenna type in receive signal chain",
			transceiver: &pb.Transceiver{
				Name: "transceivers/existing",
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
			wantError: "FailedPrecondition",
		},
		{
			name: "Fails updating a transceiver with unspecified antenna type",
			transceiver: &pb.Transceiver{
				Name: "transceivers/existing",
			},
			wantError: "FailedPrecondition",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.UpdateTransceiverRequest{
				Transceiver: tt.transceiver,
			}

			_, err := h.UpdateTransceiver(ctx, req)
			if tt.wantError == "" && err != nil {
				t.Fatalf("Updating a transceiver failed %s.", err)
			}
			if tt.wantError != "" && !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("Updating a transceiver has wrong error message %v.", err.Error())
			}
		})
	}
}

func TestPrototypeHandler_Transceivers(t *testing.T) {
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
				Name: "transceivers/existing",
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
		if !strings.Contains(err.Error(), "has bearer attached") {
			t.Fatal("error should relate to no bearer being attached but was %s", err.Error())
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

	t.Run("ListTransceivers lists all bearers", func(t *testing.T) {
		resp, err := h.ListTransceivers(ctx, &pb.ListTransceiversRequest{})

		if err != nil {
			t.Fatalf("Expected no error but was %v", err)
		}
		if len(resp.Transceivers) != 1 {
			t.Fatalf("Unexpected number of attachment circuits")
		}
	})

	t.Run("ListTransceivers throws error on filter", func(t *testing.T) {
		_, err := h.ListTransceivers(ctx, &pb.ListTransceiversRequest{
			Filter: "test filter",
		})

		if err == nil {
			t.Fatal("Error expected")
		}
		if !strings.HasPrefix(err.Error(), "rpc error: code = Unimplemented") {
			t.Fatalf("expected error did not match error %v", err)
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
			name:     "Allows creation of adjacent bearers in time",
			bearerID: "adjacent",
			bearer: &pb.Bearer{
				Name:                "bearers/adjacent",
				Target:              TARGET_NAME,
				Transceiver:         "transceivers/existing",
				Interval:            createInterval(60*60, 60*60*2),
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
			name:     "Creates bearer with edge frequencies",
			bearerID: "edgefrequencies",
			bearer: &pb.Bearer{
				Name:                "bearers/edgefrequencies",
				Target:              TARGET_NAME,
				Transceiver:         "transceivers/existing",
				Interval:            createInterval(60*60*20, 60*60*21),
				RxCenterFrequencyHz: 12000000000,
				RxBandwidthHz:       40000000,
				TxCenterFrequencyHz: 18000000000,
				TxBandwidthHz:       40000000,
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

	t.Run("ListBearers lists all bearers", func(t *testing.T) {
		resp, err := h.ListBearers(ctx, &pb.ListBearersRequest{})

		if err != nil {
			t.Fatalf("Expected no error but was %v", err)
		}
		if len(resp.Bearers) != 1 {
			t.Fatalf("Unexpected number of attachment circuits")
		}
	})

	t.Run("ListBearers throws error on filter", func(t *testing.T) {
		_, err := h.ListBearers(ctx, &pb.ListBearersRequest{
			Filter: "test filter",
		})

		if err == nil {
			t.Fatal("Error expected")
		}
		if !strings.HasPrefix(err.Error(), "rpc error: code = Unimplemented") {
			t.Fatalf("expected error did not match error %v", err)
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
			name: "Does not create a new attachment circuit if it already exists",
			acID: "unique",
			ac: &pb.AttachmentCircuit{
				Name:     "attachmentCircuits/existing",
				Interval: createInterval(60*60*4, 60*60*5),
				L2Connection: &pb.AttachmentCircuit_L2Connection{
					Bearer: "bearers/existing",
				},
			},
			wantError: true,
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

	t.Run("ListAttachmentCircuits lists all circuits", func(t *testing.T) {
		resp, err := h.ListAttachmentCircuits(ctx, &pb.ListAttachmentCircuitsRequest{})

		if err != nil {
			t.Fatalf("Expected no error but was %v", err)
		}
		if len(resp.AttachmentCircuits) != 0 { // All attachment circuits were deleted
			t.Fatalf("Unexpected number of attachment circuits")
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
			t.Fatalf("Setup for ListAttachmentCircuit failed: %v", err)
		}

		resp, err = h.ListAttachmentCircuits(ctx, &pb.ListAttachmentCircuitsRequest{})

		if err != nil {
			t.Fatalf("Expected no error but was %v", err)
		}
		if len(resp.AttachmentCircuits) != 1 {
			t.Fatalf("Unexpected number of attachment circuits")
		}
	})

	t.Run("ListAttachmentCircuits throws error on filter", func(t *testing.T) {
		_, err := h.ListAttachmentCircuits(ctx, &pb.ListAttachmentCircuitsRequest{
			Filter: "test filter",
		})

		if err == nil {
			t.Fatal("Error expected")
		}
		if !strings.HasPrefix(err.Error(), "rpc error: code = Unimplemented") {
			t.Fatalf("expected error did not match error %v", err)
		}
	})
}

func TestPrototypeHandler_ListCompatibleTransceiverTypes(t *testing.T) {
	h, ctx := createExistingTransceiver(t)

	t.Run("ListCompatibleTransceiverTypes gets us some information about compatible transceivers", func(t *testing.T) {
		response, err := h.ListCompatibleTransceiverTypes(ctx, &pb.ListCompatibleTransceiverTypesRequest{})

		if err != nil {
			t.Fatalf("expected no error, but was %v", err)
		}
		if diff := cmp.Diff(&pb.ListCompatibleTransceiverTypesResponse{
			CompatibleTransceiverTypes: []*pb.CompatibleTransceiverType{
				{
					TransceiverFilter: "transmit_signal_chain.antenna.type = OPTICAL AND receive_signal_chain.antenna.type = OPTICAL",
				},
			},
		}, response,
			protocmp.Transform()); diff != "" {
			t.Fatalf("Compatible transceivers mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestPrototypeHandler_ListContactWindows(t *testing.T) {
	h, ctx := createExistingTransceiver(t)

	t.Run("ListContactWindows gets a list of all contact windows", func(t *testing.T) {
		response, err := h.ListContactWindows(ctx, &pb.ListContactWindowsRequest{})

		if err != nil {
			t.Fatalf("expected no error, but was %v", err)
		}
		if len(response.ContactWindows) != 1 { // We created one transceiver, so there is one contact window
			t.Fatalf("Unexected number of contact windows")
		}
	})

	t.Run("ListContactWindows returns unimplemented when used with a filter", func(t *testing.T) {
		_, err := h.ListContactWindows(ctx, &pb.ListContactWindowsRequest{
			Filter: "test filter",
		})

		if err == nil {
			t.Fatal("Error expected")
		}
		if !strings.HasPrefix(err.Error(), "rpc error: code = Unimplemented") {
			t.Fatalf("expected error did not match error %v", err)
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
