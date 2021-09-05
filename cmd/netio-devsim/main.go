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
	"io"
	"log"

	"github.com/ci4rail/firmware-ota/cmd/netio-devsim/internal/firmware"
	"github.com/ci4rail/firmware-ota/cmd/netio-devsim/pkg/version"
	"github.com/ci4rail/firmware-ota/pkg/netio"
	"github.com/ci4rail/firmware-ota/pkg/netio/transport"
)

var (
	port = ":9999"
)

func main() {
	log.Printf("netio-devsim version: %s listen at port %s\n", version.Version, port)

	listener, err := transport.NewListener(port)
	if err != nil {
		log.Fatalf("Failed to create listenet: %s", err)
	}

	if err != nil {
		log.Fatalf("Failed to create devproto: %s", err)
	}
	for {
		conn, err := transport.WaitForConnect(listener)
		if err != nil {
			log.Fatalf("Failed to wait for connection: %s", err)
		}

		ch, _ := netio.NewChannel(conn)

		serveConnection(ch)
	}
}

func serveConnection(ch *netio.Channel) {
	defer ch.Close()

	for {
		c := &netio.BaseFuncCommand{}
		err := ch.ReadMessage(c)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Fatalf("Failed to read: %s", err)
		}

		var res *netio.BaseFuncResponse
		switch c.Id {
		case netio.BaseFuncCommandId_IDENTIFY_FIRMWARE:
			res = firmware.IdentifyFirmware()
		case netio.BaseFuncCommandId_LOAD_FIRMWARE_CHUNK:
			res = firmware.LoadFirmwareChunk(c.GetLoadFirmwareChunk())
		default:
			res = &netio.BaseFuncResponse{
				Id:     c.Id,
				Status: netio.BaseFuncStatus_UNKNOWN_COMMAND,
			}
		}

		err = ch.WriteMessage(res)
		if err != nil {
			log.Fatalf("Failed to write: %s", err)
		}
	}
}
