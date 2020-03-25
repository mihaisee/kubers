package cluster

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

type PoDetailsItems struct {
	data     map[string]map[string]PoResourcesSpec
	sortData []SortItem

	dataForPrint [][]string

	sortDirection string
	sortBy        string
}

type PoResourcesSpec struct {
	cpuRequest int64
	cpuLimit   int64

	memRequest int64
	memLimit   int64
}

func (items *PoDetailsItems) buildData(itemsList []v1.Pod) {
	items.data = make(map[string]map[string]PoResourcesSpec)
	items.sortData = make([]SortItem, len(itemsList))
	for ind, i := range itemsList {
		items.data[i.Name] = map[string]PoItemResources{}
		for _, c := range i.Containers {
			memUsage, _ := c.Usage.Memory().AsInt64()
			items.data[i.Name][c.Name] = PoItemResources{cpuUsage: c.Usage.Cpu().MilliValue(), memUsage: memUsage}
		}

		items.sortData[ind].resources.memUsage = items.getTotalMemByPod(i.Name)
		items.sortData[ind].resources.cpuUsage = items.getTotalCpuByPod(i.Name)
	}
}