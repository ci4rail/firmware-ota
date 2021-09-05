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

	e "github.com/ci4rail/firmware-ota/cmd/netio-cli/internal/errors"
	"github.com/ci4rail/firmware-ota/pkg/netio"
	"github.com/ci4rail/firmware-ota/pkg/netio/transport"
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

	conn, err := transport.NewConnection(host)
	e.ErrChk(err)

	ch, err := netio.NewChannel(conn)

	c := &netio.BaseFuncCommand{
		Id: netio.BaseFuncCommandId_LOAD_FIRMWARE_CHUNK,
		Data: &netio.BaseFuncCommand_LoadFirmwareChunk{
			LoadFirmwareChunk: &netio.CmdLoadFirmwareChunk{
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
		atEOF := false

		n, err := reader.Read(data)
		e.ErrChk(err)

		// check if we are at EOF
		_, err = reader.Peek(1)
		if err == io.EOF {
			atEOF = true
		}
		log.Printf("Read %d bytes at_eof=%v chunk %d\n", n, atEOF, chunkNumber)

		c.GetLoadFirmwareChunk().IsLastChunk = atEOF
		c.GetLoadFirmwareChunk().ChunkNumber = chunkNumber

		err = ch.WriteMessage(c)
		e.ErrChk(err)

		r := &netio.BaseFuncResponse{}

		err = ch.ReadMessage(r)
		e.ErrChk(err)

		fmt.Printf("%v\n", r)

		if atEOF {
			break
		}
		chunkNumber++
	}

}

func init() {
	rootCmd.AddCommand(loadFirmwareCmd)
}
