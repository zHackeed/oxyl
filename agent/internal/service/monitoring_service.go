package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/anatol/smart.go"
	"github.com/prometheus/procfs"
	"golang.org/x/sync/errgroup"
	v1 "zhacked.me/oxyl/protocol/v1"
	"zhacked.me/oxyl/protocol/v1/monitoring"
)

// I hate this

type MonitoringService struct {
	sysInfoService *SystemInfoService

	procFs procfs.FS

	byteCounterMap map[string]uint64
	pksCounterMap  map[string]uint64

	oldCpuUsage *cpuUsage

	v1.MonitoringServiceClient
}

type cpuUsage struct {
	total float64
	idle  float64
}

func NewMonitoringService(sysInfoService *SystemInfoService) (*MonitoringService, error) {
	procFSIns, err := procfs.NewDefaultFS()
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %w", err)
	}

	return &MonitoringService{
		sysInfoService: sysInfoService,
		procFs:         procFSIns,

		byteCounterMap: make(map[string]uint64),
		pksCounterMap:  make(map[string]uint64),
	}, nil
}

func (s *MonitoringService) Start(ctx context.Context) error {
	tick := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				if err := s.Consume(ctx, 1*time.Second); err != nil {
					slog.Error("failed to consume monitoring data", "error", err)
				}
			}
		}
	}()

	return nil
}

func (s *MonitoringService) Consume(ctx context.Context, tickRate time.Duration) error {
	var (
		generalInfo      *monitoring.GeneralMetrics
		mountPointInfo   []*monitoring.MountedDiskMetrics
		physicalDiskInfo []*monitoring.PhysicalDiskMetrics
		networkInfo      []*monitoring.NetworkMetrics
	)

	var g errgroup.Group

	g.Go(func() (err error) { generalInfo, err = s.getGeneralInfo(); return })
	g.Go(func() (err error) { mountPointInfo, err = s.getMountPointInfo(); return })
	g.Go(func() (err error) { physicalDiskInfo, err = s.getPhysicalMetrics(); return })
	g.Go(func() (err error) { networkInfo, err = s.getNetworkInfo(tickRate); return })

	if err := g.Wait(); err != nil {
		return fmt.Errorf("failed to consume monitoring service: %w", err)
	}

	if _, err := s.SendMetrics(ctx, &monitoring.AgentMetrics{
		GeneralMetrics:      generalInfo,
		DiskMetrics:         mountPointInfo,
		PhysicalDiskMetrics: physicalDiskInfo,
		NetworkMetrics:      networkInfo,
	}); err != nil {
		return fmt.Errorf("failed to send metrics: %w", err)
	}

	return nil
}

func (s *MonitoringService) getGeneralInfo() (*monitoring.GeneralMetrics, error) {
	timeStart := time.Now()
	defer func() {
		slog.Info("crawled general info", "time", time.Since(timeStart))
	}()

	cpuStats, err := s.procFs.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	prev := s.oldCpuUsage
	curr := &cpuUsage{
		total: cpuStats.CPUTotal.User + cpuStats.CPUTotal.Nice + cpuStats.CPUTotal.System +
			cpuStats.CPUTotal.Idle + cpuStats.CPUTotal.Iowait + cpuStats.CPUTotal.IRQ +
			cpuStats.CPUTotal.SoftIRQ + cpuStats.CPUTotal.Steal,
		idle: cpuStats.CPUTotal.Idle + cpuStats.CPUTotal.Iowait,
	}

	var cpuPercent float64
	if prev != nil {
		deltaTotal := curr.total - prev.total
		deltaIdle := curr.idle - prev.idle

		if deltaTotal > 0 {
			cpuPercent = (deltaTotal - deltaIdle) / deltaTotal * 100
		}
	}

	s.oldCpuUsage = curr

	ramMetrics, err := s.procFs.Meminfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory metrics: %w", err)
	}

	slog.Info("got memory metrics", "mem", *ramMetrics.MemTotalBytes)
	slog.Info("free", "value", *ramMetrics.MemAvailableBytes)
	slog.Info("value", *ramMetrics.MemTotalBytes-*ramMetrics.MemAvailableBytes)

	slog.Info("cpu usage", "value", cpuPercent)

	return &monitoring.GeneralMetrics{
		CpuUsage:    uint64(cpuPercent),
		MemoryUsage: *ramMetrics.MemTotalBytes - *ramMetrics.MemAvailableBytes,
		Uptime:      uint64(time.Now().Unix()) - cpuStats.BootTime,
	}, nil
}

func (s *MonitoringService) getMountPointInfo() ([]*monitoring.MountedDiskMetrics, error) {
	timeStart := time.Now()
	defer func() {
		slog.Info("crawled mount point info", "time", time.Since(timeStart))
	}()

	mountPointInfo := make([]*monitoring.MountedDiskMetrics, 0)

	for _, mountPoint := range s.sysInfoService.GetMountPoints() {
		var stat syscall.Statfs_t
		if err := syscall.Statfs(mountPoint, &stat); err != nil {
			return nil, fmt.Errorf("failed to get mount info stats: %w", err)
		}

		mountPointInfo = append(mountPointInfo, &monitoring.MountedDiskMetrics{
			MountPoint: mountPoint,
			UsedSpace:  (stat.Blocks - stat.Bfree) * uint64(stat.Bsize),
		})
	}

	return mountPointInfo, nil
}

func (s *MonitoringService) getPhysicalMetrics() ([]*monitoring.PhysicalDiskMetrics, error) {
	timeStart := time.Now()
	defer func() {
		slog.Info("crawled physical disk info", "time", time.Since(timeStart))
	}()

	var (
		mu               sync.Mutex
		wg               sync.WaitGroup
		physicalDiskInfo = make([]*monitoring.PhysicalDiskMetrics, 0)
	)

	for _, disk := range s.sysInfoService.GetBlockDevices() {
		wg.Add(1)
		go func(disk string) {
			defer wg.Done()
			metrics, err := s.consumeDisk(disk)
			if err != nil {
				slog.Debug("skipping device", "device", disk, "error", err)
				return
			}
			mu.Lock()
			physicalDiskInfo = append(physicalDiskInfo, metrics)
			mu.Unlock()
		}(disk)
	}

	wg.Wait()
	return physicalDiskInfo, nil
}

func (s *MonitoringService) consumeDisk(disk string) (*monitoring.PhysicalDiskMetrics, error) {
	physicalDiskInfo := new(monitoring.PhysicalDiskMetrics)

	smtDto, err := smart.Open("/dev/" + disk)
	if err != nil {
		return nil, fmt.Errorf("failed to get smart data of disk %s: %w", disk, err)
	}
	defer smtDto.Close()

	switch smtDto.(type) {
	case *smart.NVMeDevice:
		nvmeSmart, _ := smtDto.(*smart.NVMeDevice)
		smartData, err := nvmeSmart.ReadSMART()

		if err != nil {
			return nil, fmt.Errorf("failed to get smart data of disk %s: %w", disk, err)
		}

		physicalDiskInfo.DiskPath = disk
		physicalDiskInfo.HealthUsed = new(uint32(smartData.PercentUsed)) // nasty
		physicalDiskInfo.MediaErrors_1 = &smartData.MediaErrors.Val[0]
		physicalDiskInfo.MediaErrors_2 = &smartData.MediaErrors.Val[1]

	case *smart.SataDevice:
		sataSmart, _ := smtDto.(*smart.SataDevice)
		smartData, err := sataSmart.ReadSMARTData()
		if err != nil {
			return nil, fmt.Errorf("failed to get smart data: %w", err)
		}

		//https://en.wikipedia.org/wiki/Self-Monitoring,_Analysis_and_Reporting_Technology#Known_ATA_S.M.A.R.T._attributes
		// I despise this

		readErrorRate, exists := smartData.Attrs[1]
		if !exists {
			// So this disk does not follow standard metrics
			return nil, fmt.Errorf("disk %s does not follow standard metrics", disk)
		}

		reallocatedSectors, exists := smartData.Attrs[5]
		if !exists {
			// So this disk does not follow standard metric
			return nil, fmt.Errorf("disk %s does not follow standard metrics", disk)
		}

		physicalDiskInfo.DiskPath = disk
		physicalDiskInfo.ErrorRate = new(uint64(readErrorRate.Current))
		physicalDiskInfo.PendingSectors = new(uint32(reallocatedSectors.Current))
	default:
		slog.Debug("unknown device type", "device", disk)
	}

	return physicalDiskInfo, nil
}

func (s *MonitoringService) getNetworkInfo(rate time.Duration) ([]*monitoring.NetworkMetrics, error) {
	timeStart := time.Now()
	defer func() {
		slog.Info("crawled network info", "time", time.Since(timeStart))
	}()

	netStat, err := s.procFs.NetDev()
	if err != nil {
		return nil, fmt.Errorf("failed to get netstat: %w", err)
	}

	netMetrics := make([]*monitoring.NetworkMetrics, 0)

	for _, networkInterface := range netStat {
		if networkInterface.Name == "lo" ||
			strings.HasPrefix(networkInterface.Name, "veth") ||
			strings.HasPrefix(networkInterface.Name, "br-") ||
			strings.HasPrefix(networkInterface.Name, "docker") {
			continue
		}

		bytesIn, packetsIn := networkInterface.RxBytes, networkInterface.RxPackets
		bytesOut, packetsOut := networkInterface.TxBytes, networkInterface.TxPackets

		oldValueBpTx := s.byteCounterMap[networkInterface.Name+"Tx"]
		oldValuePkTx := s.pksCounterMap[networkInterface.Name+"Tx"]
		oldValueBpRx := s.byteCounterMap[networkInterface.Name+"Rx"]
		oldValuePkRx := s.pksCounterMap[networkInterface.Name+"Rx"]

		s.byteCounterMap[networkInterface.Name+"Tx"] = bytesOut
		s.pksCounterMap[networkInterface.Name+"Tx"] = packetsOut
		s.byteCounterMap[networkInterface.Name+"Rx"] = bytesIn
		s.pksCounterMap[networkInterface.Name+"Rx"] = packetsIn

		if oldValueBpTx == 0 || oldValuePkTx == 0 || oldValuePkRx == 0 || oldValueBpRx == 0 {
			continue
		}

		netMetrics = append(netMetrics, &monitoring.NetworkMetrics{
			InterfaceName:       networkInterface.Name,
			BytesSentRate:       (bytesOut - oldValueBpTx) / uint64(rate.Seconds()),
			PacketsSentRate:     (packetsOut - oldValuePkTx) / uint64(rate.Seconds()),
			BytesReceivedRate:   (bytesIn - oldValueBpRx) / uint64(rate.Seconds()),
			PacketsReceivedRate: (packetsIn - oldValuePkRx) / uint64(rate.Seconds()),
			BytesSent:           bytesOut,
			PacketsSent:         packetsOut,
			BytesReceived:       bytesIn,
			PacketsReceived:     packetsIn,
		})
	}

	return netMetrics, nil
}
