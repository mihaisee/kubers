package resources

import "kubers/pkg/formatter"

type Namespace struct {
	name string

	labels map[string]string

	totalCpuUsage int64
	totalMemUsage int64

	totalCpuRequest int64
	totalMemRequest int64

	totalCpuLimit int64
	totalMemLimit int64

	podItems *PodItems
}

func (ns *Namespace) calculateTotalCpuUsage() {
	ns.totalCpuUsage = ns.podItems.getTotalCpuUsage()
}

func (ns *Namespace) calculateTotalMemUsage() {
	ns.totalMemUsage = ns.podItems.getTotalMemUsage()
}

func (ns *Namespace) calculateTotalCpuSpec() {
	ns.totalCpuRequest, ns.totalCpuLimit = ns.podItems.getTotalCpuRequest(), ns.podItems.getTotalCpuLimit()
}

func (ns *Namespace) calculateTotalMemSpec() {
	ns.totalMemRequest, ns.totalMemLimit = ns.podItems.getTotalMemRequest(), ns.podItems.getTotalMemLimit()
}

func (ns *Namespace) GetTotalCpuFormatted() string {
	return formatter.FormatCpuUsageWithSpec(ns.totalCpuUsage, ns.totalCpuRequest, ns.totalCpuLimit)
}

func (ns *Namespace) GetTotalMemFormatted() string {
	return formatter.FormatMemUsageWithSpec(ns.totalMemUsage, ns.totalMemRequest, ns.totalMemLimit)
}
