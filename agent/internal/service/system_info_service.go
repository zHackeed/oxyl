package service

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"

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
	"fuse.portal": {},
	"zfs":         {},
	"zram":        {},
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
	PartitionMap     map[string]*models.DiskInfo
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
		PartitionMap:     make(map[string]*models.DiskInfo),
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

		if !strings.HasPrefix(blockDevice, "nvme") {
			nameRune := []rune(blockDevice)

			typeValue := string(nameRune[:2])
			if typeValue != "sd" && typeValue != "fd" && typeValue != "hd" {
				slog.Info("invalid partition type", slog.String("type", typeValue))
				continue
			}
		}

		if diskInfo, err := models.NewDiskInfo(nil, blockDevice, size); err == nil {
			slog.Info("block device name", slog.String("device", blockDevice))
			s.DiskPartitionMap[blockDevice] = diskInfo
		}

	}

	partitions, err := s.consumePartitions()
	if err != nil {
		return fmt.Errorf("unable to get partitions info: %v", err)
	}

	s.PartitionMap = partitions

	raidDevices, err := s.sysStats.MDStat()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// we do not have any virtual raid arrays
			return nil
		}

		return fmt.Errorf("unable to get raid devices info: %v", err)
	}

	for _, raidDevice := range raidDevices {
		s.RaidMap[raidDevice.Name] = &raidDevice
	}

	return nil
}

// So profs does not have a method to pass the partition info
// So we have to implement our own to have it done properly for the system, as the block devices would normally
// do not need to be checked only for the smartctl info.
func (s *SystemInfoService) consumePartitions() (map[string]*models.DiskInfo, error) {
	file, err := os.Open("/proc/partitions")
	if err != nil {
		return nil, fmt.Errorf("unable to open /proc/partitions: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan() // skip header
	scanner.Scan() // empty space?

	// major minor #blocks name

	partitionMap := make(map[string]*models.DiskInfo)

	for scanner.Scan() {
		line := scanner.Text()

		fields := strings.Fields(line)
		if len(fields) != 4 {
			slog.Info("invalid partition line", slog.String("line", line))
			continue
		}

		// nvmeXnXpX
		// sdX|fdX|hdXpX

		partitionName := fields[3]
		partitionSize := fields[2]

		diskName, err := s.obtainDiskFromPartition(partitionName)
		if err != nil {
			continue
		}

		sizeValue, err := strconv.Atoi(partitionSize)
		if err != nil {
			continue
		}

		diskInfo, err := models.NewDiskInfo(diskName, partitionName, uint64(sizeValue*1024))
		if err != nil {
			slog.Info("unable to create disk info", slog.String("disk", *diskName), slog.String("partition", partitionName), err)
			continue
		}

		slog.Info("partition size", slog.String("disk", *diskName), slog.String("partition", partitionName), slog.Uint64("size", diskInfo.TotalSize), slog.String("gb", strconv.FormatUint(diskInfo.TotalSize/1024/1024/1024, 10)))

		partitionMap[partitionName] = diskInfo
	}

	return partitionMap, nil
}

// https://askubuntu.com/questions/56929/what-is-the-linux-drive-naming-scheme
func (s *SystemInfoService) obtainDiskFromPartition(partition string) (*string, error) {
	if strings.HasPrefix(partition, "nvme") {
		splitter := strings.Split(partition, "p")
		return new(splitter[0]), nil
	}

	runes := []rune(partition)

	if len(runes) < 4 {
		return nil, fmt.Errorf("invalid partition: %s", partition)
	}

	typeValue := string(runes[:2])
	if typeValue != "sd" && typeValue != "fd" && typeValue != "hd" {
		slog.Info("invalid partition type", slog.String("type", typeValue))
		return nil, fmt.Errorf("invalid partition type: %s", partition)
	}

	indexValue := string(runes[2])
	partitionValue := string(runes[3:])

	slog.Info("obtained disk from partition", slog.String("type", typeValue), slog.String("index", indexValue), slog.String("partition", partitionValue))

	return new(strings.Join([]string{typeValue, indexValue}, "")), nil
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

func (s *SystemInfoService) GetPartition(device string) (*models.DiskInfo, bool) {
	partition, found := s.PartitionMap[device]
	return partition, found
}

func (s *SystemInfoService) GetMountPoint(mountPoint string) (*procfs.MountInfo, bool) {
	mountPointInfo, found := s.MountPointMap[mountPoint]
	return mountPointInfo, found
}

func (s *SystemInfoService) GetRaidDevice(raidDevice string) (*procfs.MDStat, bool) {
	raidDeviceInfo, found := s.RaidMap[raidDevice]
	return raidDeviceInfo, found
}
