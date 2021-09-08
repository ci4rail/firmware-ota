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

	"github.com/ci4rail/firmware-ota/cmd/io4edge-cli/internal/client"
	e "github.com/ci4rail/firmware-ota/cmd/io4edge-cli/internal/errors"

	"github.com/spf13/cobra"
)

var (
	chunkSize = uint(1024)
)

var loadFirmwareCmd = &cobra.Command{
	Use:     "load-firmware FW_PKG",
	Aliases: []string{"load"},
	Short:   "Upload firmware package to device",
	Long: `Upload firmware package to device.
Example:
io4edge-cli load-firmware <firmware-package-file>`,
	Run:  loadFirmware,
	Args: cobra.ExactArgs(1),
}

func loadFirmware(cmd *cobra.Command, args []string) {
	// file := args[0]
	// c, err := client.NewClient(device)
	// e.ErrChk(err)

}

var loadRawFirmwareCmd = &cobra.Command{
	Use:     "load-raw-firmware FW_FILE",
	Aliases: []string{"load-raw"},
	Short:   "Upload firmware binary file to device",
	Long: `Upload firmware binary file to device.
NO COMPATIBILITY CHECKS! USE AT YOUR OWN RISK!
Example:
io4edge-cli load-raw-firmware <firmware-file>`,
	Run:  loadRawFirmware,
	Args: cobra.ExactArgs(1),
}

func loadRawFirmware(cmd *cobra.Command, args []string) {
	file := args[0]
	c, err := client.NewClient(device)
	e.ErrChk(err)

	err = c.LoadFirmwareBinaryFromFile(file, chunkSize, time.Duration(timeoutSecs)*time.Second)
	e.ErrChk(err)

	fmt.Println("Success. Read back programmed firmware info.")
	identifyFirmware(cmd, args)
}

func init() {
	rootCmd.AddCommand(loadFirmwareCmd)
	rootCmd.AddCommand(loadRawFirmwareCmd)
}
