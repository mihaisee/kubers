package resources

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"kubers/pkg/formatter"
)

type PoItems struct {
	data     map[string]map[string]PoItemResources
	sortData []SortItem

	dataForPrint [][]string

	SortDirection string
	SortBy        string
}

type PoItemResources struct {
	cpuUsage   int64
	cpuRequest int64
	cpuLimit   int64

	memUsage   int64
	memRequest int64
	memLimit   int64
}

type SortItem struct {
	podName   string
	resources PoItemResources
}

type SortPoBy PoItems

func (m SortPoBy) Len() int {
	return len(m.sortData)
}

func (m SortPoBy) Less(i, j int) bool {
	switch m.SortBy {
	case "cpu":
		if m.SortDirection == "asc" {
			return m.sortData[i].resources.cpuUsage < m.sortData[j].resources.cpuUsage
		} else {
			return m.sortData[i].resources.cpuUsage > m.sortData[j].resources.cpuUsage
		}
	case "memUsage":
		if m.SortDirection == "asc" {
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

func (items *PoItems) BuildData(itemsList []v1beta1.PodMetrics, podsList []v1.Pod) {
	items.data = make(map[string]map[string]PoItemResources)
	items.sortData = make([]SortItem, len(itemsList))
	for ind, i := range itemsList {
		items.sortData[ind] = SortItem{podName: i.Name, resources: PoItemResources{memUsage: 0, cpuUsage: 0}}
		items.data[i.Name] = map[string]PoItemResources{}
		for _, c := range i.Containers {
			memUsage, _ := c.Usage.Memory().AsInt64()
			items.data[i.Name][c.Name] = PoItemResources{cpuUsage: c.Usage.Cpu().MilliValue(), memUsage: memUsage}
		}

		items.sortData[ind].resources.memUsage = items.getTotalMemByPod(i.Name)
		items.sortData[ind].resources.cpuUsage = items.getTotalCpuByPod(i.Name)
	}
}

func (items *PoItems) FormatForPrint(byPod bool) [][]string {
	for _, row := range items.sortData {
		if byPod == false {
			items.dataForPrint = append(items.dataForPrint, items.getRowsByContainer(row.podName)...)
		} else {
			items.dataForPrint = append(items.dataForPrint, items.getRowsByPod(row.podName)...)
		}
	}

	return items.dataForPrint
}

func (items PoItems) getRowsByContainer(podName string) [][]string {
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

func (items PoItems) getRowsByPod(podName string) [][]string {
	var itemsForPrint [][]string

	row := []string{podName, items.getAllContainersByPod(podName), items.getTotalCpuByPodFormatted(podName), items.getTotalMemByPodFormatted(podName)}
	itemsForPrint = append(itemsForPrint, row)

	return itemsForPrint
}

func (items PoItems) getAllContainersByPod(pod string) string {
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

func (items PoItems) getTotalMem() int64 {
	var tm int64
	for po := range items.data {
		tm += items.getTotalMemByPod(po)
	}

	return tm
}

func (items PoItems) GetTotalMemFormatted() string {
	return formatter.FormatMem(items.getTotalMem())
}

func (items PoItems) getTotalCpu() int64 {
	var tc int64
	for po := range items.data {
		tc += items.getTotalCpuByPod(po)
	}

	return tc
}

func (items PoItems) GetTotalCpuFormatted() string {
	return formatter.FormatCpu(items.getTotalCpu())
}

func (items PoItems) getTotalMemByPodFormatted(pod string) string {
	return formatter.FormatMem(items.getTotalMemByPod(pod))
}

func (items PoItems) getTotalCpuByPodFormatted(pod string) string {
	return formatter.FormatCpu(items.getTotalCpuByPod(pod))
}

func (items PoItems) getTotalMemByPod(pod string) int64 {
	var t int64
	for _, rs := range items.data[pod] {
		t += rs.memUsage
	}

	return t
}

func (items PoItems) getTotalCpuByPod(pod string) int64 {
	var t int64
	for _, rs := range items.data[pod] {
		t += rs.cpuUsage
	}

	return t
}

func (item PoItemResources) getMemoryFormatted() string {
	return formatter.FormatMem(item.memUsage)
}

func (item PoItemResources) getCpuFormatted() string {
	return formatter.FormatCpu(item.cpuUsage)
}
