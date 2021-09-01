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
	"fmt"
	"log"

	"github.com/ci4rail/firmware-ota/netio-devsim/pkg/socket"
	"github.com/ci4rail/firmware-ota/netio-devsim/pkg/version"
	"github.com/ci4rail/firmware-ota/pkg/netio_base"
	"github.com/golang/protobuf/proto"
)

var (
	port = ":9999"
)

func main() {
	log.Printf("netio-devsim version: %s listen at port %s\n", version.Version, port)

	sock, err := socket.WaitForConnect(port)

	if err != nil {
		log.Fatalf("Failed to wait for connection: %s", err)
	}

	for {
		payload, err := sock.Read()

		if err != nil {
			log.Fatalf("Failed to read from connection: %s", err)
		}

		fmt.Printf("got payload %v\n", payload)

		cmdData := &netio_base.Command{}
		if err := proto.Unmarshal(payload, cmdData); err != nil {
			log.Fatalf("Failed to unmarshal: %s", err)
		}
		fmt.Printf("Got command %v\n", cmdData.Id)

		var res netobase.Response

		switch cmdData.Id {
		case netio_base.CommandId_IDENTIFY_FIRMWARE:
			res = IdentifyFirmware()
		default:
			res := &netio_base.Response{
				Id:  cmdData.Id,
				Status: netio_base.Status_UNKNOWN_COMMAND,
			}
		}

	}

}

func IdentifyFirmware() netio_base.Response {

	res := &netio_base.Response{
		Id: netio_base.CommandId_IDENTIFY_FIRMWARE,
		Status: netio_base.Status_OK,
		Data: &netio_base.Response_IdentifyFirmware{
			IdentifyFirmware: &netio_base.ResIdentifyFirmware{
				Name: "bla",
				MajorVersion: 1,
				MinorVersion: 0,
				PatchVersion: 0,
			},
		}
	}
	return res	
}
