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
	"time"

	"github.com/ci4rail/firmware-ota/cmd/netio-cli/internal/client"
	e "github.com/ci4rail/firmware-ota/cmd/netio-cli/internal/errors"

	"github.com/spf13/cobra"
)

var (
	chunkSize = uint(1024)
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
	file := args[0]
	c, err := client.NewClient(device)
	e.ErrChk(err)

	err = c.LoadFirmwareFromFile(file, chunkSize, time.Duration(timeoutSecs)*time.Second)
	e.ErrChk(err)

	fmt.Printf("New ")
	identifyFirmware(cmd, args)
}

func init() {
	rootCmd.AddCommand(loadFirmwareCmd)
}
