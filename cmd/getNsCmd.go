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
	"fmt"
	"github.com/spf13/cobra"
	"kubers/pkg/printer"
)

var getNsCmd = &cobra.Command{
	Use:   "ns",
	Short: "List namespaces with resources",
	Long: `List resource usage for namespaces. By default ordered 'desc' by 'cpu'`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("kubers ns called")

		order, _ := cmd.Flags().GetString("order")
		by, _ := cmd.Flags().GetString("by")
		filter, _ := cmd.Flags().GetString("filter")

		printer.PrintNsMetrics(order, by, filter)
	},
	PreRun: getNsCmdPreRun,
}

var getNsCmdPreRun = func(getNsCmd *cobra.Command, args []string) {
	preRunFlags(getNsCmd)
}

func init() {
	getCmd.AddCommand(getNsCmd)

	getNsCmd.Flags().StringP("order", "o", "", "Order to sort by.")
	getNsCmd.Flags().StringP("by", "b", "", "Sort by cpu or memory.")
	getNsCmd.Flags().StringP("filter", "f", "", "Filter by label.")
}
