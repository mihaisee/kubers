package cluster

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"os"
	sortUtility "sort"
	"strconv"
)

func PrintPodsMetrics(namespace string, byPo bool, sort string, by string) {
	mcl := getMetricsClientOutsideCLuster()

	podMetricsList, _ := mcl.MetricsV1beta1().PodMetricses(namespace).List(metaV1.ListOptions{})

	items := PoItems{}
	items.buildData(podMetricsList.Items)

	items.sortBy = by
	items.sortDirection = sort

	sortUtility.Sort(SortPoBy(items))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Pod", "Container(s)", "CPU", "Memory"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoMergeCells(false)
	table.SetRowLine(true)
	table.AppendBulk(items.formatForPrint(byPo))
	table.SetFooter([]string{"", "", items.getTotalCpuFormatted(), items.getTotalMemFormatted()})
	table.Render()
}

func PrintNsMetrics(sort string, by string, filter string) {
	cl := getKubernetesClientOutsideCluster()

	options := metaV1.ListOptions{LabelSelector: filter}
	ns, _ := cl.CoreV1().Namespaces().List(options)

	items := NsItems{}
	items.sortBy = by
	items.sortDirection = sort
	items.buildData(ns.Items)

	mcl := getMetricsClientOutsideCLuster()
	items.getResources(mcl)
	sortUtility.Sort(SortNsBy(items))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Ns", "CPU", "Memory"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoMergeCells(false)
	table.SetRowLine(true)
	table.AppendBulk(items.formatForPrint())
	table.SetFooter([]string{"", items.getTotalCpuFormatted(), items.getTotalMemFormatted()})
	table.Render()
}

func PrintNoMetrics(sort string, by string, filter string) {
	mcl := getMetricsClientOutsideCLuster()

	options := metaV1.ListOptions{}
	noMetricsList, _ := mcl.MetricsV1beta1().NodeMetricses().List(options)

	fmt.Println(noMetricsList)
}

func formatMem(mem int64) string {
	memGb := float64(mem) * 0.000000001
	return strconv.FormatFloat(memGb, 'f', 3, 64) + "Gb"
}

func formatCpu(cpu int64) string {
	return strconv.FormatInt(cpu, 10)
}
