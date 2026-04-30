package models

type GeneralMetricsAvg struct {
	AvgCPU    float64
	AvgMemory float64
}

type DiskUsageAvg struct {
	MountPoint   string
	AvgDiskUsage float64
}

type DiskHealthAvg struct {
	DiskPath          string
	AvgHealthLeft     float64
	AvgErrorRate      float64
	AvgPendingSectors float64
}

type NetworkAvg struct {
	InterfaceName string
	AvgRXRate     float64
	AvgTXRate     float64
}
