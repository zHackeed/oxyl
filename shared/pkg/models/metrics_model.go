package models

import (
	"errors"
	"time"
)

type AgentGeneralMetrics struct {
	When        time.Time `json:"when,omitempty"`
	CPUUsage    uint64    `json:"cpu_usage"`
	MemoryUsage uint64    `json:"memory_usage"`
	Uptime      uint64    `json:"uptime"`
}

func NewGeneralMetrics(cpuUsage uint64, memoryUsage uint64, uptime uint64) (*AgentGeneralMetrics, error) {
	if cpuUsage > 100 {
		return nil, errors.New("cpu usage cannot be greater than 100")
	}

	return &AgentGeneralMetrics{
		CPUUsage:    cpuUsage,
		MemoryUsage: memoryUsage,
		Uptime:      uptime,
	}, nil
}

type AgentMountPointMetrics struct {
	When       time.Time `json:"when,omitempty"`
	MountPoint string    `json:"mount_point"`
	DiskUsage  uint64    `json:"disk_usage"`
}

func NewMountPointMetrics(mountPoint string, diskUsage uint64) (*AgentMountPointMetrics, error) {
	if len(mountPoint) == 0 {
		return nil, errors.New("mount point is empty")
	}

	return &AgentMountPointMetrics{
		MountPoint: mountPoint,
		DiskUsage:  diskUsage,
	}, nil
}

type AgentPhysicalDiskMetrics struct {
	When     time.Time `json:"when,omitempty"`
	DiskPath string    `json:"disk_path"`

	HealthUsed   uint32 `json:"health_used,omitempty"`
	MediaErrors1 uint64 `json:"media_errors_1,omitempty"`
	MediaErrors2 uint64 `json:"media_errors_2,omitempty"`

	ErrorRate      uint64 `json:"error_rate,omitempty"`
	PendingSectors uint32 `json:"pending_sectors,omitempty"`
}

func NewAgentPhysicalDiskMetrics(
	diskPath string, healthUsed *uint32,
	mediaErrors1 *uint64, mediaErrors2 *uint64,
	errorRate *uint64, pendingSectors *uint32) (*AgentPhysicalDiskMetrics, error) {

	if len(diskPath) == 0 {
		return nil, errors.New("disk path is empty")
	}

	agent := &AgentPhysicalDiskMetrics{
		DiskPath: diskPath,
	}

	if healthUsed != nil && mediaErrors1 != nil && mediaErrors2 != nil {
		agent.HealthUsed = *healthUsed
		agent.MediaErrors1 = *mediaErrors1
		agent.MediaErrors2 = *mediaErrors2
	} else {
		if errorRate == nil || pendingSectors == nil {
			return nil, errors.New("sata disk must provide error rate and pending sectors")
		}

		agent.ErrorRate = *errorRate
		agent.PendingSectors = *pendingSectors
	}

	return agent, nil
}

type AgentNetworkMetrics struct {
	When          time.Time `json:"when,omitempty"`
	InterfaceName string    `json:"interface_name"`

	RxBytes uint64 `json:"rx_bytes"`
	TxBytes uint64 `json:"tx_bytes"`

	RxPackets uint64 `json:"rx_packets"`
	TxPackets uint64 `json:"tx_packets"`

	RxRate uint64 `json:"rx_rate"`
	TxRate uint64 `json:"tx_rate"`

	RxPacketRate uint64 `json:"rx_packet_rate"`
	TxPacketRate uint64 `json:"tx_packet_rate"`
}

func NewAgentNetworkMetrics(ifName string, rxBytes, txBytes, rxPackets, txPackets, rxRate, txRate, rxPacketRate, txPacketRate uint64) (*AgentNetworkMetrics, error) {
	if len(ifName) == 0 {
		return nil, errors.New("interface name is empty")
	}

	return &AgentNetworkMetrics{
		InterfaceName: ifName,

		RxBytes: rxBytes,
		TxBytes: txBytes,

		RxPackets: rxPackets,
		TxPackets: txPackets,

		RxRate: rxRate,
		TxRate: txRate,

		RxPacketRate: rxPacketRate,
		TxPacketRate: txPacketRate,
	}, nil
}

type AgentMetrics struct {
	AgentId string `json:"agent_id"`

	GeneralMetrics      []*AgentGeneralMetrics                 `json:"general_metrics"`
	MountPointMetrics   map[string][]*AgentMountPointMetrics   `json:"mount_point_metrics"`
	PhysicalDiskMetrics map[string][]*AgentPhysicalDiskMetrics `json:"physical_disk_metrics"`
	NetworkMetrics      map[string][]*AgentNetworkMetrics      `json:"network_metrics"`
	From                time.Time                              `json:"from"`
	To                  time.Time                              `json:"to"`
}

func NewAgentMetrics(agentId string,
	generalMetrics []*AgentGeneralMetrics,
	mountMetrics []*AgentMountPointMetrics,
	physicalDiskMetrics []*AgentPhysicalDiskMetrics,
	networkMetrics []*AgentNetworkMetrics, from time.Duration) (*AgentMetrics, error) {
	if len(agentId) != 26 {
		return nil, errors.New("agent id is invalid")
	}

	if len(generalMetrics) == 0 {
		return nil, errors.New("general metrics cannot be empty")
	}

	if len(mountMetrics) == 0 {
		return nil, errors.New("mount metrics cannot be empty")
	}

	if len(physicalDiskMetrics) == 0 {
		return nil, errors.New("physical disk metrics cannot be empty")
	}

	if len(networkMetrics) == 0 {
		return nil, errors.New("network metrics cannot be empty")
	}

	mappedMountPointMetrics := make(map[string][]*AgentMountPointMetrics)
	for _, metric := range mountMetrics {
		mappedMountPointMetrics[metric.MountPoint] = append(mappedMountPointMetrics[metric.MountPoint], metric)
	}

	mappedPhysicalDisks := make(map[string][]*AgentPhysicalDiskMetrics)
	for _, metric := range physicalDiskMetrics {
		mappedPhysicalDisks[metric.DiskPath] = append(mappedPhysicalDisks[metric.DiskPath], metric)
	}

	mappedNetworkMetrics := make(map[string][]*AgentNetworkMetrics)
	for _, metric := range networkMetrics {
		mappedNetworkMetrics[metric.InterfaceName] = append(mappedNetworkMetrics[metric.InterfaceName], metric)
	}

	return &AgentMetrics{
		AgentId: agentId,

		GeneralMetrics:      generalMetrics,
		MountPointMetrics:   mappedMountPointMetrics,
		PhysicalDiskMetrics: mappedPhysicalDisks,
		NetworkMetrics:      mappedNetworkMetrics,

		From: time.Now().Add(-from),
		To:   time.Now(),
	}, nil
}
