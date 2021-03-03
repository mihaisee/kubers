package printer

import (
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"kubers/pkg/client"
	"kubers/pkg/resources"
	"os"
	sortUtility "sort"
)

func PrintNsMetrics(order string, by string, filter string) {
	cl := client.GetKubernetesClientOutsideCluster()

	ctx := context.Background()
	options := metaV1.ListOptions{LabelSelector: filter}
	ns, _ := cl.CoreV1().Namespaces().List(ctx, options)

	items := resources.NsItems{}
	items.SortBy = by
	items.SortDirection = order
	items.BuildData(ns.Items)

	mcl := client.GetMetricsClientOutsideCLuster()
	items.GetResources(mcl, cl)
	sortUtility.Sort(resources.SortNsBy(items))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Ns", "CPU[Cores]", "Memory[Gb]"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetCenterSeparator("")
	table.SetAutoMergeCells(false)
	table.AppendBulk(items.FormatForPrint())
	table.SetFooter([]string{"", items.GetTotalCpuFormatted(), items.GetTotalMemFormatted()})
	table.Render()
}

func PrintPodsMetrics(ns string, byPo bool, order string, by string) {
	mcl := client.GetMetricsClientOutsideCLuster()
	cl := client.GetKubernetesClientOutsideCluster()

	ctx := context.Background()
	podMetricsList, _ := mcl.MetricsV1beta1().PodMetricses(ns).List(ctx, metaV1.ListOptions{})
	podDetails, _ := cl.CoreV1().Pods(ns).List(ctx, metaV1.ListOptions{})

	items := resources.PoItems{}
	items.BuildData(podMetricsList.Items, podDetails.Items)

	items.SortBy = by
	items.SortDirection = order

	sortUtility.Sort(resources.SortPoBy(items))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Pod", "Container(s)", "CPU", "Memory"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetCenterSeparator("")
	table.SetAutoMergeCells(false)
	table.AppendBulk(items.FormatForPrint(byPo))
	table.SetFooter([]string{"", "", items.GetTotalCpuFormatted(), items.GetTotalMemFormatted()})
	table.Render()
}

func PrintNoMetrics(order string, by string, filter string) {
	mcl := client.GetMetricsClientOutsideCLuster()

	ctx := context.Background()
	options := metaV1.ListOptions{}
	noMetricsList, _ := mcl.MetricsV1beta1().NodeMetricses().List(ctx, options)

	fmt.Println(noMetricsList)
}
