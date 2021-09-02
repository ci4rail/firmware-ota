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

	"github.com/ci4rail/firmware-ota/netio-devsim/internal/devproto"
	"github.com/ci4rail/firmware-ota/netio-devsim/internal/firmware"
	"github.com/ci4rail/firmware-ota/netio-devsim/pkg/version"
	"github.com/ci4rail/firmware-ota/pkg/netio_base"
)

var (
	port = ":9999"
)

func main() {
	log.Printf("netio-devsim version: %s listen at port %s\n", version.Version, port)

	dp, err := devproto.NewDevProto(port)
	if err != nil {
		log.Fatalf("Failed to create devproto: %s", err)
	}
	for {
		err := dp.WaitForConnection()
		if err != nil {
			log.Fatalf("Failed to wait for connection: %s", err)
		}
		serveConnection(dp)
	}
}

func serveConnection(dp *devproto.DevProto) {
	defer dp.Close()

	for {
		c := &netio_base.Command{}
		err := dp.ReadMessage(c)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Fatalf("Failed to read: %s", err)
		}

		var res *netio_base.Response
		switch c.Id {
		case netio_base.CommandId_IDENTIFY_FIRMWARE:
			res = firmware.IdentifyFirmware()
		case netio_base.CommandId_LOAD_FIRMWARE_CHUNK:
			res = firmware.LoadFirmwareChunk(c.GetLoadFirmwareChunk())
		default:
			res = &netio_base.Response{
				Id:     c.Id,
				Status: netio_base.Status_UNKNOWN_COMMAND,
			}
		}

		err = dp.WriteMessage(res)
		if err != nil {
			log.Fatalf("Failed to write: %s", err)
		}
	}
}
