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

	e "github.com/ci4rail/firmware-ota/cmd/netio-cli/internal/errors"
	"github.com/ci4rail/firmware-ota/pkg/netio"
	"github.com/ci4rail/firmware-ota/pkg/netio/transport"
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

	conn, err := transport.NewConnection(host)
	e.ErrChk(err)

	ch, err := netio.NewChannel(conn)
	e.ErrChk(err)

	c := &netio.BaseFuncCommand{
		Id: netio.BaseFuncCommandId_IDENTIFY_FIRMWARE,
	}

	err = ch.WriteMessage(c)
	e.ErrChk(err)

	r := &netio.BaseFuncResponse{}

	err = ch.ReadMessage(r)
	e.ErrChk(err)

	fmt.Printf("%v\n", r)

}

func init() {
	rootCmd.AddCommand(identifyFirmwareCmd)
}
