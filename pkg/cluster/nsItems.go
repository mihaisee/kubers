package cluster

import (
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type NsItems struct {
	data          []NsItem
	sortBy        string
	sortDirection string

	totalCpuUsage   int64
	totalCpuRequest int64
	totalCpuLimit   int64

	totalMemUsage   int64
	totalMemRequest int64
	totalMemLimit   int64
}

type NsItem struct {
	name   string
	labels map[string]string
	rs     NsItemResources
}

type NsItemResources struct {
	memUsage   int64
	memRequest int64
	memLimit   int64

	cpuUsage   int64
	cpuRequest int64
	cpuLimit   int64
}

type SortNsBy NsItems

func (n SortNsBy) Len() int {
	return len(n.data)
}

func (n SortNsBy) Less(i, j int) bool {
	switch n.sortBy {
	case "cpu":
		if n.sortDirection == "asc" {
			return n.data[i].rs.cpuUsage < n.data[j].rs.cpuUsage
		} else {
			return n.data[i].rs.cpuUsage > n.data[j].rs.cpuUsage
		}
	case "memUsage":
		if n.sortDirection == "asc" {
			return n.data[i].rs.memUsage < n.data[j].rs.memUsage
		} else {
			return n.data[i].rs.memUsage > n.data[j].rs.memUsage
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

func (n *NsItems) getResources(mcl *metrics.Clientset, cl *kubernetes.Clientset) {
	c := make(chan NsItem)
	for _, ns := range n.data {
		go getNsResources(mcl, cl, ns.name, ns.labels, c)
	}

	for i := 0; i < len(n.data); i++ {
		nsItem := <-c
		n.data[i] = nsItem

		n.totalCpuUsage += n.data[i].rs.cpuUsage
		n.totalMemUsage += n.data[i].rs.memUsage
	}
}

func getNsResources(mcl *metrics.Clientset, cl *kubernetes.Clientset, ns string, labels map[string]string, c chan NsItem) {
	podMetricsList, _ := mcl.MetricsV1beta1().PodMetricses(ns).List(metaV1.ListOptions{})
	podDetails, _ := cl.CoreV1().Pods(ns).List(metaV1.ListOptions{})

	items := PoItems{}
	if len(podMetricsList.Items) > 0 {
		items.buildData(podMetricsList.Items, podDetails.Items)
	}

	nsItem := NsItem{
		name:   ns,
		labels: labels,
		rs: NsItemResources{
			cpuUsage: items.getTotalCpu(),
			memUsage: items.getTotalMem(),
		},
	}

	c <- nsItem
}

func (n NsItems) getTotalCpuFormatted() string {
	return formatCpu(n.totalCpuUsage)
}

func (n NsItems) getTotalMemFormatted() string {
	return formatMem(n.totalMemUsage)
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
	return formatMem(r.memUsage)
}

func (r NsItemResources) getCpuFormatted() string {
	return formatCpu(r.cpuUsage)
}
