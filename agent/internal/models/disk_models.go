package models

import "fmt"

type DiskInfo struct {
	Name      string
	TotalSize uint64
}

func NewDiskInfo(name string, totalSize uint64) (*DiskInfo, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	if totalSize == 0 {
		return nil, fmt.Errorf("total size cannot be zero")
	}

	return &DiskInfo{
		Name:      name,
		TotalSize: totalSize,
	}, nil
}
