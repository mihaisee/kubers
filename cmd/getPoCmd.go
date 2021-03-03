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

var getPoCmd = &cobra.Command{
	Use:   "po",
	Short: "List resources for a namespace grouped by pods of containers",
	Long: `You can list resources for a namespace for pods or for containers

This kubectl plugin is meant to facilitate the inspection of the printer used resources by namespace.
It can list resource usages for each pod or in a more detailed way for each container of a pod.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rsp called")

		ns, _ := cmd.Flags().GetString("ns")
		byPo, _ := cmd.Flags().GetBool("by-pod")
		order, _ := cmd.Flags().GetString("order")
		by, _ := cmd.Flags().GetString("by")

		printer.PrintPodsMetrics(ns, byPo, order, by)
	},
	PreRun: getPoCmdPreRun,
}

var getPoCmdPreRun = func(getPoCmd *cobra.Command, args []string) {
	preRunFlags(getPoCmd)
}

func init() {
	getCmd.AddCommand(getPoCmd)

	getPoCmd.Flags().StringP("order", "o", "", "Order by resources used.")
	getPoCmd.Flags().StringP("by", "b", "", "Order by cpu or memory resources.")
	getPoCmd.Flags().StringP("filter", "f", "", "Filter by label.")
}
