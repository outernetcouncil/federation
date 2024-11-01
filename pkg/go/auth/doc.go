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

// Package auth provides server-side authentication mechanisms for gRPC
// services using JSON Web Tokens (JWT) and RSA public/private key pairs.
//
// To use this package, create a Config and initialize the interceptors:
//
//	config := &auth.Config{
//		PublicKey: publicKeyReader,
//	}
//	unaryInterceptor, streamInterceptor, err := auth.InitializeServerInterceptors(config)
//	if err != nil {
//		// Handle error
//	}
//	server := grpc.NewServer(
//		grpc.UnaryInterceptor(unaryInterceptor),
//		grpc.StreamInterceptor(streamInterceptor),
//	)
//
// For more information about authentication, see https://docs.spacetime.aalyria.com/authentication
package auth
