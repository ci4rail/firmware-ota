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
	"time"

	"github.com/ci4rail/firmware-ota/cmd/netio-devsim/internal/firmware"
	"github.com/ci4rail/firmware-ota/cmd/netio-devsim/pkg/version"
	"github.com/ci4rail/firmware-ota/pkg/netio"
	"github.com/ci4rail/firmware-ota/pkg/netio/basefunc"
	"github.com/ci4rail/firmware-ota/pkg/netio/transport"
)

var (
	port = ":9999"
)

func main() {
	log.Printf("netio-devsim version: %s listen at port %s\n", version.Version, port)

	listener, err := transport.NewSocketListener(port)
	if err != nil {
		log.Fatalf("Failed to create listener: %s", err)
	}

	if err != nil {
		log.Fatalf("Failed to create devproto: %s", err)
	}
	for {
		conn, err := transport.WaitForSocketConnect(listener)
		if err != nil {
			log.Fatalf("Failed to wait for connection: %s", err)
		}
		log.Printf("new connection!\n")

		ms, _ := transport.NewMsgStreamFromConnection(conn)

		ch, _ := netio.NewChannel(ms)

		serveConnection(ch)
		time.Sleep(4 * time.Second) // simulate reboot
	}
}

func serveConnection(ch *netio.Channel) {
	defer ch.Close()

	for {
		c := &basefunc.BaseFuncCommand{}
		err := ch.ReadMessage(c, 0)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Fatalf("Failed to read: %s", err)
		}

		var res *basefunc.BaseFuncResponse
		doreset := false
		switch c.Id {
		case basefunc.BaseFuncCommandId_IDENTIFY_FIRMWARE:
			res = firmware.IdentifyFirmware()
		case basefunc.BaseFuncCommandId_LOAD_FIRMWARE_CHUNK:
			res, doreset = firmware.LoadFirmwareChunk(c.GetLoadFirmwareChunk())
		default:
			res = &basefunc.BaseFuncResponse{
				Id:     c.Id,
				Status: basefunc.BaseFuncStatus_UNKNOWN_COMMAND,
			}
		}

		err = ch.WriteMessage(res)
		if err != nil {
			log.Printf("Failed to write: %s", err)
			return
		}
		if doreset {
			return
		}
	}
}
