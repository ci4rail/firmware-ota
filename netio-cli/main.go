/*
Copyright © 2021 Ci4Rail GmbH
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"log"

	"github.com/ci4rail/firmware-ota/netio-cli/pkg/socket"
	"github.com/ci4rail/firmware-ota/netio-cli/pkg/version"
	"github.com/ci4rail/firmware-ota/pkg/netio_base"
	"github.com/golang/protobuf/proto"
)

var (
	host = "localhost:9999"
)

func main() {
	log.Printf("netio-cli version: %s. Host: %s\n", version.Version, host)

	sock, err := socket.NewConnection(host)

	if err != nil {
		log.Fatalf("Failed to create connection: %s", err)
	}

	protoData := &netio_base.Command{}
	protoData.Id = netio_base.CommandId_IDENTIFY_FIRMWARE
	payload, err := proto.Marshal(protoData)

	if err != nil {
		log.Fatalf("Failed to marshall: %s", err)
	}

	err = sock.Write(payload)
	if err != nil {
		log.Fatalf("Failed to write %s", err)
	}

}
