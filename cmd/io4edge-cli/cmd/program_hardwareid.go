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
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var programHardwareIdentificationCmd = &cobra.Command{
	Use:     "program-hwid NAME MAJOR SERIAL",
	Aliases: []string{"hwid"},
	Short:   "Upload firmware to device",
	Long: `Upload firmware to device.
Example:
io4edge-cli program-hwid S101-IOU04 1 6ba7b810-9dad-11d1-80b4-00c04fd430c8`,
	Run:  programHardwareIdentification,
	Args: cobra.ExactArgs(1),
}

func programHardwareIdentification(cmd *cobra.Command, args []string) {
	//file := args[0]
	//c, err := client.NewClient(device)
	//e.ErrChk(err)

	//err = c.ProgramHardwareIdentification()
}

func init() {
	rootCmd.AddCommand(programHardwareIdentificationCmd)
}
