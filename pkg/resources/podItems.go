package resources

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"kubers/pkg/formatter"
)

type PodItems struct {
	data []Pod

	SortBy    string
	SortOrder string
}

type SortPoBy PodItems

func (m SortPoBy) Len() int {
	return len(m.data)
}

func (m SortPoBy) Less(i, j int) bool {
	switch m.SortBy {
	case "cpu":
		if m.SortOrder == "asc" {
			return m.data[i].totalCpuUsage < m.data[j].totalCpuUsage
		} else {
			return m.data[i].totalCpuUsage > m.data[j].totalCpuUsage
		}
	case "mem":
		if m.SortOrder == "asc" {
			return m.data[i].totalMemUsage < m.data[j].totalMemUsage
		} else {
			return m.data[i].totalMemUsage > m.data[j].totalMemUsage
		}
	default:
		return true
	}
}

func (m SortPoBy) Swap(i, j int) {
	m.data[i], m.data[j] = m.data[j], m.data[i]
}

func (items *PodItems) BuildData(itemsPodMetricsList []v1beta1.PodMetrics, itemsPodList []v1.Pod) {
	items.data = make([]Pod, len(itemsPodMetricsList))
	for i, poM := range itemsPodMetricsList {
		items.data[i] = Pod{
			name:       poM.Name,
			containers: []*ContainerResources{},
		}

		for _, c := range poM.Containers {
			podContainer := &ContainerResources{}
			podContainer.name = c.Name
			podContainer.cpuUsage = c.Usage.Cpu().MilliValue()
			podContainer.memUsage, _ = c.Usage.Memory().AsInt64()

			items.data[i].containers = append(items.data[i].containers, podContainer)
		}

		// Get corresponding pod
		for _, po := range itemsPodList {
			for _, c := range po.Spec.Containers {
				for cIndex, cMetrics := range items.data[i].containers {
					if c.Name == cMetrics.name {
						items.data[i].containers[cIndex].cpuRequest = c.Resources.Requests.Cpu().MilliValue()
						items.data[i].containers[cIndex].cpuLimit = c.Resources.Limits.Cpu().MilliValue()
						items.data[i].containers[cIndex].memRequest, _ = c.Resources.Requests.Memory().AsInt64()
						items.data[i].containers[cIndex].memLimit, _ = c.Resources.Limits.Memory().AsInt64()

						break
					}
				}
			}
		}

		items.data[i].calculateTotalCpuUsage()
		items.data[i].calculateTotalMemUsage()
		items.data[i].calculateTotalCpuSpec()
		items.data[i].calculateTotalMemSpec()
	}
}

func (items *PodItems) FormatForPrint(byCo bool) [][]string {
	var itemsForPrint [][]string
	for _, pod := range items.data {
		if byCo == false {
			itemsForPrint = append(itemsForPrint, pod.getRowsByPod()...)
		} else {
			itemsForPrint = append(itemsForPrint, pod.getRowsByContainer()...)
		}
	}

	return itemsForPrint
}

func (items *PodItems) getTotalMemUsage() int64 {
	var mem int64
	for _, po := range items.data {
		mem += po.totalMemUsage
	}

	return mem
}

func (items *PodItems) getTotalMemRequest() int64 {
	var mem int64
	for _, po := range items.data {
		mem += po.totalMemRequest
	}

	return mem
}

func (items *PodItems) getTotalMemLimit() int64 {
	var mem int64
	for _, po := range items.data {
		mem += po.totalMemLimit
	}

	return mem
}

func (items *PodItems) getTotalCpuUsage() int64 {
	var cpu int64
	for _, po := range items.data {
		cpu += po.totalCpuUsage
	}

	return cpu
}

func (items *PodItems) getTotalCpuRequest() int64 {
	var cpu int64
	for _, po := range items.data {
		cpu += po.totalCpuRequest
	}

	return cpu
}

func (items *PodItems) getTotalCpuLimit() int64 {
	var cpu int64
	for _, po := range items.data {
		cpu += po.totalCpuLimit
	}

	return cpu
}

func (items *PodItems) GetTotalCpuFormatted() string {
	return formatter.FormatCpuUsageWithSpec(items.getTotalCpuUsage(), items.getTotalCpuRequest(), items.getTotalCpuLimit())
}

func (items *PodItems) GetTotalMemFormatted() string {
	return formatter.FormatMemUsageWithSpec(items.getTotalMemUsage(), items.getTotalMemRequest(), items.getTotalMemLimit())
}
