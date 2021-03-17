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

var getPoCmd = &cobra.Command{
	Use:   "po",
	Short: "List pods with containers and resources.",
	Long:  `Usage:
kubers get po
kubers get po -n staging
kubers get po -b mem -o asc
kubers get po -l app=www
kubers get po -c`,
	Run: func(cmd *cobra.Command, args []string) {
		ns, _ := cmd.Flags().GetString("ns")
		order, _ := cmd.Flags().GetString("order")
		by, _ := cmd.Flags().GetString("by")
		label, _ := cmd.Flags().GetString("label")
		byContainer, _ := cmd.Flags().GetBool("by-container")

		printer.PrintPodsMetrics(ns, order, by, label, byContainer)
	},
	PreRun: getPoCmdPreRun,
}

var getPoCmdPreRun = func(getPoCmd *cobra.Command, args []string) {
	preRunFlags(getPoCmd)
}

func init() {
	getCmd.AddCommand(getPoCmd)

	getPoCmd.Flags().StringP("ns", "n", "", "Filter by namespace.")
	getPoCmd.Flags().StringP("by", "b", "cpu", "Order by cpu or memory [cpu|mem].")
	getPoCmd.Flags().StringP("order", "o", "desc", "Order to sort by [asc|desc].")
	getPoCmd.Flags().StringP("label", "l", "", "Filter by label.")
	getPoCmd.Flags().BoolP("by-container", "c", false, "List by container or grouped.")
}
