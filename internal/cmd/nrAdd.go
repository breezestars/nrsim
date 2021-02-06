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

type Nr struct {
	gnbId int
	mcc   int
	mnc   int
	tac   int
	sst   int
	sd    int
}

// nrAddCmd represents the add command
var nrAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add NR",
	Long:  `A command to add NR.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("NR add called,\ngNB struct value is: \n%+v\n", nr)

	},
}

func init() {
	nrCmd.AddCommand(nrAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	nrAddCmd.Flags().IntVarP(&nr.gnbId, "id", "i", 1, "Id of NR")
	nrAddCmd.Flags().IntVar(&nr.mcc, "mcc", 208, "MCC")
	nrAddCmd.Flags().IntVar(&nr.mnc, "mnc", 93, "MNC")
	nrAddCmd.Flags().IntVar(&nr.tac, "tac", 1, "TAC")
	nrAddCmd.Flags().IntVar(&nr.sst, "sst", 1, "SST")
	nrAddCmd.Flags().IntVar(&nr.sd, "sd", 1, "SD")

}
