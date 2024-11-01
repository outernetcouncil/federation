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

package serverauth

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1"
	"google.golang.org/grpc/status"
)

var testKey = generateRSAPrivateKey()

type rsaKeyForTesting struct {
	privateKey *rsa.PrivateKey
	privatePEM []byte
	publicPEM  []byte
}

func generateRSAPrivateKey() rsaKeyForTesting {
	bitSize := 2048
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		panic(err)
	}

	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	pubBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		panic(err)
	}
	pubPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: pubBytes,
		},
	)

	return rsaKeyForTesting{
		privateKey: key,
		privatePEM: keyPEM,
		publicPEM:  pubPEM,
	}
}

func TestLoadPublicKey(t *testing.T) {
	for _, tc := range []struct {
		name    string
		keyData []byte
		wantErr string
	}{
		{
			name:    "valid public key",
			keyData: testKey.publicPEM,
			wantErr: "",
		},
		{
			name:    "empty public key",
			keyData: []byte{},
			wantErr: "empty public key",
		},
		{
			name:    "invalid public key",
			keyData: []byte("invalid key data"),
			wantErr: "public key not PEM-encoded",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := LoadPublicKey(bytes.NewReader(tc.keyData))
			if err != nil && err.Error() != tc.wantErr {
				t.Errorf("unexpected error: got %v, want %q", err, tc.wantErr)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	validToken, err := generateValidJWT(testKey.privateKey)
	if err != nil {
		t.Fatalf("failed to generate valid JWT: %v", err)
	}

	for _, tc := range []struct {
		name      string
		token     string
		publicKey interface{}
		wantErr   string
	}{
		{
			name:      "valid token",
			token:     validToken,
			publicKey: testKey.privateKey.Public(),
			wantErr:   "",
		},
		{
			name:      "invalid token",
			token:     "invalid.token",
			publicKey: testKey.privateKey.Public(),
			wantErr:   "parsing token: token is malformed: token contains an invalid number of segments",
		},
		{
			name:      "invalid signing method",
			token:     validToken,
			publicKey: nil,
			wantErr:   "parsing token: token signature is invalid: key is of invalid type: RSA verify expects *rsa.PublicKey",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, err := validateJWT(tc.token, tc.publicKey)
			if err != nil {
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Errorf("unexpected error: got %v, want %v", err, tc.wantErr)
				}
			} else if tc.wantErr != "" {
				t.Errorf("expected error %q, got nil", tc.wantErr)
			}
		})
	}
}

func generateValidJWT(privateKey *rsa.PrivateKey) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": "test",
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour)),
	})
	return token.SignedString(privateKey)
}

func TestUnaryServerInterceptor(t *testing.T) {
	validToken, err := generateValidJWT(testKey.privateKey)
	if err != nil {
		t.Fatalf("failed to generate valid JWT: %v", err)
	}

	for _, tc := range []struct {
		name      string
		md        metadata.MD
		publicKey interface{}
		wantCode  codes.Code
	}{
		{
			name:      "valid token",
			md:        metadata.Pairs(authHeader, "Bearer "+validToken),
			publicKey: testKey.privateKey.Public(),
			wantCode:  codes.OK,
		},
		{
			name:      "missing metadata",
			md:        nil,
			publicKey: testKey.privateKey.Public(),
			wantCode:  codes.Unauthenticated,
		},
		{
			name:      "missing authorization header",
			md:        metadata.Pairs("some-header", "some-value"),
			publicKey: testKey.privateKey.Public(),
			wantCode:  codes.Unauthenticated,
		},
		{
			name:      "invalid token",
			md:        metadata.Pairs(authHeader, "Bearer invalid.token"),
			publicKey: testKey.privateKey.Public(),
			wantCode:  codes.Unauthenticated,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			interceptor := UnaryServerInterceptor(tc.publicKey)
			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				return "response", nil
			}

			ctx := metadata.NewIncomingContext(context.Background(), tc.md)
			_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{}, handler)
			st, _ := status.FromError(err)
			if st.Code() != tc.wantCode {
				t.Errorf("unexpected status code: got %v, want %v", st.Code(), tc.wantCode)
			}
		})
	}
}

func TestStreamServerInterceptor(t *testing.T) {
	validToken, err := generateValidJWT(testKey.privateKey)
	if err != nil {
		t.Fatalf("failed to generate valid JWT: %v", err)
	}

	for _, tc := range []struct {
		name      string
		md        metadata.MD
		publicKey interface{}
		wantCode  codes.Code
	}{
		{
			name:      "valid token",
			md:        metadata.Pairs(authHeader, "Bearer "+validToken),
			publicKey: testKey.privateKey.Public(),
			wantCode:  codes.OK,
		},
		{
			name:      "missing metadata",
			md:        nil,
			publicKey: testKey.privateKey.Public(),
			wantCode:  codes.Unauthenticated,
		},
		{
			name:      "missing authorization header",
			md:        metadata.Pairs("some-header", "some-value"),
			publicKey: testKey.privateKey.Public(),
			wantCode:  codes.Unauthenticated,
		},
		{
			name:      "invalid token",
			md:        metadata.Pairs(authHeader, "Bearer invalid.token"),
			publicKey: testKey.privateKey.Public(),
			wantCode:  codes.Unauthenticated,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			interceptor := StreamServerInterceptor(tc.publicKey)
			handler := func(srv interface{}, stream grpc.ServerStream) error {
				return nil
			}

			ctx := metadata.NewIncomingContext(context.Background(), tc.md)
			ss := &mockServerStream{ctx: ctx}
			err := interceptor(nil, ss, &grpc.StreamServerInfo{}, handler)
			st, _ := status.FromError(err)
			if st.Code() != tc.wantCode {
				t.Errorf("unexpected status code: got %v, want %v", st.Code(), tc.wantCode)
			}
		})
	}
}

type mockServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (m *mockServerStream) Context() context.Context {
	return m.ctx
}

func TestGRPCServerWithReflection(t *testing.T) {
	// Generate a valid JWT token
	validToken, err := generateValidJWT(testKey.privateKey)
	if err != nil {
		t.Fatalf("failed to generate valid JWT: %v", err)
	}

	// Load the public key
	publicKey, err := LoadPublicKey(bytes.NewReader(testKey.publicPEM))
	if err != nil {
		t.Fatalf("failed to load public key: %v", err)
	}

	// Create the gRPC server with the interceptors
	unaryInterceptor := UnaryServerInterceptor(publicKey)
	streamInterceptor := StreamServerInterceptor(publicKey)
	server := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor),
		grpc.StreamInterceptor(streamInterceptor),
	)

	// Register the reflection service on the gRPC server
	reflection.Register(server)

	// Start the gRPC server in a goroutine
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	go server.Serve(lis)
	defer server.Stop()

	// Create a gRPC client connection
	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial server: %v", err)
	}
	defer conn.Close()

	// Create a client for the reflection service
	client := reflectpb.NewServerReflectionClient(conn)

	tests := []struct {
		name      string
		token     string
		wantCode  codes.Code
		wantError string
	}{
		{
			name:     "valid token",
			token:    validToken,
			wantCode: codes.OK,
		},
		{
			name:      "invalid token",
			token:     "invalid.token",
			wantCode:  codes.Unauthenticated,
			wantError: "invalid token",
		},
		{
			name:      "missing token",
			token:     "",
			wantCode:  codes.Unauthenticated,
			wantError: "missing authorization header",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			if tc.token != "" {
				ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(authHeader, "Bearer "+tc.token))
			}

			err := makeReflectionRequest(ctx, client)
			if err != nil {
				st, _ := status.FromError(err)
				if st.Code() != tc.wantCode {
					t.Errorf("unexpected status code: got %v, want %v", st.Code(), tc.wantCode)
				}
				if tc.wantError != "" && !strings.Contains(st.Message(), tc.wantError) {
					t.Errorf("unexpected error message: got %v, want %v", st.Message(), tc.wantError)
				}
			} else if tc.wantCode != codes.OK {
				t.Errorf("expected error, got success")
			}
		})
	}
}

func makeReflectionRequest(ctx context.Context, client reflectpb.ServerReflectionClient) error {
	stream, err := client.ServerReflectionInfo(ctx)
	if err != nil {
		return err
	}

	if err := stream.Send(&reflectpb.ServerReflectionRequest{
		Host:           "",
		MessageRequest: &reflectpb.ServerReflectionRequest_ListServices{},
	}); err != nil {
		return err
	}

	_, err = stream.Recv()
	return err
}
