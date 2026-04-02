package service

import (
	"fmt"
	"log/slog"
	"maps"
	"slices"

	"github.com/acobaugh/osrelease"
	"github.com/prometheus/procfs"
	"github.com/prometheus/procfs/blockdevice"
	"github.com/prometheus/procfs/sysfs"
	"zhacked.me/oxyl/agent/internal/models"
)

var excludedFilesystems = map[string]struct{}{
	"sysfs":       {},
	"proc":        {},
	"devtmpfs":    {},
	"devpts":      {},
	"tmpfs":       {},
	"efivarfs":    {},
	"securityfs":  {},
	"cgroup2":     {},
	"pstore":      {},
	"bpf":         {},
	"autofs":      {},
	"mqueue":      {},
	"hugetlbfs":   {},
	"debugfs":     {},
	"tracefs":     {},
	"fusectl":     {},
	"configfs":    {},
	"binfmt_misc": {},
	"nsfs":        {},
	"overlay":     {},
}

type SystemInfoService struct {
	sysStats  *procfs.FS
	sysInfo   *sysfs.FS
	blockInfo *blockdevice.FS

	systemOs string

	cpuName  string
	cpuCount uint64

	totalMemory uint64

	// ---- Mappings for my own use ----
	MountPointMap    map[string]*procfs.MountInfo
	DiskPartitionMap map[string]*models.DiskInfo
	RaidMap          map[string]*procfs.MDStat
}

func NewSystemInfoService() (*SystemInfoService, error) {
	sysStats, err := procfs.NewDefaultFS()
	if err != nil {
		return nil, fmt.Errorf("unable to create procfs: %v", err)
	}

	sysInfo, err := sysfs.NewDefaultFS()
	if err != nil {
		return nil, fmt.Errorf("unable to create procfs: %v", err)
	}

	blockDeviceInfo, err := blockdevice.NewDefaultFS()
	if err != nil {
		return nil, fmt.Errorf("unable to create procfs: %v", err)
	}

	return &SystemInfoService{
		sysStats:         &sysStats,
		sysInfo:          &sysInfo,
		blockInfo:        &blockDeviceInfo,
		MountPointMap:    make(map[string]*procfs.MountInfo),
		DiskPartitionMap: make(map[string]*models.DiskInfo),
		RaidMap:          make(map[string]*procfs.MDStat),
	}, nil
}

func (s *SystemInfoService) CaptureData() error {
	osInfo, err := osrelease.Read()
	if err != nil {
		return fmt.Errorf("unable to read os info: %v", err)
	}

	s.systemOs = osInfo["PRETTY_NAME"]

	cpus, err := s.sysStats.CPUInfo()
	if err != nil {
		return fmt.Errorf("unable to get cpu info: %v", err)
	}

	s.cpuCount = uint64(len(cpus))
	s.cpuName = cpus[0].ModelName

	memInfo, err := s.sysStats.Meminfo()
	if err != nil {
		return fmt.Errorf("unable to get memory info: %v", err)
	}

	s.totalMemory = *memInfo.MemTotal

	mountPoints, err := s.sysStats.GetMounts()
	if err != nil {
		return fmt.Errorf("unable to get mount points info: %v", err)
	}

	// we need to exclude the partitions that are from the system itself
	// Excluding tmpfs filesystems, proc filesystems, sys filesystems, etc.
	for _, mountPointInfo := range mountPoints {
		if _, ok := excludedFilesystems[mountPointInfo.FSType]; ok {
			continue
		}

		s.MountPointMap[mountPointInfo.Source] = mountPointInfo
	}

	blockDevices, err := s.blockInfo.SysBlockDevices()
	if err != nil {
		return fmt.Errorf("unable to get block devices info: %v", err)
	}

	for _, blockDevice := range blockDevices {
		size, err := s.blockInfo.SysBlockDeviceSize(blockDevice)

		if err != nil {
			slog.Error("failed to stat size of disk", slog.String("device", blockDevice))
			continue
		}

		if diskInfo, err := models.NewDiskInfo(blockDevice, size); err == nil {
			s.DiskPartitionMap[blockDevice] = diskInfo
		}

	}

	raidDevices, err := s.sysStats.MDStat()
	if err != nil {
		return fmt.Errorf("unable to get raid devices info: %v", err)
	}

	for _, raidDevice := range raidDevices {
		s.RaidMap[raidDevice.Name] = &raidDevice
	}

	return nil
}

func (s *SystemInfoService) GetSystemOs() string {
	return s.systemOs
}

func (s *SystemInfoService) GetTotalMemory() uint64 {
	return s.totalMemory
}

func (s *SystemInfoService) GetCPUCount() uint64 {
	return s.cpuCount
}

func (s *SystemInfoService) GetCPUName() string {
	return s.cpuName
}

func (s *SystemInfoService) GetBlockDevices() []string {
	return slices.Collect(maps.Keys(s.DiskPartitionMap))
}

func (s *SystemInfoService) GetMountPoints() []string {
	return slices.Collect(maps.Keys(s.MountPointMap))
}

func (s *SystemInfoService) GetRaidDevices() []string {
	return slices.Collect(maps.Keys(s.RaidMap))
}

func (s *SystemInfoService) GetBlockDevice(device string) (*models.DiskInfo, bool) {
	blockDevice, found := s.DiskPartitionMap[device]
	return blockDevice, found
}

func (s *SystemInfoService) GetMountPoint(mountPoint string) (*procfs.MountInfo, bool) {
	mountPointInfo, found := s.MountPointMap[mountPoint]
	return mountPointInfo, found
}

func (s *SystemInfoService) GetRaidDevice(raidDevice string) (*procfs.MDStat, bool) {
	raidDeviceInfo, found := s.RaidMap[raidDevice]
	return raidDeviceInfo, found
}
