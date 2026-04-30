package storage

import (
	"context"

	"zhacked.me/oxyl/service/thresholds/internal/models"
	"zhacked.me/oxyl/shared/pkg/datasource"
)

type AgentStorage struct {
	conn *datasource.TimescaleConnection
}

func NewAgentStorage(conn *datasource.TimescaleConnection) *AgentStorage {
	return &AgentStorage{
		conn: conn,
	}
}

func (s *AgentStorage) GetAll(ctx context.Context) (map[string]*models.AgentMetadata, error) {
	sql := `SELECT id, holder, total_memory, total_disk FROM agents`

	rows, err := s.conn.Pool().Query(ctx, sql)
	if err != nil {
		return nil, err
	}

	agents := make(map[string]*models.AgentMetadata)
	defer rows.Close()

	for rows.Next() {
		agent := new(models.AgentMetadata)
		agent.Partitions = make(map[string]*models.AgentPartition)

		var id string

		if err := rows.Scan(&id, &agent.Holder, &agent.TotalMemory, &agent.TotalDisk); err != nil {
			return nil, err
		}

		agents[id] = agent
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	sql = `SELECT agent, mount_point, total_size FROM agent_partition_scheme`

	rows, err = s.conn.Pool().Query(ctx, sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		agentPartition := new(models.AgentPartition)

		var agentID string
		var mountPoint string

		if err := rows.Scan(&agentID, &mountPoint, &agentPartition.TotalSize); err != nil {
			return nil, err
		}

		if a, ok := agents[agentID]; ok {
			a.Partitions[mountPoint] = agentPartition
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return agents, nil
}

func (s *AgentStorage) GetMetadata(ctx context.Context, agentID string) (*models.AgentMetadata, error) {
	sql := `SELECT holder, total_memory, total_disk FROM agents WHERE id = $1`

	agent := &models.AgentMetadata{
		Partitions: make(map[string]*models.AgentPartition),
	}

	if err := s.conn.Pool().QueryRow(ctx, sql, agentID).Scan(&agent.Holder, &agent.TotalMemory, &agent.TotalDisk); err != nil {
		return nil, err
	}

	sql = `SELECT mount_point, total_size FROM agent_partition_scheme WHERE agent = $1`

	rows, err := s.conn.Pool().Query(ctx, sql, agentID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		agentPartition := new(models.AgentPartition)
		var mountPoint string

		if err := rows.Scan(&mountPoint, &agentPartition.TotalSize); err != nil {
			return nil, err
		}

		agent.Partitions[mountPoint] = agentPartition
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return agent, nil
}
