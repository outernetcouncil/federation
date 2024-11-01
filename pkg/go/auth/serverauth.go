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

// Package serverauth provides server-side authentication mechanisms for gRPC
// services using JSON Web Tokens (JWT) and RSA public/private key pairs. It
// includes functionality to load and parse public keys, validate JWT tokens,
// and integrate with gRPC interceptors for both unary and stream RPCs. The
// package ensures that incoming requests are authenticate by verifying the JWT
// tokens against the provided public key, enhancing the security of the gRPC
// services.
package serverauth // import "aalyria.com/spacetime/serverauth"

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authHeader = "authorization"
)

type Config struct {
	// Client's public key, used to validate signature created with client private key
	PublicKey io.Reader
}

// Validate JWT client-provided tokens against the client's publicKey
func validateJWT(tokenString string, publicKey interface{}) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parsing token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func InitializeServerInterceptors(sc *Config) (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor, error) {
	publicKey, err := LoadPublicKey(sc.PublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load public key: %w", err)
	}

	return UnaryServerInterceptor(publicKey), StreamServerInterceptor(publicKey), nil
}

// Unary interceptor for server-side JWT validation
func UnaryServerInterceptor(publicKey any) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
		}

		authHeader, ok := md[authHeader]
		if !ok || len(authHeader) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "missing authorization header")
		}

		tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")
		_, err := validateJWT(tokenString, publicKey)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		return handler(ctx, req)
	}
}

// Stream interceptor for server-side JWT validation
func StreamServerInterceptor(publicKey any) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return status.Errorf(codes.InvalidArgument, "missing metadata")
		}

		authHeader, ok := md[authHeader]
		if !ok || len(authHeader) == 0 {
			return status.Errorf(codes.Unauthenticated, "missing authorization header")
		}

		tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")
		_, err := validateJWT(tokenString, publicKey)
		if err != nil {
			return status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		return handler(srv, ss)
	}
}

// LoadPublicKey reads, decodes, and parses the public key from the provided io.Reader
func LoadPublicKey(r io.Reader) (interface{}, error) {
	pubKeyBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading public key: %w", err)
	}
	if len(pubKeyBytes) == 0 {
		return nil, errors.New("empty public key")
	}

	pubKeyBlock, _ := pem.Decode(pubKeyBytes)
	if pubKeyBlock == nil {
		return nil, errors.New("public key not PEM-encoded")
	}

	pubKey, err := x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parsing public key: %w", err)
	}

	return pubKey, nil
}
