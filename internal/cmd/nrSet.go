/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var nrSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set NR",
	Long:  `A command to set NR.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("NR set called,\ngNB struct value is: \n%+v\n", nr)
	},
}

func init() {
	nrCmd.AddCommand(nrSetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	nrSetCmd.Flags().IntVarP(&nr.gnbId, "id", "i", 1, "Id of NR")
	nrSetCmd.Flags().IntVar(&nr.mcc, "mcc", 208, "MCC")
	nrSetCmd.Flags().IntVar(&nr.mnc, "mnc", 93, "MNC")
	nrSetCmd.Flags().IntVar(&nr.tac, "tac", 1, "TAC")
	nrSetCmd.Flags().IntVar(&nr.sst, "sst", 1, "SST")
	nrSetCmd.Flags().IntVar(&nr.sd, "sd", 1, "SD")

}
