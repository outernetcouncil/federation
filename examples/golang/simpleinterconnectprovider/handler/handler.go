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
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"outernetcouncil.org/nmts/v1alpha/proto/ek/physical"

	pb "github.com/outernetcouncil/federation/gen/go/federation/interconnect/v1alpha"
)

type PrototypeHandler struct {
	pb.UnimplementedInterconnectServiceServer
	mu                 sync.Mutex
	transceivers       map[string]*pb.Transceiver
	bearers            map[string]*pb.Bearer
	target             map[string]*pb.Target
	attachmentCircuits map[string]*pb.AttachmentCircuit
}

func NewPrototypeHandler() *PrototypeHandler {
	return &PrototypeHandler{
		transceivers:       make(map[string]*pb.Transceiver),
		bearers:            make(map[string]*pb.Bearer),
		target:             make(map[string]*pb.Target),
		attachmentCircuits: make(map[string]*pb.AttachmentCircuit),
	}
}

func (p *PrototypeHandler) ListCompatibleTransceiverTypes(context.Context, *pb.ListCompatibleTransceiverTypesRequest) (*pb.ListCompatibleTransceiverTypesResponse, error) {
	// TODO: Is this enough as an example? Constraining the antenna types?
	return &pb.ListCompatibleTransceiverTypesResponse{
		CompatibleTransceiverTypes: []*pb.CompatibleTransceiverType{
			{
				TransceiverFilter: "transmit_signal_chain.antenna.type = OPTICAL AND receive_signal_chain.antenna.type = OPTICAL",
			},
		},
	}, nil
}

func (p *PrototypeHandler) GetTransceiver(ctx context.Context, trans *pb.GetTransceiverRequest) (*pb.Transceiver, error) {
	// TODO: Is name here the same as ID? It's a bit weird that the request has a separate ID field, but the update, get and Delete does not...
	if p.transceivers[trans.Name] == nil {
		return nil, status.Errorf(codes.NotFound, "transceiver with requested ID was not found")
	}

	return p.transceivers[trans.Name], nil
}

func (p *PrototypeHandler) CreateTransceiver(_ context.Context, trans *pb.CreateTransceiverRequest) (*pb.Transceiver, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.transceivers[trans.TransceiverId] != nil {
		return nil, status.Errorf(codes.AlreadyExists, "transceiver with requested ID was already created")
	}
	if trans.Transceiver.ReceiveSignalChain == nil ||
		trans.Transceiver.TransmitSignalChain == nil ||
		trans.Transceiver.ReceiveSignalChain.Antenna == nil ||
		trans.Transceiver.TransmitSignalChain.Antenna == nil ||
		trans.Transceiver.TransmitSignalChain.Antenna.Type != physical.Antenna_OPTICAL ||
		trans.Transceiver.ReceiveSignalChain.Antenna.Type != physical.Antenna_OPTICAL {
		return nil, status.Errorf(codes.FailedPrecondition, "transceiver is not compatible, see ListCompatibleTransceiverTypes for details")
	}
	p.transceivers[trans.TransceiverId] = trans.Transceiver

	return trans.Transceiver, nil
}

func (p *PrototypeHandler) UpdateTransceiver(_ context.Context, trans *pb.UpdateTransceiverRequest) (*pb.Transceiver, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// TODO: Is name here the same as ID? It's a bit weird that the request has a separate ID field, but the update, get and Delete does not...
	if p.transceivers[trans.Transceiver.Name] == nil {
		return nil, status.Errorf(codes.NotFound, "transceiver with requested ID was not found")
	}
	p.transceivers[trans.Transceiver.Name] = trans.Transceiver

	return trans.Transceiver, nil
}

func (p *PrototypeHandler) DeleteTransceiver(_ context.Context, trans *pb.DeleteTransceiverRequest) (*emptypb.Empty, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// TODO: Is name here the same as ID? It's a bit weird that the request has a separate ID field, but the update, get and Delete does not...
	if p.transceivers[trans.Name] == nil {
		return nil, status.Errorf(codes.NotFound, "transceiver with requested ID was not found")
	}
	delete(p.transceivers, trans.Name)

	return &emptypb.Empty{}, nil
}

func (p *PrototypeHandler) ListContactWindows(context.Context, *pb.ListContactWindowsRequest) (*pb.ListContactWindowsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListContactWindows not implemented")
}
func (p *PrototypeHandler) ListBearers(context.Context, *pb.ListBearersRequest) (*pb.ListBearersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListBearers not implemented")
}
func (p *PrototypeHandler) CreateBearer(context.Context, *pb.CreateBearerRequest) (*pb.Bearer, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBearer not implemented")
}
func (p *PrototypeHandler) DeleteBearer(context.Context, *pb.DeleteBearerRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteBearer not implemented")
}
func (p *PrototypeHandler) ListAttachmentCircuits(context.Context, *pb.ListAttachmentCircuitsRequest) (*pb.ListAttachmentCircuitsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAttachmentCircuits not implemented")
}
func (p *PrototypeHandler) CreateAttachmentCircuit(context.Context, *pb.CreateAttachmentCircuitRequest) (*pb.AttachmentCircuit, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAttachmentCircuit not implemented")
}
func (p *PrototypeHandler) DeleteAttachmentCircuit(context.Context, *pb.DeleteAttachmentCircuitRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAttachmentCircuit not implemented")
}
func (p *PrototypeHandler) GetTarget(context.Context, *pb.GetTargetRequest) (*pb.Target, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTarget not implemented")
}
func (p *PrototypeHandler) ListTargets(context.Context, *pb.ListTargetsRequest) (*pb.ListTargetsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTargets not implemented")
}
