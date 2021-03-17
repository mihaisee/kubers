package resources

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	"kubers/pkg/formatter"
	"sync"
)

type NsItems struct {
	data []Namespace

	SortBy    string
	SortOrder string
}

type SortNsBy NsItems

func (n SortNsBy) Len() int {
	return len(n.data)
}

func (n SortNsBy) Less(i, j int) bool {
	switch n.SortBy {
	case "cpu":
		if n.SortOrder == "asc" {
			return n.data[i].totalCpuUsage < n.data[j].totalCpuUsage
		} else {
			return n.data[i].totalCpuUsage > n.data[j].totalCpuUsage
		}
	case "mem":
		if n.SortOrder == "asc" {
			return n.data[i].totalMemUsage < n.data[j].totalMemUsage
		} else {
			return n.data[i].totalMemUsage > n.data[j].totalMemUsage
		}
	default:
		return true
	}
}

func (n SortNsBy) Swap(i, j int) {
	n.data[i], n.data[j] = n.data[j], n.data[i]
}

func (items *NsItems) BuildData(itemsList []v1.Namespace) {
	for _, ns := range itemsList {
		items.data = append(items.data, Namespace{name: ns.Name, labels: ns.Labels})
	}
}

func (items *NsItems) GetResources(mcl *metrics.Clientset, cl *kubernetes.Clientset) {
	c := make(chan Namespace)
	for _, ns := range items.data {
		go getNsResources(mcl, cl, ns.name, ns.labels, c)
	}

	for i := 0; i < len(items.data); i++ {
		nsItem := <-c
		items.data[i] = nsItem
	}
}

func getNsResources(mcl *metrics.Clientset, cl *kubernetes.Clientset, ns string, labels map[string]string, c chan Namespace) {
	ctx := context.Background()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	chPml := make(chan v1beta1.PodMetricsList)
	go func(wg *sync.WaitGroup, ch chan v1beta1.PodMetricsList) {
		podMetricsList, _ := mcl.MetricsV1beta1().PodMetricses(ns).List(ctx, metaV1.ListOptions{})
		wg.Done()

		ch <- *podMetricsList
	}(wg, chPml)

	chPl := make(chan v1.PodList)
	go func(wg *sync.WaitGroup, ch chan v1.PodList) {
		podDetails, _ := cl.CoreV1().Pods(ns).List(ctx, metaV1.ListOptions{})
		wg.Done()

		ch <- *podDetails
	}(wg, chPl)

	wg.Wait()

	podMetricsList := <-chPml
	podDetails := <-chPl

	podItems := &PodItems{}
	if len(podMetricsList.Items) > 0 {
		podItems.BuildData(podMetricsList.Items, podDetails.Items)
	}

	nsItem := Namespace{
		name:     ns,
		labels:   labels,
		podItems: podItems,
	}

	nsItem.calculateTotalCpuUsage()
	nsItem.calculateTotalMemUsage()
	nsItem.calculateTotalCpuSpec()
	nsItem.calculateTotalMemSpec()

	c <- nsItem
}

func (items *NsItems) getTotalCpuUsage() int64 {
	var cpu int64
	for _, ns := range items.data {
		cpu += ns.totalCpuUsage
	}

	return cpu
}

func (items *NsItems) getTotalCpuRequest() int64 {
	var cpu int64
	for _, ns := range items.data {
		cpu += ns.totalCpuRequest
	}

	return cpu
}

func (items *NsItems) getTotalCpuLimit() int64 {
	var cpu int64
	for _, ns := range items.data {
		cpu += ns.totalCpuLimit
	}

	return cpu
}

func (items *NsItems) getTotalMemUsage() int64 {
	var mem int64
	for _, ns := range items.data {
		mem += ns.totalMemUsage
	}

	return mem
}

func (items *NsItems) getTotalMemRequest() int64 {
	var mem int64
	for _, ns := range items.data {
		mem += ns.totalMemRequest
	}

	return mem
}

func (items *NsItems) getTotalMemLimit() int64 {
	var mem int64
	for _, ns := range items.data {
		mem += ns.totalMemLimit
	}

	return mem
}

func (items *NsItems) GetTotalCpuFormatted() string {
	return formatter.FormatCpuUsageWithSpec(items.getTotalCpuUsage(), items.getTotalCpuRequest(), items.getTotalCpuLimit())
}

func (items *NsItems) GetTotalMemFormatted() string {
	return formatter.FormatMemUsageWithSpec(items.getTotalMemUsage(), items.getTotalMemRequest(), items.getTotalMemLimit())
}

func (items *NsItems) FormatForPrint() [][]string {
	var itemsForPrint [][]string
	for _, ns := range items.data {
		row := []string{
			ns.name,
			ns.GetTotalCpuFormatted(),
			ns.GetTotalMemFormatted(),
		}
		itemsForPrint = append(itemsForPrint, row)
	}

	return itemsForPrint
}

