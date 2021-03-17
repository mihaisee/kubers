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

var getNsCmd = &cobra.Command{
	Use:   "ns",
	Short: "List namespaces with resources",
	Long: `Usage: 
kubers get ns
kubers get ns -n staging
kubers get ns -b mem -o asc
kubers get ns -l selector=team-ns`,
	Run: func(cmd *cobra.Command, args []string) {
		ns, _ := cmd.Flags().GetString("ns")
		order, _ := cmd.Flags().GetString("order")
		by, _ := cmd.Flags().GetString("by")
		label, _ := cmd.Flags().GetString("label")

		printer.PrintNsMetrics(ns, order, by, label)
	},
	PreRun: getNsCmdPreRun,
}

var getNsCmdPreRun = func(getNsCmd *cobra.Command, args []string) {
	preRunFlags(getNsCmd)
}

func init() {
	getCmd.AddCommand(getNsCmd)

	getNsCmd.Flags().StringP("ns", "n", "", "Filter by namespace.")
	getNsCmd.Flags().StringP("by", "b", "cpu", "Sort by cpu or memory [cpu|mem].")
	getNsCmd.Flags().StringP("order", "o", "desc", "Order to sort by [asc|desc].")
	getNsCmd.Flags().StringP("label", "l", "", "Filter by label.")
}
