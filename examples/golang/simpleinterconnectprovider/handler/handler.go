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
	"log"
	"math"
	"strings"
	"sync"
	"time"

	"google.golang.org/genproto/googleapis/type/interval"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"outernetcouncil.org/nmts/v1alpha/proto/ek/physical"
	"outernetcouncil.org/nmts/v1alpha/proto/types/geophys"

	pb "github.com/outernetcouncil/federation/gen/go/federation/interconnect/v1alpha"
)

const TARGET_NAME = "target/mysat"

type PrototypeHandler struct {
	pb.UnimplementedInterconnectServiceServer
	mu                 sync.Mutex
	transceivers       map[string]*pb.Transceiver
	bearers            map[string]*pb.Bearer
	target             map[string]*pb.Target
	attachmentCircuits map[string]*pb.AttachmentCircuit
	// Contact windows are mostly going to be computed on the fly. To simplify the example, we just "compute" them whenever a transceiver is created.
	contactWindows []*pb.ContactWindow
}

// We pretend to be a very simple provider with one target only.
func NewPrototypeHandler() *PrototypeHandler {
	providerTarget := pb.Target{
		Name:   TARGET_NAME,
		Motion: &geophys.Motion{},
	}
	targets := make(map[string]*pb.Target)
	targets[providerTarget.Name] = &providerTarget

	return &PrototypeHandler{
		transceivers:       make(map[string]*pb.Transceiver),
		bearers:            make(map[string]*pb.Bearer),
		target:             targets,
		attachmentCircuits: make(map[string]*pb.AttachmentCircuit),
		contactWindows:     make([]*pb.ContactWindow, 0, 1),
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
	if p.transceivers[trans.Name] == nil {
		return nil, status.Errorf(codes.NotFound, "transceiver with requested ID was not found")
	}

	return p.transceivers[trans.Name], nil
}

func (p *PrototypeHandler) CreateTransceiver(_ context.Context, trans *pb.CreateTransceiverRequest) (*pb.Transceiver, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	transceiverName := fmt.Sprintf("transceivers/%s", trans.TransceiverId)
	if p.transceivers[transceiverName] != nil {
		return nil, status.Errorf(codes.AlreadyExists, "transceiver with requested ID was already created")
	}
	if err := checkForAdmissibleTransceiver(trans.Transceiver); err != nil {
		return nil, err
	}
	// Override the name of the attachment circuit to ensure that it has the correct resource name.
	// It is up to the API to either validate the correctness of the name or just override it on creation.
	trans.Transceiver.Name = transceiverName
	p.transceivers[transceiverName] = trans.Transceiver

	for _, target := range p.target {
		targetID := strings.Split(target.Name, "/")[1]
		p.contactWindows = append(p.contactWindows, &pb.ContactWindow{
			Name: fmt.Sprintf("contactWindow/%s%s", trans.TransceiverId, targetID),
			Interval: &interval.Interval{
				StartTime: &timestamppb.Timestamp{
					Seconds: int64(time.Now().Unix()),
				},
				EndTime: &timestamppb.Timestamp{
					Seconds: int64(time.Now().Unix()) + 60*60*24, // let's just have a one day window everywhere
				},
			},
			Transceiver:            fmt.Sprintf("transceivers/%s", trans.TransceiverId),
			Target:                 target.Name,
			MinRxCenterFrequencyHz: 12000000000,
			MaxRxCenterFrequencyHz: 18000000000,
			MinRxBandwidthHz:       20000000,
			MaxRxBandwidthHz:       40000000,
			MinTxCenterFrequencyHz: 12000000000,
			MaxTxCenterFrequencyHz: 18000000000,
			MinTxBandwidthHz:       20000000,
			MaxTxBandwidthHz:       40000000,
		})
	}

	return trans.Transceiver, nil
}

func checkForAdmissibleTransceiver(trans *pb.Transceiver) error {
	if trans.ReceiveSignalChain == nil ||
		trans.TransmitSignalChain == nil ||
		trans.ReceiveSignalChain.Antenna == nil ||
		trans.TransmitSignalChain.Antenna == nil ||
		trans.TransmitSignalChain.Antenna.Type != physical.Antenna_OPTICAL ||
		trans.ReceiveSignalChain.Antenna.Type != physical.Antenna_OPTICAL {
		return status.Errorf(codes.FailedPrecondition, "transceiver is not compatible, see ListCompatibleTransceiverTypes for details")
	}

	return nil
}

func (p *PrototypeHandler) UpdateTransceiver(_ context.Context, trans *pb.UpdateTransceiverRequest) (*pb.Transceiver, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.transceivers[trans.Transceiver.Name] == nil {
		return nil, status.Errorf(codes.NotFound, "transceiver with requested ID was not found")
	}
	if err := checkForAdmissibleTransceiver(trans.Transceiver); err != nil {
		return nil, err
	}

	for _, bearer := range p.bearers {
		// In this example, we simply prohibit that a client update their transceiver if it is used in a connection.
		// In a real API implementation, more complicated logic could be applied to ensure that it is actually possible to update.
		if bearer.Transceiver == trans.Transceiver.Name {
			return nil, status.Error(codes.FailedPrecondition, "transceiver has bearer attached and cannot be updated")
		}
	}

	p.transceivers[trans.Transceiver.Name] = trans.Transceiver

	return trans.Transceiver, nil
}

func (p *PrototypeHandler) DeleteTransceiver(_ context.Context, trans *pb.DeleteTransceiverRequest) (*emptypb.Empty, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.transceivers[trans.Name] == nil {
		return nil, status.Errorf(codes.NotFound, "transceiver with requested ID was not found")
	}
	for _, bearer := range p.bearers {
		// In order to ensure that the connection setup is valid, we need to check for attached bearers.
		if bearer.Transceiver == trans.Name {
			return nil, status.Error(codes.FailedPrecondition, "transceiver has bearer attached and cannot be deleted")
		}
	}

	delete(p.transceivers, trans.Name)

	newWindows := make([]*pb.ContactWindow, 0, len(p.contactWindows))
	for _, window := range p.contactWindows {
		if window.Transceiver != trans.Name {
			newWindows = append(newWindows, window)
		}
	}
	p.contactWindows = newWindows

	return &emptypb.Empty{}, nil
}

func (p *PrototypeHandler) ListContactWindows(context.Context, *pb.ListContactWindowsRequest) (*pb.ListContactWindowsResponse, error) {
	// TODO: Handle filter
	return &pb.ListContactWindowsResponse{
		ContactWindows: p.contactWindows,
	}, nil
}

func (p *PrototypeHandler) ListBearers(context.Context, *pb.ListBearersRequest) (*pb.ListBearersResponse, error) {
	// TODO: Handle filter
	bearers := make([]*pb.Bearer, 0, len(p.attachmentCircuits))
	for _, bearer := range p.bearers {
		bearers = append(bearers, bearer)
	}

	return &pb.ListBearersResponse{Bearers: bearers}, nil
}

func (p *PrototypeHandler) GetBearer(_ context.Context, bearer *pb.GetBearerRequest) (*pb.Bearer, error) {
	if p.bearers[bearer.Name] == nil {
		return nil, status.Errorf(codes.NotFound, "bearer with requested ID was not found")
	}

	return p.bearers[bearer.Name], nil
}

func (p *PrototypeHandler) CreateBearer(_ context.Context, bearer *pb.CreateBearerRequest) (*pb.Bearer, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	bearerName := fmt.Sprintf("bearers/%s", bearer.BearerId)
	if p.bearers[bearerName] != nil {
		return nil, status.Errorf(codes.AlreadyExists, "bearer with requested ID was already created")
	}
	if bearer.Bearer.Interval.StartTime.AsTime().After(bearer.Bearer.Interval.EndTime.AsTime()) {
		return nil, status.Errorf(codes.InvalidArgument, "bearer has negative time interval argument")
	}
	if !p.checkForSufficientContactWindow(bearer.Bearer) {
		return nil, status.Errorf(codes.FailedPrecondition, "bearer has no sufficient contact window")
	}
	// Override the name of the attachment circuit to ensure that it has the correct resource name.
	// It is up to the API to either validate the correctness of the name or just override it on creation.
	bearer.Bearer.Name = bearerName

	p.bearers[bearerName] = bearer.Bearer

	return bearer.Bearer, nil
}

func (p *PrototypeHandler) checkForSufficientContactWindow(bearer *pb.Bearer) bool {
	for _, contactWindow := range p.contactWindows {
		if contactWindow.Target != bearer.Target || contactWindow.Transceiver != bearer.Transceiver {
			continue
		}

		if contactWindow.Interval.StartTime.AsTime().After(bearer.Interval.StartTime.AsTime()) ||
			contactWindow.Interval.EndTime.AsTime().Before(bearer.Interval.EndTime.AsTime()) {
			continue
		}

		if contactWindow.MinRxCenterFrequencyHz >= bearer.RxCenterFrequencyHz ||
			contactWindow.MaxRxCenterFrequencyHz <= bearer.RxCenterFrequencyHz ||
			contactWindow.MinRxBandwidthHz >= bearer.RxBandwidthHz ||
			contactWindow.MaxRxBandwidthHz <= bearer.RxBandwidthHz ||
			contactWindow.MinTxCenterFrequencyHz >= bearer.TxCenterFrequencyHz ||
			contactWindow.MaxTxCenterFrequencyHz <= bearer.TxCenterFrequencyHz ||
			contactWindow.MinTxBandwidthHz >= bearer.TxBandwidthHz ||
			contactWindow.MaxTxBandwidthHz <= bearer.TxBandwidthHz {
			continue
		}

		for _, knownBearer := range p.bearers {
			if bearer.Target != knownBearer.Target || bearer.Transceiver != knownBearer.Transceiver {
				continue
			}

			if math.Abs(float64(knownBearer.RxCenterFrequencyHz-bearer.RxCenterFrequencyHz))-(float64(knownBearer.RxBandwidthHz)+float64(bearer.RxBandwidthHz))/2 > 0 &&
				math.Abs(float64(knownBearer.TxCenterFrequencyHz-bearer.TxCenterFrequencyHz))-(float64(knownBearer.TxBandwidthHz)+float64(bearer.TxBandwidthHz))/2 > 0 {
				return true
			}

			if knownBearer.Interval.EndTime.AsTime().After(bearer.Interval.StartTime.AsTime()) && bearer.Interval.StartTime.AsTime().After(knownBearer.Interval.StartTime.AsTime()) ||
				knownBearer.Interval.EndTime.AsTime().After(bearer.Interval.EndTime.AsTime()) && bearer.Interval.EndTime.AsTime().After(knownBearer.Interval.StartTime.AsTime()) ||
				bearer.Interval.EndTime.AsTime().After(knownBearer.Interval.EndTime.AsTime()) && knownBearer.Interval.StartTime.AsTime().After(bearer.Interval.StartTime.AsTime()) {
				log.Println("Bearers are overlapping.")

				return false
			}
		}

		return true
	}

	log.Println("Could not find any sufficient contact window.")

	return false
}

func (p *PrototypeHandler) DeleteBearer(_ context.Context, bearer *pb.DeleteBearerRequest) (*emptypb.Empty, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.bearers[bearer.Name] == nil {
		return nil, status.Errorf(codes.NotFound, "bearer with requested ID was not found")
	}
	for _, ac := range p.attachmentCircuits {
		// In order to ensure that the connection setup is valid, we need to check attached bearers.
		if bearer.Name == ac.L2Connection.Bearer {
			return nil, status.Error(codes.FailedPrecondition, "bearer has attachment circuit attached and cannot be deleted")
		}
	}
	delete(p.bearers, bearer.Name)

	return &emptypb.Empty{}, nil
}

func (p *PrototypeHandler) ListAttachmentCircuits(context.Context, *pb.ListAttachmentCircuitsRequest) (*pb.ListAttachmentCircuitsResponse, error) {
	// TODO: Handle filter
	attachmentCircuits := make([]*pb.AttachmentCircuit, 0, len(p.attachmentCircuits))
	for _, circuit := range p.attachmentCircuits {
		attachmentCircuits = append(attachmentCircuits, circuit)
	}

	return &pb.ListAttachmentCircuitsResponse{AttachmentCircuits: attachmentCircuits}, nil
}

func (p *PrototypeHandler) GetAttachmentCircuit(_ context.Context, circuit *pb.GetAttachmentCircuitRequest) (*pb.AttachmentCircuit, error) {
	if p.attachmentCircuits[circuit.Name] == nil {
		return nil, status.Errorf(codes.NotFound, "attachment circuit with requested ID was not found")
	}

	return p.attachmentCircuits[circuit.Name], nil
}

func (p *PrototypeHandler) CreateAttachmentCircuit(_ context.Context, ac *pb.CreateAttachmentCircuitRequest) (*pb.AttachmentCircuit, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	attachmentCircuitName := fmt.Sprintf("attachmentCircuits/%s", ac.AttachmentCircuitId)
	if p.bearers[attachmentCircuitName] != nil {
		return nil, status.Errorf(codes.AlreadyExists, "attachment circuit with requested ID was already created")
	}
	if !p.checkForSufficientBearer(ac.AttachmentCircuit) {
		return nil, status.Errorf(codes.FailedPrecondition, "attachment circuit is not attached to existing bearer covering the provisioning window")
	}

	// Override the name of the attachment circuit to ensure that it has the correct resource name.
	// It is up to the API to either validate the correctness of the name or just override it on creation.
	ac.AttachmentCircuit.Name = attachmentCircuitName
	p.attachmentCircuits[attachmentCircuitName] = ac.AttachmentCircuit

	return ac.AttachmentCircuit, nil
}

func (p *PrototypeHandler) checkForSufficientBearer(attachmentCircuit *pb.AttachmentCircuit) bool {
	for _, bearer := range p.bearers {
		if bearer.Name != attachmentCircuit.L2Connection.Bearer {
			continue
		}

		if bearer.Interval.StartTime.AsTime().After(attachmentCircuit.Interval.StartTime.AsTime()) ||
			bearer.Interval.EndTime.AsTime().Before(attachmentCircuit.Interval.EndTime.AsTime()) {
			return false
		}

		return true
	}

	return false
}

func (p *PrototypeHandler) DeleteAttachmentCircuit(_ context.Context, request *pb.DeleteAttachmentCircuitRequest) (*emptypb.Empty, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.attachmentCircuits[request.Name] == nil {
		return nil, status.Errorf(codes.NotFound, "attachment circuit with requested ID was not found")
	}
	delete(p.attachmentCircuits, request.Name)

	return &emptypb.Empty{}, nil
}

func (p *PrototypeHandler) GetTarget(_ context.Context, targetRequest *pb.GetTargetRequest) (*pb.Target, error) {
	if p.target[targetRequest.Name] != nil {
		return nil, status.Errorf(codes.NotFound, "could not find target")
	}

	return p.target[targetRequest.Name], nil
}

func (p *PrototypeHandler) ListTargets(context.Context, *pb.ListTargetsRequest) (*pb.ListTargetsResponse, error) {
	targets := make([]*pb.Target, 0, len(p.target))
	for _, target := range targets {
		targets = append(targets, target)
	}
	targetResponse := &pb.ListTargetsResponse{
		Targets: targets,
	}

	return targetResponse, nil
}
