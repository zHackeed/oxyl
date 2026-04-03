package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
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

	oldCpuUsage uint64

	v1.MonitoringServiceClient
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
		oldCpuUsage:    9999,
	}, nil
}

func (s *MonitoringService) Start(ctx context.Context) error {
	tick := time.NewTicker(10 * time.Second)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				if err := s.Consume(ctx, 10*time.Second); err != nil {
					slog.Error("failed to consume monitoring data", "error", err)
				}
			}
		}
	}()

	return nil
}

func (s *MonitoringService) Consume(ctx context.Context, tickRate time.Duration) error {
	var g errgroup.Group

	generalizeChan := make(chan *monitoring.GeneralMetrics, 1)
	mountPointInfoChan := make(chan []*monitoring.MountedDiskMetrics, 1)
	//physicalDiskInfoChan := make(chan []*monitoring.PhysicalDiskMetrics, 1)
	networkInfoChan := make(chan []*monitoring.NetworkMetrics, 1)

	defer func() {
		close(generalizeChan)
		close(mountPointInfoChan)
		//close(physicalDiskInfoChan)
		close(networkInfoChan)
	}()

	g.Go(func() error {
		return s.getGeneralInfo(generalizeChan, tickRate)
	})

	g.Go(func() error {
		return s.getMountPointInfo(mountPointInfoChan)
	})

	/*
		g.Go(func() error {
			return s.getPhysicalMetrics(physicalDiskInfoChan)
		})

	*/

	g.Go(func() error {
		return s.getNetworkInfo(networkInfoChan, tickRate)
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("failed to consume monitoring service: %w", err)
	}

	generalInfo, ok := <-generalizeChan
	if !ok {
		return fmt.Errorf("failed to get general info")
	}

	mountPointInfo, ok := <-mountPointInfoChan
	if !ok {
		return fmt.Errorf("failed to get mount point info")
	}

	/*
		physicalDiskInfo, ok := <-physicalDiskInfoChan
		if !ok {
			return fmt.Errorf("failed to get physical disk info")
		}
	*/

	networkInfo, ok := <-networkInfoChan
	if !ok {
		return fmt.Errorf("failed to get network info")
	}

	// we do not care if the metrics are sent or not. We will always send them regardless
	if _, err := s.SendMetrics(ctx, &monitoring.AgentMetrics{
		GeneralMetrics: generalInfo,
		DiskMetrics:    mountPointInfo,
		//PhysicalDiskMetrics: physicalDiskInfo,
		NetworkMetrics: networkInfo,
	}); err != nil {
		return fmt.Errorf("failed to send metrics: %w", err)
	}

	return nil
}

func (s *MonitoringService) getGeneralInfo(consumer chan<- *monitoring.GeneralMetrics, rate time.Duration) error {
	timeStart := time.Now()
	defer func() {
		slog.Info("crawled general info", "time", time.Since(timeStart))
	}()

	generalData := new(monitoring.GeneralMetrics)

	cpuStats, err := s.procFs.Stat()
	if err != nil {
		return fmt.Errorf("failed to get stats: %w", err)
	}

	currentCpuUsage := uint64(cpuStats.CPUTotal.User + cpuStats.CPUTotal.Nice + cpuStats.CPUTotal.System)

	if s.oldCpuUsage == 9999 {
		s.oldCpuUsage = currentCpuUsage
	}

	generalData.CpuUsage = (currentCpuUsage - s.oldCpuUsage) / uint64(rate.Seconds())
	s.oldCpuUsage = currentCpuUsage

	ramMetrics, err := s.procFs.Meminfo()
	if err != nil {
		return fmt.Errorf("failed to get memory metrics: %w", err)
	}

	generalData.MemoryUsage = *ramMetrics.Active
	generalData.Uptime = uint64(time.Now().Unix()) - cpuStats.BootTime

	consumer <- generalData
	return nil
}

func (s *MonitoringService) getMountPointInfo(consumer chan<- []*monitoring.MountedDiskMetrics) error {
	timeStart := time.Now()
	defer func() {
		slog.Info("crawled mount point info", "time", time.Since(timeStart))
	}()

	mountPointInfo := make([]*monitoring.MountedDiskMetrics, 0)

	for _, mountPoint := range s.sysInfoService.GetMountPoints() {
		slog.Info("disk partition", "mountPoint", mountPoint)
		var stat syscall.Statfs_t // So, the kernel directly does not have any way to obtain the filesystem stats
		if err := syscall.Statfs(mountPoint, &stat); err != nil {
			return fmt.Errorf("failed to get mount info stats: %w", err)
		}

		mountPointInfo = append(mountPointInfo, &monitoring.MountedDiskMetrics{
			MountPoint: mountPoint,
			UsedSpace:  (stat.Blocks - stat.Bfree) * uint64(stat.Bsize),
		})
	}

	consumer <- mountPointInfo
	return nil
}

func (s *MonitoringService) getPhysicalMetrics(consumer chan<- []*monitoring.PhysicalDiskMetrics) error {
	timeStart := time.Now()
	defer func() {
		slog.Info("crawled physical disk info", "time", time.Since(timeStart))
	}()

	physicalDiskInfo := make([]*monitoring.PhysicalDiskMetrics, 0)

	for _, disk := range s.sysInfoService.GetBlockDevices() {
		smtDto, err := smart.Open("/dev/" + nvmeController(disk))
		if err != nil {
			return fmt.Errorf("failed to get smart data of disk %s: %w", disk, err)
		}

		if err != nil {
			_ = smtDto.Close()
			return fmt.Errorf("failed to get smart data of disk %s: %w", disk, err)
		}

		switch smtDto.(type) {
		case *smart.NVMeDevice:
			nvmeSmart, _ := smtDto.(*smart.NVMeDevice)
			smartData, err := nvmeSmart.ReadSMART()

			if err != nil {
				_ = smtDto.Close()
				return fmt.Errorf("failed to get smart data of disk %s: %w", disk, err)
			}

			physicalDiskInfo = append(physicalDiskInfo, &monitoring.PhysicalDiskMetrics{
				DiskPath:      disk,
				HealthUsed:    new(uint32(smartData.PercentUsed)), // nasty
				MediaErrors_1: &smartData.MediaErrors.Val[0],
				MediaErrors_2: &smartData.MediaErrors.Val[1],
			})

		case *smart.SataDevice:
			sataSmart, _ := smtDto.(*smart.SataDevice)
			smartData, err := sataSmart.ReadSMARTData()
			if err != nil {
				_ = smtDto.Close()
				return fmt.Errorf("failed to get smart data: %w", err)
			}

			//https://en.wikipedia.org/wiki/Self-Monitoring,_Analysis_and_Reporting_Technology#Known_ATA_S.M.A.R.T._attributes
			// I despise this

			readErrorRate, exists := smartData.Attrs[1]
			if !exists {
				_ = smtDto.Close()
				// So this disk does not follow standard metrics
				continue
			}

			reallocatedSectors, exists := smartData.Attrs[5]
			if !exists {
				_ = smtDto.Close()
				// So this disk does not follow standard metrics
				continue
			}

			physicalDiskInfo = append(physicalDiskInfo, &monitoring.PhysicalDiskMetrics{
				DiskPath:       disk,
				ErrorRate:      new(uint64(readErrorRate.Current)),
				PendingSectors: new(uint32(reallocatedSectors.Current)),
			})
		default:
			slog.Debug("unknown device type", "device", disk)
		}

		_ = smtDto.Close()
	}

	consumer <- physicalDiskInfo
	return nil
}

func (s *MonitoringService) getNetworkInfo(consumer chan<- []*monitoring.NetworkMetrics, rate time.Duration) error {
	start := time.Now()

	netMetrics := make([]*monitoring.NetworkMetrics, 0)

	netStat, err := s.procFs.NetDev()
	if err != nil {
		return fmt.Errorf("failed to get netstat: %w", err)
	}

	for _, networkInterface := range netStat {
		if networkInterface.Name == "lo" || strings.HasPrefix(networkInterface.Name, "veth") || strings.HasPrefix(networkInterface.Name, "br-") || strings.HasPrefix(networkInterface.Name, "docker") {
			continue
		}

		bytesIn, packetsIn := networkInterface.RxBytes, networkInterface.RxPackets
		bytesOut, packetsOut := networkInterface.TxBytes, networkInterface.TxPackets

		oldValueBpTx, oldValuePkTx := s.byteCounterMap[networkInterface.Name+"Tx"], s.pksCounterMap[networkInterface.Name+"Tx"]
		oldValueBpRx, oldValuePkRx := s.byteCounterMap[networkInterface.Name+"Rx"], s.pksCounterMap[networkInterface.Name+"Rx"]

		if oldValueBpTx == 0 || oldValuePkTx == 0 || oldValuePkRx == 0 || oldValueBpRx == 0 {
			s.byteCounterMap[networkInterface.Name+"Tx"] = bytesOut
			s.pksCounterMap[networkInterface.Name+"Tx"] = packetsOut
			s.byteCounterMap[networkInterface.Name+"Rx"] = bytesIn
			s.pksCounterMap[networkInterface.Name+"Rx"] = packetsIn
			continue
		}

		valueBpTx := (networkInterface.TxBytes - oldValueBpTx) / uint64(rate.Seconds())
		valuePkTx := (networkInterface.TxPackets - oldValuePkTx) / uint64(rate.Seconds())
		valueBpRx := (networkInterface.RxBytes - oldValueBpRx) / uint64(rate.Seconds())
		valuePkRx := (networkInterface.RxPackets - oldValuePkRx) / uint64(rate.Seconds())

		s.byteCounterMap[networkInterface.Name+"Tx"] = bytesOut
		s.pksCounterMap[networkInterface.Name+"Tx"] = packetsOut
		s.byteCounterMap[networkInterface.Name+"Rx"] = bytesIn
		s.pksCounterMap[networkInterface.Name+"Rx"] = packetsIn

		netMetrics = append(netMetrics, &monitoring.NetworkMetrics{
			InterfaceName: networkInterface.Name,

			BytesSentRate:   valueBpTx,
			PacketsSentRate: valuePkTx,

			BytesReceivedRate:   valueBpRx,
			PacketsReceivedRate: valuePkRx,

			BytesSent:   bytesOut,
			PacketsSent: packetsOut,

			BytesReceived:   bytesIn,
			PacketsReceived: packetsIn,
		})
	}

	consumer <- netMetrics

	slog.Info("crawled network info", "time", time.Since(start))
	return nil
}

func nvmeController(device string) string {
	// nvme0n1 → nvme0
	if strings.HasPrefix(device, "nvme") {
		if idx := strings.Index(device, "n"); idx != -1 {
			return device[:idx]
		}
	}
	return device
}
