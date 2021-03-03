package formatter

import "strconv"

func FormatMem(mem int64) string {
	memGb := float64(mem) * 0.000000001
	return strconv.FormatFloat(memGb, 'f', 3, 64)
}

func FormatCpu(cpu int64) string {
	cpuCores := float64(cpu) * 0.001
	return strconv.FormatFloat(cpuCores, 'f', 2, 64)
}
