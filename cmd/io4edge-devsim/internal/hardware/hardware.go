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

package hardware

import (
	"github.com/ci4rail/io4edge-client-go/pkg/io4edge/basefunc"
)

type serialNumber struct {
	Hi uint64
	Lo uint64
}

type hardwareID struct {
	RootArticle  string
	SerialNumber serialNumber
	MajorVersion uint32
}

var (
	hwID = hardwareID{"S101-CPU01:UC", serialNumber{0x1234567887654321, 0x4567456745674567}, 1}
)

// IdentifyHardware reports the current hardware inventory data
func IdentifyHardware() *basefunc.BaseFuncResponse {

	res := &basefunc.BaseFuncResponse{
		Id:     basefunc.BaseFuncCommandId_IDENTIFY_HARDWARE,
		Status: basefunc.BaseFuncStatus_OK,
		Data: &basefunc.BaseFuncResponse_IdentifyHardware{
			IdentifyHardware: &basefunc.ResIdentifyHardware{
				RootArticle: hwID.RootArticle,
				SerialNumber: &basefunc.SerialNumber{
					Hi: hwID.SerialNumber.Hi,
					Lo: hwID.SerialNumber.Lo,
				},
				MajorVersion: hwID.MajorVersion,
			},
		},
	}
	return res
}
