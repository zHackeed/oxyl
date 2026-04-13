package models

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/oklog/ulid/v2"
)

type AgentPartition struct {
	MountPoint string `json:"mount_point"`

	TotalSize uint64 `json:"total_size"` // Bytes, might have to change

	Raid      bool `json:"raid"`
	RaidLevel int  `json:"raid_level,omitempty"`
}

type Agent struct {
	ID string `json:"id"`

	Holder string `json:"holder"` // Company ID that is the owner of the agent

	DisplayName  string `json:"display_name"`
	RegisteredIP net.IP `json:"registered_ip"`

	Status AgentStatus `json:"status"`

	Metadata *AgentMetadata `json:"metadata,omitempty"`

	LastHandshake time.Time `json:"last_handshake,omitempty"`
	LastUpdated   time.Time `json:"last_updated,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type AgentMetadata struct {
	SystemOS string `json:"system_os"`
	CPUModel string `json:"cpu_model"`

	TotalMemory uint64 `json:"total_memory"` // Bytes, might have to change
	TotalDisk   uint64 `json:"total_disk"`   // Bytes, might have to change

	Partitions []*AgentPartition `json:"partitions"`

	EnrollmentToken string `json:"-"`
}

func NewAgentMetadata(systemOS, cpuModel string, totalMemory, totalDisk uint64, partitions []*AgentPartition, enrollmentToken string) (*AgentMetadata, error) {
	if systemOS == "" || cpuModel == "" || totalMemory == 0 || totalDisk == 0 || len(partitions) <= 0 || enrollmentToken == "" {
		return nil, errors.New("invalid agent metadata")
	}

	return &AgentMetadata{
		SystemOS:        systemOS,
		CPUModel:        cpuModel,
		TotalDisk:       totalDisk,
		Partitions:      partitions,
		EnrollmentToken: enrollmentToken,
	}, nil
}

func NewPartition(mountPoint string, totalSize uint64, raid bool, raidLevel int) (*AgentPartition, error) {
	if mountPoint == "" {
		return nil, errors.New("mount point is empty")
	}

	if totalSize <= 0 {
		return nil, errors.New("total size is not valid")
	}

	if raid && raidLevel <= 0 {
		return nil, errors.New("raid level is not valid")
	}

	return &AgentPartition{
		MountPoint: mountPoint,
		TotalSize:  totalSize,
		Raid:       raid,
		RaidLevel:  raidLevel,
	}, nil
}

func NewAgent(displayName, registeredIP, holder string) (*Agent, error) {
	if holder == "" {
		return nil, errors.New("holder is empty")
	}

	if len(holder) > 26 {
		return nil, errors.New("holder is too long, maybe malformed")
	}

	if displayName == "" {
		return nil, errors.New("display name is empty")
	}

	if len(displayName) > 255 {
		return nil, errors.New("display name is too long")
	}

	if registeredIP == "" {
		return nil, errors.New("registered ip is empty")
	}

	parsedIp := net.ParseIP(registeredIP)
	if parsedIp == nil {
		return nil, fmt.Errorf("the registered ip %q is not valid, maybe malformed", registeredIP)
	}

	return &Agent{
		ID:           ulid.Make().String(),
		Holder:       holder,
		DisplayName:  displayName,
		RegisteredIP: parsedIp,
		Status:       AgentStatusEnrolling,
		LastUpdated:  time.Now(),
		CreatedAt:    time.Now(),
	}, nil
}

func (a *Agent) UpdateDisplayName(displayName string) error {
	if displayName == "" {
		return errors.New("display name is empty")
	}

	if len(displayName) > 255 {
		return errors.New("display name is too long")
	}

	a.DisplayName = displayName
	a.LastUpdated = time.Now()
	return nil
}

func (a *Agent) UpdateRegisteredIP(registeredIP string) error {
	if registeredIP == "" {
		return errors.New("registered ip is empty")
	}

	parsedIp := net.ParseIP(registeredIP)
	if parsedIp == nil {
		return fmt.Errorf("the registered ip %q is not valid, maybe malformed", registeredIP)
	}

	a.RegisteredIP = parsedIp
	a.LastUpdated = time.Now()
	return nil
}

func (a *Agent) SetMetadata(metadata *AgentMetadata) error {
	if metadata == nil {
		return errors.New("metadata is empty")
	}

	a.Metadata = metadata
	return nil
}

func (a *AgentMetadata) AddPartition(partition *AgentPartition) {
	a.Partitions = append(a.Partitions, partition)
}

func (a *AgentMetadata) RemovePartition(mountPoint string) {
	for i, p := range a.Partitions {
		if p.MountPoint == mountPoint {
			a.Partitions = append(a.Partitions[:i], a.Partitions[i+1:]...)
			break
		}
	}
}

func (a *Agent) UpdateStatus(status AgentStatus) {
	a.Status = status
	a.LastUpdated = time.Now()
}
