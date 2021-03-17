package printer

import (
	"context"
	"github.com/olekukonko/tablewriter"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"kubers/pkg/client"
	"kubers/pkg/resources"
	"os"
	"sort"
	"sync"
)

func PrintNsMetrics(ns string, order string, by string, label string) {
	mcl, cl := getClients()

	ctx := context.Background()
	options := metaV1.ListOptions{LabelSelector: label}
	if ns != "" {
		options.FieldSelector = "metadata.name=" + ns
	}
	namespace, _ := cl.CoreV1().Namespaces().List(ctx, options)

	items := resources.NsItems{
		SortBy:    by,
		SortOrder: order,
	}
	items.BuildData(namespace.Items)
	items.GetResources(mcl, cl)

	sort.Sort(resources.SortNsBy(items))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Ns", "CPU (Request/Limit) Cores", "Memory (Request/Limit) Gb"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetCenterSeparator("")
	table.SetAutoMergeCells(false)
	table.AppendBulk(items.FormatForPrint())
	table.SetFooter([]string{"", items.GetTotalCpuFormatted(), items.GetTotalMemFormatted()})
	table.Render()
}

func PrintPodsMetrics(ns string, order string, by string, label string, byCo bool) {
	mcl, cl := getClients()

	ctx := context.Background()
	options := metaV1.ListOptions{LabelSelector: label}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	chPml := make(chan v1beta1.PodMetricsList)
	go func(wg *sync.WaitGroup, ch chan v1beta1.PodMetricsList) {
		podMetricsList, _ := mcl.MetricsV1beta1().PodMetricses(ns).List(ctx, options)
		wg.Done()

		ch <- *podMetricsList
	}(wg, chPml)

	chPdl := make(chan v1.PodList)
	go func(wg *sync.WaitGroup, ch chan v1.PodList) {
		podDetailsList, _ := cl.CoreV1().Pods(ns).List(ctx, options)
		wg.Done()

		ch <- *podDetailsList
	}(wg, chPdl)

	wg.Wait()

	podMetricsList := <-chPml
	podDetailsList := <-chPdl

	items := resources.PodItems{
		SortBy:    by,
		SortOrder: order,
	}
	items.BuildData(podMetricsList.Items, podDetailsList.Items)

	sort.Sort(resources.SortPoBy(items))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Pod", "Container(s)", "CPU (Request/Limit) Cores", "Memory (Request/Limit) Gb"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetCenterSeparator("")
	table.SetAutoMergeCells(false)
	table.AppendBulk(items.FormatForPrint(byCo))
	table.SetFooter([]string{"", "", items.GetTotalCpuFormatted(), items.GetTotalMemFormatted()})
	table.Render()
}

func getClients() (*versioned.Clientset, *kubernetes.Clientset) {
	return client.GetMetricsClientOutsideCLuster(), client.GetKubernetesClientOutsideCluster()
}
