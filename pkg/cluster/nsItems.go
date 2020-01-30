package cluster

import (
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type NsItems struct {
	data          []NsItem
	totalCpu      int64
	totalMem      int64
	sortBy        string
	sortDirection string
}

type NsItem struct {
	name   string
	labels map[string]string
	rs     NsItemResources
}

type NsItemResources struct {
	mem int64
	cpu int64
}

type SortNsBy NsItems

func (n SortNsBy) Len() int {
	return len(n.data)
}

func (n SortNsBy) Less(i, j int) bool {
	switch n.sortBy {
	case "cpu":
		if n.sortDirection == "asc" {
			return n.data[i].rs.cpu < n.data[j].rs.cpu
		} else {
			return n.data[i].rs.cpu > n.data[j].rs.cpu
		}
	case "mem":
		if n.sortDirection == "asc" {
			return n.data[i].rs.mem < n.data[j].rs.mem
		} else {
			return n.data[i].rs.mem > n.data[j].rs.mem
		}
	default:
		return true
	}
}

func (n SortNsBy) Swap(i, j int) {
	n.data[i], n.data[j] = n.data[j], n.data[i]
}

func (n *NsItems) buildData(itemsList []v1.Namespace) {
	for _, ns := range itemsList {
		n.data = append(n.data, NsItem{name: ns.Name, labels: ns.Labels})
	}
}

func (n *NsItems) getResources(mcl *metrics.Clientset) {
	c := make(chan NsItemResources)
	for _, ns := range n.data {
		go getNsResources(mcl, ns.name, c)
	}

	for i := 0; i < len(n.data); i++ {
		rs := <-c
		n.data[i].rs = rs

		n.totalCpu += rs.cpu
		n.totalMem += rs.mem
	}
}

func getNsResources(mcl *metrics.Clientset, ns string, c chan NsItemResources) {
	podMetricsList, _ := mcl.MetricsV1beta1().PodMetricses(ns).List(metaV1.ListOptions{})

	items := PoItems{}
	items.buildData(podMetricsList.Items)

	rs := NsItemResources{
		cpu: items.getTotalCpu(),
		mem: items.getTotalMem(),
	}

	c <- rs
}

func (n NsItems) getTotalCpuFormatted() string {
	return formatCpu(n.totalCpu)
}

func (n NsItems) getTotalMemFormatted() string {
	return formatMem(n.totalMem)
}

func (n NsItems) formatForPrint() [][]string {
	var itemsForPrint [][]string
	for _, ns := range n.data {
		row := []string{ns.name, ns.rs.getCpuFormatted(), ns.rs.getMemFormatted()}
		itemsForPrint = append(itemsForPrint, row)
	}

	return itemsForPrint
}

func (r NsItemResources) getMemFormatted() string {
	return formatMem(r.mem)
}

func (r NsItemResources) getCpuFormatted() string {
	return formatCpu(r.cpu)
}
