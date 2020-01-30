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
	"kubectl-ext/pkg/cluster"
)

// rsrCmd represents the rsr command
var rspCmd = &cobra.Command{
	Use:   "rsp",
	Short: "List resources for a namespace grouped by pods of containers",
	Long: `You can list resources for a namespace for pods or for containers

This kubectl plugin is meant to facilitate the inspection of the cluster used resources by namespace.
It can list resource usages for each pod or in a more detailed way for each container of a pod.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rsp called")

		ns, _ := cmd.Flags().GetString("ns")
		byPo, _ := cmd.Flags().GetBool("by-pod")
		sort, _ := cmd.Flags().GetString("sort")
		by, _ := cmd.Flags().GetString("by")

		cluster.PrintPodsMetrics(ns, byPo, sort, by)
	},
	PreRun: preRunFlags,
}

var rsnCmd = &cobra.Command{
	Use:   "rsn",
	Short: "List namespace resources",
	Long: `You can list resource usage for namespaces

This kubectl plugin is meant to facilitate the inspection of the cluster used resources by namespace.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rsn called")

		sort, _ := cmd.Flags().GetString("sort")
		by, _ := cmd.Flags().GetString("by")
		filter, _ := cmd.Flags().GetString("filter")

		cluster.PrintNsMetrics(sort, by, filter)
	},
	PreRun: preRunFlags,
}

var rsnoCmd = &cobra.Command{
	Use:   "rsno",
	Short: "List node resources",
	Long: `You can list resource usage for nodes

This kubectl plugin is meant to facilitate the inspection of the cluster used resources by node.`,
	Run: func(cmd *cobra.Command, args []string) {
		sort, _ := cmd.Flags().GetString("sort")
		by, _ := cmd.Flags().GetString("by")
		filter, _ := cmd.Flags().GetString("filter")

		cluster.PrintNsMetrics(sort, by, filter)
	},
	PreRun: preRunFlags,
}

var preRunFlags = func(cmd *cobra.Command, args []string) {
	sort, _ := cmd.Flags().GetString("sort")
	by, _ := cmd.Flags().GetString("by")

	if sort != "" && (sort == "asc" || sort == "desc") && by == "" {
		cmd.Flags().Set("by", "mem")
	}
}

func init() {
	// Po resources command
	rootCmd.AddCommand(rspCmd)

	rspCmd.Flags().StringP("ns", "n", "", "Show resources for namespace.")
	_ = rspCmd.MarkFlagRequired("ns")
	rspCmd.Flags().BoolP("by-pod", "p", false, "Group resources by Pods.")
	rspCmd.Flags().StringP("sort", "s", "", "Order by pod resources.")
	rspCmd.Flags().StringP("by", "b", "", "Order by cpu or memory resources.")

	// Ns resources command
	rootCmd.AddCommand(rsnCmd)

	rsnCmd.Flags().StringP("sort", "s", "", "Order by resources used.")
	rsnCmd.Flags().StringP("by", "b", "", "Order by cpu or memory resources.")
	rsnCmd.Flags().StringP("filter", "f", "", "Filter by label.")

	// Node resources commmand
	rootCmd.AddCommand(rsnoCmd)

	rsnoCmd.Flags().StringP("sort", "s", "", "Order by resources used.")
	rsnoCmd.Flags().StringP("by", "b", "", "Order by cpu or memory resources.")
	rsnoCmd.Flags().StringP("filter", "f", "", "Filter by label.")
}
