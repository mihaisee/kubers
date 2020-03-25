package cluster

import (
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

type PoMetricsItems struct {
	data     map[string]map[string]PoResourcesUsage
	sortData []SortItem

	dataForPrint [][]string

	sortDirection string
	sortBy        string
}

type PoResourcesUsage struct {
	cpuUsage   int64
	memUsage   int64
}

type SortItem struct {
	podName   string
	resources PoResourcesUsage
}

type SortPoBy PoMetricsItems

func (m SortPoBy) Len() int {
	return len(m.sortData)
}

func (m SortPoBy) Less(i, j int) bool {
	switch m.sortBy {
	case "cpu":
		if m.sortDirection == "asc" {
			return m.sortData[i].resources.cpuUsage < m.sortData[j].resources.cpuUsage
		} else {
			return m.sortData[i].resources.cpuUsage > m.sortData[j].resources.cpuUsage
		}
	case "memUsage":
		if m.sortDirection == "asc" {
			return m.sortData[i].resources.memUsage < m.sortData[j].resources.memUsage
		} else {
			return m.sortData[i].resources.memUsage > m.sortData[j].resources.memUsage
		}
	default:
		return true
	}
}

func (m SortPoBy) Swap(i, j int) {
	m.sortData[i], m.sortData[j] = m.sortData[j], m.sortData[i]
}

func (items *PoMetricsItems) buildData(itemsList []v1beta1.PodMetrics) {
	items.data = make(map[string]map[string]PoResourcesUsage)
	items.sortData = make([]SortItem, len(itemsList))
	for ind, i := range itemsList {
		items.sortData[ind] = SortItem{podName: i.Name, resources: PoResourcesUsage{memUsage: 0, cpuUsage: 0}}
		items.data[i.Name] = map[string]PoResourcesUsage{}
		for _, c := range i.Containers {
			memUsage, _ := c.Usage.Memory().AsInt64()
			items.data[i.Name][c.Name] = PoResourcesUsage{cpuUsage: c.Usage.Cpu().MilliValue(), memUsage: memUsage}
		}

		items.sortData[ind].resources.memUsage = items.getTotalMemByPod(i.Name)
		items.sortData[ind].resources.cpuUsage = items.getTotalCpuByPod(i.Name)
	}
}

func (items *PoMetricsItems) formatForPrint(byPod bool) [][]string {
	for _, row := range items.sortData {
		if byPod == false {
			items.dataForPrint = append(items.dataForPrint, items.getRowsByContainer(row.podName)...)
		} else {
			items.dataForPrint = append(items.dataForPrint, items.getRowsByPod(row.podName)...)
		}
	}

	return items.dataForPrint
}

func (items PoMetricsItems) getRowsByContainer(podName string) [][]string {
	var itemsForPrint [][]string
	var i = 0
	for c, rs := range items.data[podName] {
		if i != 0 {
			podName = ""
		}

		row := []string{podName, c, rs.getCpuFormatted(), rs.getMemoryFormatted()}
		itemsForPrint = append(itemsForPrint, row)

		i++
	}

	return itemsForPrint
}

func (items PoMetricsItems) getRowsByPod(podName string) [][]string {
	var itemsForPrint [][]string

	row := []string{podName, items.getAllContainersByPod(podName), items.getTotalCpuByPodFormatted(podName), items.getTotalMemByPodFormatted(podName)}
	itemsForPrint = append(itemsForPrint, row)

	return itemsForPrint
}

func (items PoMetricsItems) getAllContainersByPod(pod string) string {
	var c string
	for cn := range items.data[pod] {
		if c == "" {
			c = cn
		} else {
			c += ", " + cn
		}
	}

	return c
}

func (items PoMetricsItems) getTotalMem() int64 {
	var tm int64
	for po := range items.data {
		tm += items.getTotalMemByPod(po)
	}

	return tm
}

func (items PoMetricsItems) getTotalMemFormatted() string {
	return formatMem(items.getTotalMem())
}

func (items PoMetricsItems) getTotalCpu() int64 {
	var tc int64
	for po := range items.data {
		tc += items.getTotalCpuByPod(po)
	}

	return tc
}

func (items PoMetricsItems) getTotalCpuFormatted() string {
	return formatCpu(items.getTotalCpu())
}

func (items PoMetricsItems) getTotalMemByPodFormatted(pod string) string {
	return formatMem(items.getTotalMemByPod(pod))
}

func (items PoMetricsItems) getTotalCpuByPodFormatted(pod string) string {
	return formatCpu(items.getTotalCpuByPod(pod))
}

func (items PoMetricsItems) getTotalMemByPod(pod string) int64 {
	var t int64
	for _, rs := range items.data[pod] {
		t += rs.memUsage
	}

	return t
}

func (items PoMetricsItems) getTotalCpuByPod(pod string) int64 {
	var t int64
	for _, rs := range items.data[pod] {
		t += rs.cpuUsage
	}

	return t
}

func (item PoResourcesUsage) getMemoryFormatted() string {
	return formatMem(item.memUsage)
}

func (item PoResourcesUsage) getCpuFormatted() string {
	return formatCpu(item.cpuUsage)
}
