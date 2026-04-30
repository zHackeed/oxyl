package models

type AgentMetadata struct {
	Holder string

	TotalMemory uint64
	TotalDisk   uint64

	Partitions map[string]*AgentPartition
}

type AgentPartition struct {
	TotalSize uint64
}
