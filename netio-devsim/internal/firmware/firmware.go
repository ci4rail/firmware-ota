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

package firmware

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/ci4rail/firmware-ota/pkg/netio_base"
)

type FirmwareId struct {
	Name  string
	Major uint
	Minor uint
	Patch uint
}

type FirmwareHeader struct {
	Name           string
	Major          uint
	Minor          uint
	Patch          uint
	Firmware_Works bool
}

var (
	FwId            = FirmwareId{"default", 1, 0, 0}
	NextChunkNumber = uint32(0)
	NextFlashOffset = uint32(0)
	Flash           = make([]byte, 200000)
)

func IdentifyFirmware() *netio_base.Response {

	res := &netio_base.Response{
		Id:     netio_base.CommandId_IDENTIFY_FIRMWARE,
		Status: netio_base.Status_OK,
		Data: &netio_base.Response_IdentifyFirmware{
			IdentifyFirmware: &netio_base.ResIdentifyFirmware{
				Name:         FwId.Name,
				MajorVersion: uint32(FwId.Major),
				MinorVersion: uint32(FwId.Minor),
				PatchVersion: uint32(FwId.Patch),
			},
		},
	}
	return res
}

func LoadFirmwareChunk(c *netio_base.CmdLoadFirmwareChunk) *netio_base.Response {

	var status = netio_base.Status_OK

	if NextChunkNumber != c.ChunkNumber {
		status = netio_base.Status_CHUNK_SEQ_ERROR
	} else {
		log.Printf("Loading chunk %d @%08x\n", NextChunkNumber, NextFlashOffset)

		// simulate flash programming
		copy(Flash[NextFlashOffset:], c.Data)
		NextFlashOffset += uint32(len(c.Data))
		NextChunkNumber++

		if c.IsLastChunk {
			NextFlashOffset = 0
			NextChunkNumber = 0
			header, err := fwHeaderFromFlash(Flash)
			if err != nil {
				log.Printf("firmware header not ok %v\n", err)
			} else if !header.Firmware_Works {
				log.Printf("firmware not working\n")
			} else {
				log.Printf("activating new firmware %v (slow)\n", header)
				time.Sleep(4 * time.Second)
				FwId.Name = header.Name
				FwId.Major = header.Major
				FwId.Minor = header.Minor
				FwId.Patch = header.Patch
			}
		}

	}

	res := &netio_base.Response{
		Id:     netio_base.CommandId_LOAD_FIRMWARE_CHUNK,
		Status: status,
	}
	return res
}

func fwHeaderFromFlash(flash []byte) (*FirmwareHeader, error) {
	// find end of json. This works only if json has no nested {}
	idx := strings.Index(string(flash), "}")
	if idx == -1 {
		return nil, errors.New("bad json")
	}
	flash_json := flash[:idx+1]

	var header FirmwareHeader
	err := json.Unmarshal(flash_json, &header)
	if err != nil {
		return nil, err
	}
	return &header, nil
}
