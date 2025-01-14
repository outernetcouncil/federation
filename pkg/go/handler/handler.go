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

// package handler defines the interface for handling Federation gRPC requests.
package handler

import (
	"context"

	pb "github.com/outernetcouncil/federation/gen/go/federation/v1alpha"
)

// FederationHandler defines the interface for handling Federation gRPC requests.
type FederationHandler interface {
	pb.FederationServiceServer
	StreamInterconnectionPoints(*pb.StreamInterconnectionPointsRequest, pb.FederationService_StreamInterconnectionPointsServer) error
	ListServiceOptions(*pb.ListServiceOptionsRequest, pb.FederationService_ListServiceOptionsServer) error
	ScheduleService(context.Context, *pb.ScheduleServiceRequest) (*pb.ScheduleServiceResponse, error)
	MonitorServices(pb.FederationService_MonitorServicesServer) error
	CancelService(context.Context, *pb.CancelServiceRequest) (*pb.CancelServiceResponse, error)
}
