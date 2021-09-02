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
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ci4rail/firmware-ota/netio-cli/internal/devproto"
	e "github.com/ci4rail/firmware-ota/netio-cli/internal/errors"
	"github.com/ci4rail/firmware-ota/pkg/netio_base"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var loadFirmwareCmd = &cobra.Command{
	Use:     "load-firmware FW_FILE",
	Aliases: []string{"load"},
	Short:   "Upload firmware to device",
	Long: `Upload firmware to device.
Example:
netio-cli <firmware-file>`,
	Run:  loadFirmware,
	Args: cobra.ExactArgs(1),
}

func loadFirmware(cmd *cobra.Command, args []string) {

	dp, err := devproto.NewDevProtoClientConnection(host)
	e.ErrChk(err)

	c := &netio_base.Command{
		Id: netio_base.CommandId_LOAD_FIRMWARE_CHUNK,
		Data: &netio_base.Command_LoadFirmwareChunk{
			LoadFirmwareChunk: &netio_base.CmdLoadFirmwareChunk{
				Data: make([]byte, 1024),
			},
		},
	}
	data := c.GetLoadFirmwareChunk().Data

	f, err := os.Open(args[0])
	e.ErrChk(err)

	defer f.Close()

	reader := bufio.NewReader(f)
	chunkNumber := uint32(0)

	for {
		at_eof := false

		n, err := reader.Read(data)
		e.ErrChk(err)

		// check if we are at EOF
		_, err = reader.Peek(1)
		if err == io.EOF {
			at_eof = true
		}
		log.Printf("Read %d bytes at_eof=%v chunk %d\n", n, at_eof, chunkNumber)

		c.GetLoadFirmwareChunk().IsLastChunk = at_eof
		c.GetLoadFirmwareChunk().ChunkNumber = chunkNumber

		err = dp.WriteMessage(c)
		e.ErrChk(err)

		r := &netio_base.Response{}

		err = dp.ReadMessage(r)
		e.ErrChk(err)

		fmt.Printf("%v\n", r)

		if at_eof {
			break
		}
		chunkNumber++
	}

}

func init() {
	rootCmd.AddCommand(loadFirmwareCmd)
}
