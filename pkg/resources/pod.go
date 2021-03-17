package resources

import "kubers/pkg/formatter"

type Pod struct {
	name string

	totalCpuUsage int64
	totalMemUsage int64

	totalCpuRequest int64
	totalMemRequest int64

	totalCpuLimit int64
	totalMemLimit int64

	containers []*ContainerResources
}

type ContainerResources struct {
	name string

	cpuUsage int64
	memUsage int64

	cpuRequest int64
	cpuLimit   int64

	memRequest int64
	memLimit   int64
}

func (pod *Pod) getRowsByContainer() [][]string {
	var itemsForPrint [][]string
	for i, container := range pod.containers {
		var podName = pod.name
		if i != 0 {
			podName = ""
		}

		row := []string{
			podName,
			container.name,
			container.GetTotalCpuFormatted(),
			container.GetTotalMemFormatted(),
		}
		itemsForPrint = append(itemsForPrint, row)

		i++
	}

	return itemsForPrint
}

func (pod *Pod) getRowsByPod() [][]string {
	var itemsForPrint [][]string

	row := []string{
		pod.name,
		pod.getAllContainersFormatted(),
		pod.GetTotalCpuFormatted(),
		pod.GetTotalMemFormatted(),
	}
	itemsForPrint = append(itemsForPrint, row)

	return itemsForPrint
}

func (pod *Pod) getAllContainersFormatted() string {
	var c string
	for _, cn := range pod.containers {
		if c == "" {
			c = cn.name
		} else {
			c += ", " + cn.name
		}
	}

	return c
}

func (pod *Pod) GetTotalMemFormatted() string {
	return formatter.FormatMemUsageWithSpec(pod.totalMemUsage, pod.totalMemRequest, pod.totalMemLimit)
}

func (pod *Pod) GetTotalCpuFormatted() string {
	return formatter.FormatCpuUsageWithSpec(pod.totalCpuUsage, pod.totalCpuRequest, pod.totalCpuLimit)
}

func (pod *Pod) calculateTotalMemUsage() {
	for _, c := range pod.containers {
		pod.totalMemUsage += c.memUsage
	}
}

func (pod *Pod) calculateTotalMemSpec() {
	for _, c := range pod.containers {
		pod.totalMemRequest += c.memRequest
		pod.totalMemLimit += c.memLimit
	}
}

func (pod *Pod) calculateTotalCpuUsage() {
	for _, c := range pod.containers {
		pod.totalCpuUsage += c.cpuUsage
	}
}

func (pod *Pod) calculateTotalCpuSpec() {
	for _, c := range pod.containers {
		pod.totalCpuRequest += c.cpuRequest
		pod.totalCpuLimit += c.cpuLimit
	}
}

func (c *ContainerResources) GetTotalCpuFormatted() string {
	return formatter.FormatCpuUsageWithSpec(c.cpuUsage, c.cpuRequest, c.cpuLimit)
}

func (c *ContainerResources) GetTotalMemFormatted() string {
	return formatter.FormatMemUsageWithSpec(c.memUsage, c.memRequest, c.memLimit)
}
