package models

import "fmt"

type DiskInfo struct {
	Holder    *string
	Name      string
	TotalSize uint64
}

func NewDiskInfo(holder *string, name string, totalSize uint64) (*DiskInfo, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	if totalSize == 0 {
		return nil, fmt.Errorf("total size cannot be zero")
	}

	return &DiskInfo{
		Holder:    holder,
		Name:      name,
		TotalSize: totalSize,
	}, nil
}
