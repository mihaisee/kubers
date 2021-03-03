/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	"kubers/pkg/printer"
)

var getNoCmd = &cobra.Command{
	Use:   "no",
	Short: "List node resources",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

		order, _ := cmd.Flags().GetString("order")
		by, _ := cmd.Flags().GetString("by")
		filter, _ := cmd.Flags().GetString("filter")

		printer.PrintNoMetrics(order, by, filter)
	},
	PreRun: getNoCmdPreRun,
}

var getNoCmdPreRun = func(getNoCmd *cobra.Command, args []string) {
	preRunFlags(getNoCmd)
}

func init() {
	getCmd.AddCommand(getNoCmd)

	getNoCmd.Flags().StringP("order", "o", "", "Order by resources used.")
	getNoCmd.Flags().StringP("by", "b", "", "Order by cpu or memory resources.")
	getNoCmd.Flags().StringP("filter", "f", "", "Filter by label.")
}
