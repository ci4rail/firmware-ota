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

	"github.com/ci4rail/firmware-ota/pkg/netio"
)

type firmwareID struct {
	Name  string
	Major uint
	Minor uint
	Patch uint
}

type firmwareHeader struct {
	Name          string
	Major         uint
	Minor         uint
	Patch         uint
	FirmwareWorks bool `json:"firmware_works"`
}

var (
	fwID            = firmwareID{"default", 1, 0, 0}
	nextChunkNumber = uint32(0)
	nextFlashOffset = uint32(0)
	flash           = make([]byte, 200000)
)

// IdentifyFirmware reports the currently active firmware name and version
func IdentifyFirmware() *netio.BaseFuncResponse {

	res := &netio.BaseFuncResponse{
		Id:     netio.BaseFuncCommandId_IDENTIFY_FIRMWARE,
		Status: netio.BaseFuncStatus_OK,
		Data: &netio.BaseFuncResponse_IdentifyFirmware{
			IdentifyFirmware: &netio.ResIdentifyFirmware{
				Name:         fwID.Name,
				MajorVersion: uint32(fwID.Major),
				MinorVersion: uint32(fwID.Minor),
				PatchVersion: uint32(fwID.Patch),
			},
		},
	}
	return res
}

// LoadFirmwareChunk loads the chunk in c to the virtual flash
func LoadFirmwareChunk(c *netio.CmdLoadFirmwareChunk) *netio.BaseFuncResponse {

	var status = netio.BaseFuncStatus_OK

	if nextChunkNumber != c.ChunkNumber {
		status = netio.BaseFuncStatus_CHUNK_SEQ_ERROR
	} else {
		log.Printf("Loading chunk %d @%08x\n", nextChunkNumber, nextFlashOffset)

		// simulate flash programming
		copy(flash[nextFlashOffset:], c.Data)
		nextFlashOffset += uint32(len(c.Data))
		nextChunkNumber++

		if c.IsLastChunk {
			nextFlashOffset = 0
			nextChunkNumber = 0
			header, err := fwHeaderFromFlash(flash)
			if err != nil {
				log.Printf("firmware header not ok %v\n", err)
			} else if !header.FirmwareWorks {
				log.Printf("firmware not working\n")
			} else {
				log.Printf("activating new firmware %v (slow)\n", header)
				time.Sleep(4 * time.Second)
				fwID.Name = header.Name
				fwID.Major = header.Major
				fwID.Minor = header.Minor
				fwID.Patch = header.Patch
			}
		}

	}

	res := &netio.BaseFuncResponse{
		Id:     netio.BaseFuncCommandId_LOAD_FIRMWARE_CHUNK,
		Status: status,
	}
	return res
}

func fwHeaderFromFlash(flash []byte) (*firmwareHeader, error) {
	// find end of json. This works only if json has no nested {}
	idx := strings.Index(string(flash), "}")
	if idx == -1 {
		return nil, errors.New("bad json")
	}
	flashJSON := flash[:idx+1]

	var header firmwareHeader
	err := json.Unmarshal(flashJSON, &header)
	if err != nil {
		return nil, err
	}
	return &header, nil
}
