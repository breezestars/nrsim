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
	"context"
	"fmt"
	"github.com/cmingou/nrsim/internal/api"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"

	"github.com/spf13/cobra"
)

// nrGetCmd represents the get command
var nrGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get NR",
	Long:  `A command to get NR.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("NR get called,\ngNB struct value is: \n%+v\n", nr)

		ctx, cancel := context.WithTimeout(context.Background(), GrpcConnectTimeout)
		defer cancel()

		var gnbCfgList *api.GnbConfigList
		client := GetCliServerClient()
		gnbCfgList, err := client.ListGnb(ctx, &emptypb.Empty{})
		if err != nil {
			dealError(err)
		}

		if nr.gnbId == -1 {
			// Print all NR
			for _, v := range gnbCfgList.GnbConfig {
				log.Printf("%+v\n", v.GlobalGNBID)
			}
		} else {
			// Print specific NR
			for _, v := range gnbCfgList.GnbConfig {
				if v.GlobalGNBID.Gnbid == uint32(nr.gnbId) {
					log.Printf("%+v\n", v.GlobalGNBID)
					break
				}
			}
		}
	},
}

func init() {
	nrCmd.AddCommand(nrGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	nrGetCmd.Flags().IntVarP(&nr.gnbId, "id", "i", -1, "Id of NR")
}
