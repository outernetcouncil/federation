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
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"

	pb "github.com/outernetcouncil/federation/gen/go/federation/v1alpha"
	inet "outernetcouncil.org/nmts/proto/types/ietf"
)

func TestPrototypeHandler_ScheduleService(t *testing.T) {
	h := NewPrototypeHandler()
	ctx := context.Background()

	req := &pb.ScheduleServiceRequest{
		Type: &pb.ScheduleServiceRequest_RequestorEdgeToIpNetwork{
			RequestorEdgeToIpNetwork: &pb.RequestorEdgeToIpNetwork{},
		},
	}

	resp, err := h.ScheduleService(ctx, req)
	if err != nil {
		t.Fatalf("ScheduleService failed: %v", err)
	}

	if resp.ServiceId == "" {
		t.Error("Expected non-empty ServiceId")
	}

	if len(h.services) != 1 {
		t.Errorf("Expected 1 service, got %d", len(h.services))
	}

	service, exists := h.services[resp.ServiceId]
	if !exists {
		t.Errorf("Service %s not found in handler", resp.ServiceId)
	}

	if service.Id != resp.ServiceId {
		t.Errorf("Service ID mismatch. Expected %s, got %s", resp.ServiceId, service.Id)
	}

	if update, ok := service.Type.(*pb.ServiceStatus_ServiceUpdate_); !ok || !update.ServiceUpdate.IsActive {
		t.Error("Expected service to be active")
	}
}

func TestPrototypeHandler_CancelService(t *testing.T) {
	h := NewPrototypeHandler()
	ctx := context.Background()

	// Schedule a service first
	scheduleReq := &pb.ScheduleServiceRequest{
		Type: &pb.ScheduleServiceRequest_RequestorEdgeToIpNetwork{
			RequestorEdgeToIpNetwork: &pb.RequestorEdgeToIpNetwork{},
		},
	}
	scheduleResp, _ := h.ScheduleService(ctx, scheduleReq)

	tests := []struct {
		name        string
		serviceID   string
		expectError bool
		errorCode   codes.Code
	}{
		{
			name:        "Cancel existing service",
			serviceID:   scheduleResp.ServiceId,
			expectError: false,
		},
		{
			name:        "Cancel non-existent service",
			serviceID:   "non-existent-id",
			expectError: true,
			errorCode:   codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &pb.CancelServiceRequest{
				ServiceId: tt.serviceID,
			}

			resp, err := h.CancelService(ctx, req)

			if tt.expectError {
				if err == nil {
					t.Fatal("Expected an error, but got nil")
				}
				if status.Code(err) != tt.errorCode {
					t.Errorf("Expected error code %v, got %v", tt.errorCode, status.Code(err))
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if !resp.Cancelled {
					t.Error("Expected Cancelled to be true")
				}
				if _, exists := h.services[tt.serviceID]; exists {
					t.Error("Service still exists after cancellation")
				}
			}
		})
	}
}

func TestPrototypeHandler_StreamInterconnectionPoints(t *testing.T) {
	tests := []struct {
		name         string
		snapshotOnly bool
		points       []*pb.InterconnectionPoint
		wantErr      bool
	}{
		{
			name:         "Empty snapshot",
			snapshotOnly: true,
			points:       []*pb.InterconnectionPoint{},
		},
		{
			name:         "Single point snapshot",
			snapshotOnly: true,
			points: []*pb.InterconnectionPoint{
				{
					Uuid: "test-point-1",
				},
			},
		},
		{
			name:         "Multiple points streaming",
			snapshotOnly: false,
			points: []*pb.InterconnectionPoint{
				{Uuid: "test-point-1"},
				{Uuid: "test-point-2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewPrototypeHandler()

			// Populate test data
			h.mu.Lock()
			for _, p := range tt.points {
				h.interconnections[p.Uuid] = p
			}
			h.mu.Unlock()

			// Create mock stream
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			stream := &mockStreamInterconnectionPointsServer{
				ctx:      ctx,
				recvChan: make(chan *pb.StreamInterconnectionPointsRequest, 1),
				sendChan: make(chan *pb.StreamInterconnectionPointsResponse, 10),
			}

			// Send initial request
			req := &pb.StreamInterconnectionPointsRequest{
				SnapshotOnly: &tt.snapshotOnly,
			}

			// Run stream in goroutine
			errChan := make(chan error, 1)
			go func() {
				errChan <- h.StreamInterconnectionPoints(req, stream)
			}()

			// Verify received chunks
			chunk := <-stream.sendChan
			if len(chunk.Mutations) != len(tt.points) {
				t.Errorf("Got %d mutations, want %d", len(chunk.Mutations), len(tt.points))
			}

			if !tt.snapshotOnly {
				// Verify streaming continues until context cancelled
				select {
				case <-errChan:
					t.Error("Stream ended prematurely")
				case <-ctx.Done():
					// Expected behavior
				}
			}
		})
	}
}

// Mock stream implementation
type mockStreamInterconnectionPointsServer struct {
	ctx      context.Context
	recvChan chan *pb.StreamInterconnectionPointsRequest
	sendChan chan *pb.StreamInterconnectionPointsResponse
	grpc.ServerStream
}

func (m *mockStreamInterconnectionPointsServer) Context() context.Context {
	return m.ctx
}

func (m *mockStreamInterconnectionPointsServer) Send(chunk *pb.StreamInterconnectionPointsResponse) error {
	select {
	case m.sendChan <- chunk:
		return nil
	case <-m.ctx.Done():
		return m.ctx.Err()
	}
}

func TestPrototypeHandler_ListServiceOptions(t *testing.T) {
	tests := []struct {
		name string
		req  *pb.ListServiceOptionsRequest
		want []*pb.ServiceOption
	}{
		{
			name: "Basic service options",
			req:  &pb.ListServiceOptionsRequest{},
			want: []*pb.ServiceOption{
				{
					Id: "sample-service-option",
					EndpointX: &pb.ServiceOption_XProviderInterconnection{
						XProviderInterconnection: &pb.InterconnectionPoint{
							Uuid: "sample-provider-interconnection",
						},
					},
					EndpointY: &pb.ServiceOption_IpNetwork{
						IpNetwork: &inet.IPNetwork{
							Prefix: &inet.IPPrefix{
								Version: &inet.IPPrefix_Ipv4{
									Ipv4: &inet.IPv4Prefix{
										Str: "192.168.1.0/24",
									},
								},
							},
						},
					},
					Directionality: pb.Directionality_DIRECTIONALITY_BIDIRECTIONAL,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewPrototypeHandler()
			stream := &mockListServiceOptionsServer{
				ctx:       context.Background(),
				responses: make(chan *pb.ListServiceOptionsResponse, 1),
			}

			err := h.ListServiceOptions(tt.req, stream)
			if err != nil {
				t.Fatalf("ListServiceOptions() error = %v", err)
			}

			select {
			case resp := <-stream.responses:
				// Add proper comparison options for protobuf messages
				opts := []cmp.Option{
					cmpopts.IgnoreUnexported(
						pb.ServiceOption{},
						pb.InterconnectionPoint{},
						pb.ServiceOption_XProviderInterconnection{},
						pb.ServiceOption_IpNetwork{},
						inet.IPNetwork{},
						inet.IPPrefix{},
						inet.IPv4Prefix{},
					),
					protocmp.Transform(),
				}
				if diff := cmp.Diff(tt.want, resp.ServiceOptions, opts...); diff != "" {
					t.Errorf("ListServiceOptions() mismatch (-want +got):\n%s", diff)
				}
			case <-time.After(time.Second):
				t.Fatal("timeout waiting for response")
			}
		})
	}
}

type mockListServiceOptionsServer struct {
	ctx       context.Context
	responses chan *pb.ListServiceOptionsResponse
	grpc.ServerStream
}

func (m *mockListServiceOptionsServer) Context() context.Context {
	return m.ctx
}

func (m *mockListServiceOptionsServer) Send(resp *pb.ListServiceOptionsResponse) error {
	select {
	case m.responses <- resp:
		return nil
	case <-m.ctx.Done():
		return m.ctx.Err()
	}
}

func TestPrototypeHandler_MonitorServices(t *testing.T) {
	tests := []struct {
		name     string
		services map[string]*pb.ServiceStatus
		updates  []*pb.MonitorServicesRequest
		want     []*pb.ServiceStatus
	}{
		{
			name: "Monitor single service",
			services: map[string]*pb.ServiceStatus{
				"service-1": {
					Id: "service-1",
					Type: &pb.ServiceStatus_ServiceUpdate_{
						ServiceUpdate: &pb.ServiceStatus_ServiceUpdate{
							IsActive: true,
						},
					},
				},
			},
			updates: []*pb.MonitorServicesRequest{
				{AddServiceIds: []string{"service-1"}},
			},
			want: []*pb.ServiceStatus{
				{
					Id: "service-1",
					Type: &pb.ServiceStatus_ServiceUpdate_{
						ServiceUpdate: &pb.ServiceStatus_ServiceUpdate{
							IsActive: true,
						},
					},
				},
			},
		},
		// Add more test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewPrototypeHandler()

			// Setup initial services
			h.mu.Lock()
			h.services = tt.services
			h.mu.Unlock()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			stream := newMockMonitorServicesServer(ctx, len(tt.updates))

			// Start monitoring in background
			errChan := make(chan error, 1)
			go func() {
				errChan <- h.MonitorServices(stream)
			}()

			// Send updates
			for _, update := range tt.updates {
				if err := stream.sendRequest(update, time.Second); err != nil {
					t.Fatalf("Failed to send request: %v", err)
				}
			}

			// Verify responses
			for i := 0; i < len(tt.want); i++ {
				resp, err := stream.receiveUpdate(time.Second)
				if err != nil {
					t.Fatalf("Failed to receive update: %v", err)
				}

				if len(resp.UpdatedServices) != 1 {
					t.Fatalf("Expected 1 updated service, got %d", len(resp.UpdatedServices))
				}

				if diff := cmp.Diff(tt.want[i], resp.UpdatedServices[0],
					protocmp.Transform(),
					cmpopts.IgnoreUnexported(
						pb.ServiceStatus{},
						pb.ServiceStatus_ServiceUpdate{},
						pb.ServiceStatus_ServiceUpdate_{},
					)); diff != "" {
					t.Errorf("MonitorServices() mismatch (-want +got):\n%s", diff)
				}
			}

			// Clean shutdown
			cancel()
			if err := <-errChan; err != nil && err != context.Canceled {
				t.Errorf("MonitorServices() unexpected error: %v", err)
			}
		})
	}
}

type mockMonitorServicesServer struct {
	ctx      context.Context
	requests chan *pb.MonitorServicesRequest
	updates  chan *pb.MonitorServicesResponse
	grpc.ServerStream
}

func (m *mockMonitorServicesServer) Context() context.Context {
	return m.ctx
}

func (m *mockMonitorServicesServer) Send(resp *pb.MonitorServicesResponse) error {
	select {
	case m.updates <- resp:
		return nil
	case <-m.ctx.Done():
		return m.ctx.Err()
	}
}

func (m *mockMonitorServicesServer) Recv() (*pb.MonitorServicesRequest, error) {
	select {
	case req := <-m.requests:
		return req, nil
	case <-m.ctx.Done():
		return nil, m.ctx.Err()
	}
}

// Helper method to create a new mock server with reasonable defaults
func newMockMonitorServicesServer(ctx context.Context, bufferSize int) *mockMonitorServicesServer {
	if ctx == nil {
		ctx = context.Background()
	}
	if bufferSize <= 0 {
		bufferSize = 10
	}
	return &mockMonitorServicesServer{
		ctx:      ctx,
		requests: make(chan *pb.MonitorServicesRequest, bufferSize),
		updates:  make(chan *pb.MonitorServicesResponse, bufferSize),
	}
}

// Helper methods for testing
func (m *mockMonitorServicesServer) sendRequest(req *pb.MonitorServicesRequest, timeout time.Duration) error {
	if timeout <= 0 {
		timeout = time.Second
	}
	select {
	case m.requests <- req:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timeout sending request")
	case <-m.ctx.Done():
		return m.ctx.Err()
	}
}

func (m *mockMonitorServicesServer) receiveUpdate(timeout time.Duration) (*pb.MonitorServicesResponse, error) {
	if timeout <= 0 {
		timeout = time.Second
	}
	select {
	case update := <-m.updates:
		return update, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout receiving update")
	case <-m.ctx.Done():
		return nil, m.ctx.Err()
	}
}
