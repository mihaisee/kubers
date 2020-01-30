package cluster

import "k8s.io/metrics/pkg/apis/metrics/v1beta1"

type PoItems struct {
	data     map[string]map[string]PoItemResources
	sortData []SortItem

	dataForPrint [][]string

	sortDirection string
	sortBy        string
}

type PoItemResources struct {
	cpu int64
	mem int64
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
	switch m.sortBy {
	case "cpu":
		if m.sortDirection == "asc" {
			return m.sortData[i].resources.cpu < m.sortData[j].resources.cpu
		} else {
			return m.sortData[i].resources.cpu > m.sortData[j].resources.cpu
		}
	case "mem":
		if m.sortDirection == "asc" {
			return m.sortData[i].resources.mem < m.sortData[j].resources.mem
		} else {
			return m.sortData[i].resources.mem > m.sortData[j].resources.mem
		}
	default:
		return true
	}
}

func (m SortPoBy) Swap(i, j int) {
	m.sortData[i], m.sortData[j] = m.sortData[j], m.sortData[i]
}

func (items *PoItems) buildData(itemsList []v1beta1.PodMetrics) {
	items.data = make(map[string]map[string]PoItemResources)
	items.sortData = make([]SortItem, len(itemsList))
	for ind, i := range itemsList {
		items.sortData[ind] = SortItem{podName: i.Name, resources: PoItemResources{mem: 0, cpu: 0}}
		items.data[i.Name] = map[string]PoItemResources{}
		for _, c := range i.Containers {
			mem, _ := c.Usage.Memory().AsInt64()
			items.data[i.Name][c.Name] = PoItemResources{c.Usage.Cpu().MilliValue(), mem}
		}

		items.sortData[ind].resources.mem = items.getTotalMemByPod(i.Name)
		items.sortData[ind].resources.cpu = items.getTotalCpuByPod(i.Name)
	}
}

func (items *PoItems) formatForPrint(byPod bool) [][]string {
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

func (items PoItems) getTotalMemFormatted() string {
	return formatMem(items.getTotalMem())
}

func (items PoItems) getTotalCpu() int64 {
	var tc int64
	for po := range items.data {
		tc += items.getTotalCpuByPod(po)
	}

	return tc
}

func (items PoItems) getTotalCpuFormatted() string {
	return formatCpu(items.getTotalCpu())
}

func (items PoItems) getTotalMemByPodFormatted(pod string) string {
	return formatMem(items.getTotalMemByPod(pod))
}

func (items PoItems) getTotalCpuByPodFormatted(pod string) string {
	return formatCpu(items.getTotalCpuByPod(pod))
}

func (items PoItems) getTotalMemByPod(pod string) int64 {
	var t int64
	for _, rs := range items.data[pod] {
		t += rs.mem
	}

	return t
}

func (items PoItems) getTotalCpuByPod(pod string) int64 {
	var t int64
	for _, rs := range items.data[pod] {
		t += rs.cpu
	}

	return t
}

func (item PoItemResources) getMemoryFormatted() string {
	return formatMem(item.mem)
}

func (item PoItemResources) getCpuFormatted() string {
	return formatCpu(item.cpu)
}
