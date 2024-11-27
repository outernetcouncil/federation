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
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/outernetcouncil/federation/gen/go/federation/v1alpha"
	inet "outernetcouncil.org/nmts/proto/types/ietf"
)

type PrototypeHandler struct {
	pb.UnimplementedFederationServiceServer
	mu               sync.Mutex
	services         map[string]*pb.ServiceStatus
	interconnections map[string]*pb.InterconnectionPoint
}

func NewPrototypeHandler() *PrototypeHandler {
	return &PrototypeHandler{
		services:         make(map[string]*pb.ServiceStatus),
		interconnections: make(map[string]*pb.InterconnectionPoint),
	}
}

func (h *PrototypeHandler) StreamInterconnectionPoints(req *pb.StreamInterconnectionPointsRequest, stream pb.FederationService_StreamInterconnectionPointsServer) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	chunk := &pb.StreamInterconnectionPointsResponse{
		SnapshotComplete: &pb.StreamInterconnectionPointsResponse_SnapshotComplete{},
	}

	for _, ip := range h.interconnections {
		chunk.Mutations = append(chunk.Mutations, &pb.InterconnectionMutation{
			Type: &pb.InterconnectionMutation_Upsert{
				Upsert: ip,
			},
		})
	}

	if err := stream.Send(chunk); err != nil {
		return err
	}

	if req.SnapshotOnly != nil && *req.SnapshotOnly {
		return nil
	}

	<-stream.Context().Done()
	return stream.Context().Err()
}

func (h *PrototypeHandler) ListServiceOptions(req *pb.ListServiceOptionsRequest, stream pb.FederationService_ListServiceOptionsServer) error {
	ipNetwork := &inet.IPNetwork{
		Prefix: &inet.IPPrefix{
			Version: &inet.IPPrefix_Ipv4{
				Ipv4: &inet.IPv4Prefix{
					Str: "192.168.1.0/24",
				},
			},
		},
		// Optionally set realm if needed
		Realm: "",
	}

	response := &pb.ListServiceOptionsResponse{
		ServiceOptions: []*pb.ServiceOption{
			{
				Id: "sample-service-option",
				EndpointX: &pb.ServiceOption_XProviderInterconnection{
					XProviderInterconnection: &pb.InterconnectionPoint{
						Uuid: "sample-provider-interconnection",
					},
				},
				EndpointY: &pb.ServiceOption_IpNetwork{
					IpNetwork: ipNetwork,
				},
				Directionality: pb.Directionality_DIRECTIONALITY_BIDIRECTIONAL_UNSPECIFIED,
			},
		},
	}

	return stream.Send(response)
}

func (h *PrototypeHandler) ScheduleService(ctx context.Context, req *pb.ScheduleServiceRequest) (*pb.ScheduleServiceResponse, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	serviceID := fmt.Sprintf("service-%d", len(h.services)+1)
	h.services[serviceID] = &pb.ServiceStatus{
		Id: serviceID,
		Type: &pb.ServiceStatus_ServiceUpdate_{
			ServiceUpdate: &pb.ServiceStatus_ServiceUpdate{
				IsActive: true,
			},
		},
	}

	return &pb.ScheduleServiceResponse{
		ServiceId: serviceID,
	}, nil
}

func (h *PrototypeHandler) MonitorServices(stream pb.FederationService_MonitorServicesServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}

		h.mu.Lock()
		for _, serviceID := range req.AddServiceIds {
			if status, exists := h.services[serviceID]; exists {
				if err := stream.Send(&pb.MonitorServicesResponse{
					UpdatedServices: []*pb.ServiceStatus{status},
				}); err != nil {
					h.mu.Unlock()
					return err
				}
			}
		}
		h.mu.Unlock()
	}
}

func (h *PrototypeHandler) CancelService(ctx context.Context, req *pb.CancelServiceRequest) (*pb.CancelServiceResponse, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.services[req.ServiceId]; !exists {
		return nil, status.Errorf(codes.NotFound, "service %s not found", req.ServiceId)
	}

	delete(h.services, req.ServiceId)
	return &pb.CancelServiceResponse{
		Cancelled: true,
	}, nil
}
