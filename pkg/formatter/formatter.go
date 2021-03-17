package formatter

import "strconv"

func formatMem(mem int64) string {
	memGb := float64(mem) * 0.000000001
	return strconv.FormatFloat(memGb, 'f', 3, 64)
}

func formatCpu(cpu int64) string {
	cpuCores := float64(cpu) * 0.001
	return strconv.FormatFloat(cpuCores, 'f', 3, 64)
}

func FormatCpuUsageWithSpec(usage int64, request int64, limit int64) string {
	return formatCpu(usage) + " (" + formatCpu(request) + "/" + formatCpu(limit) + ")"
}

func FormatMemUsageWithSpec(usage int64, request int64, limit int64) string {
	return formatMem(usage) + " (" + formatMem(request) + "/" + formatMem(limit) + ")"
}
