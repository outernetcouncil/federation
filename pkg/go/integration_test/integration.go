package integration_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/outernetcouncil/federation/gen/go/federation/interconnect/v1alpha"
	"github.com/outernetcouncil/federation/pkg/go/cosmicconnector"
	"github.com/outernetcouncil/federation/pkg/go/server"
)

// testHandler implements the FederationHandler interface for testing
type testHandler struct {
	pb.UnimplementedInterconnectServiceServer
}

func createTestHandler(t *testing.T) *testHandler {
	return &testHandler{}
}

func (h *testHandler) ListCompatibleTransceiverTypes(context.Context, *pb.ListCompatibleTransceiverTypesRequest) (*pb.ListCompatibleTransceiverTypesResponse, error) {
	return nil, nil
}

func (h *testHandler) GetTransceiver(context.Context, *pb.GetTransceiverRequest) (*pb.Transceiver, error) {
	return nil, nil
}

func (h *testHandler) CreateTransceiver(context.Context, *pb.CreateTransceiverRequest) (*pb.Transceiver, error) {
	return nil, nil
}

func (h *testHandler) UpdateTransceiver(context.Context, *pb.UpdateTransceiverRequest) (*pb.Transceiver, error) {
	return nil, nil
}

func (h *testHandler) DeleteTransceiver(context.Context, *pb.DeleteTransceiverRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (h *testHandler) ListContactWindows(context.Context, *pb.ListContactWindowsRequest) (*pb.ListContactWindowsResponse, error) {
	return nil, nil
}

func (h *testHandler) ListBearers(context.Context, *pb.ListBearersRequest) (*pb.ListBearersResponse, error) {
	return nil, nil
}

func (h *testHandler) CreateBearer(context.Context, *pb.CreateBearerRequest) (*pb.Bearer, error) {
	return nil, nil
}

func (h *testHandler) DeleteBearer(context.Context, *pb.DeleteBearerRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (h *testHandler) ListAttachmentCircuits(context.Context, *pb.ListAttachmentCircuitsRequest) (*pb.ListAttachmentCircuitsResponse, error) {
	return nil, nil
}

func (h *testHandler) CreateAttachmentCircuit(context.Context, *pb.CreateAttachmentCircuitRequest) (*pb.AttachmentCircuit, error) {
	return nil, nil
}

func (h *testHandler) DeleteAttachmentCircuit(context.Context, *pb.DeleteAttachmentCircuitRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (h *testHandler) GetTarget(context.Context, *pb.GetTargetRequest) (*pb.Target, error) {
	return nil, nil
}

func (h *testHandler) ListTargets(context.Context, *pb.ListTargetsRequest) (*pb.ListTargetsResponse, error) {
	return nil, nil
}

type testServers struct {
	grpcAddr     string
	channelzAddr string
	pprofAddr    string
	cc           *cosmicconnector.CosmicConnector
	cleanup      func()
}

// Helper function to fail test with message
func failTest(t *testing.T, format string, args ...interface{}) {
	t.Helper()
	t.Fatalf(format, args...)
}

func setupTestServers(t *testing.T) *testServers {
	t.Helper()

	// Get random available ports
	grpcPort := getFreePort(t)
	channelzPort := getFreePort(t)
	pprofPort := getFreePort(t)

	logger := zerolog.New(zerolog.NewTestWriter(t))

	// Create real handler implementation
	handler := createTestHandler(t)

	// Create servers with random ports
	grpcServer := server.NewGrpcServer(grpcPort, handler, logger)
	channelzServer := server.NewChannelzServer(fmt.Sprintf(":%d", channelzPort), logger)
	pprofServer := server.NewPprofServer(fmt.Sprintf(":%d", pprofPort), logger)

	cc := cosmicconnector.NewCosmicConnector(
		logger,
		grpcServer,
		channelzServer,
		pprofServer,
	)

	// Start servers in background
	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		errCh <- cc.Run(ctx)
	}()

	// Wait for servers to start
	time.Sleep(100 * time.Millisecond)

	return &testServers{
		grpcAddr:     fmt.Sprintf(":%d", grpcPort),
		channelzAddr: fmt.Sprintf(":%d", channelzPort),
		pprofAddr:    fmt.Sprintf(":%d", pprofPort),
		cc:           cc,
		cleanup: func() {
			cancel()
			<-errCh
		},
	}
}

func TestBasicServerLifecycle(t *testing.T) {
	servers := setupTestServers(t)
	defer servers.cleanup()

	// Test 1: Verify all servers are listening
	t.Run("servers_listening", func(t *testing.T) {
		// Check gRPC server
		conn, err := net.Dial("tcp", servers.grpcAddr)
		if err != nil {
			failTest(t, "failed to connect to gRPC server: %v", err)
		}
		conn.Close()

		// Check Channelz server
		conn, err = net.Dial("tcp", servers.channelzAddr)
		if err != nil {
			failTest(t, "failed to connect to Channelz server: %v", err)
		}
		conn.Close()

		// Check pprof server
		resp, err := http.Get(fmt.Sprintf("http://localhost%s/debug/pprof/", servers.pprofAddr))
		if err != nil {
			failTest(t, "failed to connect to pprof server: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			failTest(t, "unexpected pprof response status: got %d, want %d",
				resp.StatusCode, http.StatusOK)
		}
	})

	// Test 2: Verify gRPC connectivity
	t.Run("grpc_connectivity", func(t *testing.T) {
		conn, err := grpc.Dial(
			servers.grpcAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			failTest(t, "failed to create gRPC connection: %v", err)
		}
		defer conn.Close()

		client := pb.NewInterconnectServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Make a simple RPC call
		_, err = client.ListBearers(ctx, &pb.ListBearersRequest{})
		if err != nil {
			failTest(t, "ListServiceOptions failed: %v", err)
		}
	})
}

// Helper function to get an available port
func getFreePort(t *testing.T) int {
	t.Helper()

	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		failTest(t, "failed to resolve TCP address: %v", err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		failTest(t, "failed to listen on TCP port: %v", err)
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port
}
