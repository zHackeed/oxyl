//go:build linux

package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"

	v1 "zhacked.me/oxyl/protocol/v1"
	"zhacked.me/oxyl/protocol/v1/enrollment"
)

const enrollmentFile = "/etc/oxyl/enrollment.id"

type EnrollmentService struct {
	enrollmentToken *string

	systemInfoService *SystemInfoService

	v1.EnrollmentServiceClient
}

func NewEnrollmentService(systemInfoService *SystemInfoService) *EnrollmentService {
	return &EnrollmentService{
		systemInfoService: systemInfoService,
	}
}

func (s *EnrollmentService) Start(ctx context.Context) error {
	if _, err := os.Stat(enrollmentFile); err == nil {
		f, err := os.Open(enrollmentFile)
		if err != nil {
			return fmt.Errorf("unable to open enrollment file: %v", err)
		}
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("unable to read enrollment file: %v", err)
		}

		s.enrollmentToken = new(string(data))
		return nil
	}

	partitions := make([]*enrollment.DiskPartition, 0)
	totalSize := uint64(0)

	for _, partition := range s.systemInfoService.GetMountPoints() {
		partitionData, _ := s.systemInfoService.GetMountPoint(partition)

		blockValue, trimmed := strings.CutPrefix(partitionData.Source, "/dev/")
		if !trimmed {
			continue
		}

		blockDevice, _ := s.systemInfoService.GetBlockDevice(blockValue)

		slog.Info(partitionData.Source)

		totalSize = totalSize + blockDevice.TotalSize

		partitionWrapper := &enrollment.DiskPartition{
			MountPoint: partition,
			TotalSize:  blockDevice.TotalSize,
		}

		if strings.HasPrefix(blockDevice.Name, "md") {
			partitionWrapper.IsRaid = true
			mdData, _ := s.systemInfoService.GetRaidDevice(blockDevice.Name)

			level, _ := strings.CutPrefix(mdData.Type, "raid")

			raidLevel, err := strconv.Atoi(level)
			if err != nil {
				return fmt.Errorf("unable to parse raid level: %v", err)
			}

			partitionWrapper.RaidLevel = new(uint32(raidLevel))
		}

		partitions = append(partitions, partitionWrapper)
	}

	token, err := s.GetEnrollmentToken(ctx, &enrollment.EnrollmentRequest{
		CpuModel:       s.systemInfoService.GetCPUName(),
		OsVariant:      s.systemInfoService.GetSystemOs(),
		TotalMemory:    s.systemInfoService.GetTotalMemory(),
		TotalDisk:      totalSize,
		DiskPartitions: partitions,
	})

	if err != nil {
		return fmt.Errorf("unable to get enrollment token: %v", err)
	}

	s.enrollmentToken = &token.EnrollmentId

	if err := os.WriteFile(enrollmentFile, []byte(*s.enrollmentToken), 0o600); err != nil {
		return fmt.Errorf("unable to write enrollment file: %v", err)
	}

	return nil
}

func (s *EnrollmentService) ProvideEnrollmentIdentifier() (string, error) {
	if s.enrollmentToken == nil {
		return "", fmt.Errorf("enrollment token is not set")
	}

	return *s.enrollmentToken, nil
}
