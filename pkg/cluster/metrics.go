package cluster

import (
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"os"
	sortUtility "sort"
	"strconv"
)

func PrintPodsMetrics(ns string, byPo bool, sort string, by string) {
	mcl := getMetricsClientOutsideCLuster()
	cl := getKubernetesClientOutsideCluster()

	ctx := context.Background()
	podMetricsList, _ := mcl.MetricsV1beta1().PodMetricses(ns).List(ctx, metaV1.ListOptions{})
	podDetails, _ := cl.CoreV1().Pods(ns).List(ctx, metaV1.ListOptions{})

	items := PoItems{}
	items.buildData(podMetricsList.Items, podDetails.Items)

	items.sortBy = by
	items.sortDirection = sort

	sortUtility.Sort(SortPoBy(items))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Pod", "Container(s)", "CPU", "Memory"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetCenterSeparator("")
	table.SetAutoMergeCells(false)
	table.AppendBulk(items.formatForPrint(byPo))
	table.SetFooter([]string{"", "", items.getTotalCpuFormatted(), items.getTotalMemFormatted()})
	table.Render()
}

func PrintNsMetrics(sort string, by string, filter string) {
	cl := getKubernetesClientOutsideCluster()

	ctx := context.Background()
	options := metaV1.ListOptions{LabelSelector: filter}
	ns, _ := cl.CoreV1().Namespaces().List(ctx, options)

	items := NsItems{}
	items.sortBy = by
	items.sortDirection = sort
	items.buildData(ns.Items)

	mcl := getMetricsClientOutsideCLuster()
	items.getResources(mcl, cl)
	sortUtility.Sort(SortNsBy(items))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Ns", "CPU", "Memory"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetCenterSeparator("")
	table.SetAutoMergeCells(false)
	table.AppendBulk(items.formatForPrint())
	table.SetFooter([]string{"", items.getTotalCpuFormatted(), items.getTotalMemFormatted()})
	table.Render()
}

func PrintNoMetrics(sort string, by string, filter string) {
	mcl := getMetricsClientOutsideCLuster()

	ctx := context.Background()
	options := metaV1.ListOptions{}
	noMetricsList, _ := mcl.MetricsV1beta1().NodeMetricses().List(ctx, options)

	fmt.Println(noMetricsList)
}

func formatMem(mem int64) string {
	memGb := float64(mem) * 0.000000001
	return strconv.FormatFloat(memGb, 'f', 3, 64) + "Gb"
}

func formatCpu(cpu int64) string {
	return strconv.FormatInt(cpu, 10)
}
