/*
Copyright © 2021 Ci4Rail GmbH <engineering@ci4rail.com>

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

package cmd

import (
	"fmt"

	"github.com/ci4rail/firmware-ota/netio-cli/internal/devproto"
	e "github.com/ci4rail/firmware-ota/netio-cli/internal/errors"
	"github.com/ci4rail/firmware-ota/pkg/netio_base"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var identifyFirmwareCmd = &cobra.Command{
	Use:     "identify-firmware",
	Aliases: []string{"id-fw", "fw"},
	Short:   "Get firmware infos from device",
	Run:     identifyFirmware,
}

func identifyFirmware(cmd *cobra.Command, args []string) {
	fmt.Println("Identify FW called")

	dp, err := devproto.NewDevProtoClientConnection(host)
	e.ErrChk(err)

	c := &netio_base.Command{
		Id: netio_base.CommandId_IDENTIFY_FIRMWARE,
	}

	err = dp.WriteMessage(c)
	e.ErrChk(err)

	r := &netio_base.Response{}

	err = dp.ReadMessage(r)
	e.ErrChk(err)

	fmt.Printf("%v\n", r)

}

func init() {
	rootCmd.AddCommand(identifyFirmwareCmd)
}
