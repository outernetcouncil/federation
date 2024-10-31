// Copyright (c) Outernet Council and Contributors.
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

package main

import (
	"os"

	federationpb "github.com/outernetcouncil/federation/gen/go/federation/v1alpha"
	"google.golang.org/protobuf/proto"
)

func main() {
	e := federationpb.CancelServiceRequest{
		ServiceId: "service-id",
	}

	buf, _ := proto.Marshal(&e)
	os.WriteFile("federation.protobuf", buf, os.ModePerm)
}
