package storage

import (
	"context"
	sql2 "database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"zhacked.me/oxyl/shared/pkg/datasource"
	"zhacked.me/oxyl/shared/pkg/models"
)

var (
	ErrAgentNotFound    = errors.New("agent not found")
	ErrNoAgents         = errors.New("no agents found for this company")
	ErrAgentNotEnrolled = errors.New("agent not enrolled")
)

type AgentStorage struct {
	conn *datasource.TimescaleConnection
}

func NewAgentStorage(persistence *datasource.TimescaleConnection) *AgentStorage {
	return &AgentStorage{
		conn: persistence,
	}
}

func (a *AgentStorage) CreateAgent(ctx context.Context, agent *models.Agent) error {
	tx, err := a.conn.Pool().Begin(ctx)

	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	sql := `INSERT INTO agents (id, holder, display_name, registered_ip) VALUES ($1, $2, $3, $4)`
	if _, err := tx.Exec(ctx, sql, agent.ID, agent.Holder, agent.DisplayName, agent.RegisteredIP); err != nil {
		return fmt.Errorf("unable to create agent: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}

	return nil
}

func (a *AgentStorage) GetAgentWithMetadata(ctx context.Context, agentID string) (*models.Agent, *models.AgentMetadata, error) {
	sql := `
        SELECT id, holder, display_name, registered_ip, status, last_handshake, created_at,
               system_os, cpu_model, total_memory, total_disk, enrollment_token
        FROM agents WHERE id = $1
    `
	row := a.conn.Pool().QueryRow(ctx, sql, agentID)

	var agent models.Agent
	var lastHandshake sql2.NullTime
	var systemOS, cpuModel, enrollmentToken sql2.NullString
	var totalMemory, totalDisk sql2.NullInt64

	if err := row.Scan(
		&agent.ID, &agent.Holder, &agent.DisplayName,
		&agent.RegisteredIP, &agent.Status,
		&lastHandshake, &agent.CreatedAt,
		&systemOS, &cpuModel, &totalMemory, &totalDisk, &enrollmentToken,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, ErrAgentNotFound
		}
		return nil, nil, fmt.Errorf("unable to get agent: %w", err)
	}

	if lastHandshake.Valid {
		agent.LastHandshake = lastHandshake.Time
	}

	if !systemOS.Valid {
		return &agent, nil, nil
	}

	rows, err := a.conn.Pool().Query(ctx,
		`SELECT mount_point, total_size, is_raid, raid_level FROM agent_partition_scheme WHERE agent = $1`,
		agentID,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get agent partitions: %w", err)
	}
	defer rows.Close()

	partitions := make([]*models.AgentPartition, 0)
	for rows.Next() {
		var p models.AgentPartition
		if err := rows.Scan(&p.MountPoint, &p.TotalSize, &p.Raid, &p.RaidLevel); err != nil {
			return nil, nil, fmt.Errorf("unable to scan partition: %w", err)
		}
		partitions = append(partitions, &p)
	}

	if rows.Err() != nil {
		return nil, nil, fmt.Errorf("unable to get agent partitions: %w", rows.Err())
	}

	return &agent, &models.AgentMetadata{
		SystemOS:        systemOS.String,
		CPUModel:        cpuModel.String,
		TotalMemory:     uint64(totalMemory.Int64),
		TotalDisk:       uint64(totalDisk.Int64),
		EnrollmentToken: enrollmentToken.String,
		Partitions:      partitions,
	}, nil
}

func (a *AgentStorage) GetAgent(ctx context.Context, agentID string) (*models.Agent, error) {
	sql := `SELECT id, holder, display_name, registered_ip, status, last_handshake, created_at FROM agents WHERE id = $1`
	row := a.conn.Pool().QueryRow(ctx, sql, agentID)

	var agent models.Agent
	var lastHandshake sql2.NullTime

	if err := row.Scan(
		&agent.ID, &agent.Holder, &agent.DisplayName,
		&agent.RegisteredIP, &agent.Status,
		&lastHandshake, &agent.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAgentNotFound
		}
		return nil, fmt.Errorf("unable to get agent: %w", err)
	}

	if lastHandshake.Valid {
		agent.LastHandshake = lastHandshake.Time
	}

	return &agent, nil
}

func (a *AgentStorage) GetAgentsOfCompany(ctx context.Context, companyID string) ([]*models.Agent, error) {
	sql := `
        SELECT id, holder, display_name, registered_ip, status, last_handshake, created_at
        FROM agents
        WHERE holder = $1
        ORDER BY id
    `

	rows, err := a.conn.Pool().Query(ctx, sql, companyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoAgents
		}
		return nil, fmt.Errorf("unable to get agents of company: %w", err)
	}
	defer rows.Close()

	agents := make([]*models.Agent, 0)

	for rows.Next() {
		var agent models.Agent
		var lastHandshake sql2.NullTime

		if err := rows.Scan(
			&agent.ID, &agent.Holder, &agent.DisplayName,
			&agent.RegisteredIP, &agent.Status,
			&lastHandshake, &agent.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("unable to scan agent row: %w", err)
		}

		if lastHandshake.Valid {
			agent.LastHandshake = lastHandshake.Time
		}

		agents = append(agents, &agent)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("unable to get agents of company: %w", rows.Err())
	}

	return agents, nil
}
func (a *AgentStorage) GetHolderOfAgent(ctx context.Context, agentId string) (*string, error) {
	sql := `SELECT holder FROM  agents WHERE id = $1`
	row := a.conn.Pool().QueryRow(ctx, sql, agentId)

	holder := new(string)

	if err := row.Scan(holder); err != nil {
		return nil, fmt.Errorf("unable to get agent holder: %w", err)
	}

	return holder, nil
}

func (a *AgentStorage) GetAgentState(ctx context.Context, agentID string) (models.AgentStatus, error) {
	sql := `SELECT status FROM agents WHERE id = $1`
	row := a.conn.Pool().QueryRow(ctx, sql, agentID)
	var status models.AgentStatus
	if err := row.Scan(&status); err != nil {
		return "", fmt.Errorf("unable to get agent status: %w", err)
	}

	return status, nil
}

func (a *AgentStorage) AddPartitionToAgent(ctx context.Context, agentID string, partition *models.AgentPartition) error {
	sql := `INSERT INTO agent_partition_scheme (agent, mount_point, total_size, is_raid, raid_level) VALUES ($1, $2, $3, $4, $5)`
	if _, err := a.conn.Pool().Exec(ctx, sql, agentID, partition.MountPoint, partition.TotalSize, partition.Raid, partition.RaidLevel); err != nil {
		return fmt.Errorf("unable to add partition to agent: %w", err)
	}

	return nil
}

func (a *AgentStorage) RemovePartitionFromAgent(ctx context.Context, agentID, mountPoint string) error {
	sql := `DELETE FROM agent_partition_scheme WHERE agent = $1 AND mount_point = $2`
	if _, err := a.conn.Pool().Exec(ctx, sql, agentID, mountPoint); err != nil {
		return fmt.Errorf("unable to remove partition from agent: %w", err)
	}

	return nil
}

func (a *AgentStorage) UpdateAgentStatus(ctx context.Context, agentID string, status models.AgentStatus) error {
	sql := `UPDATE agents SET status = $1, last_handshake = CURRENT_TIMESTAMP WHERE id = $2`
	if _, err := a.conn.Pool().Exec(ctx, sql, status, agentID); err != nil {
		return fmt.Errorf("unable to update agent status: %w", err)
	}

	return nil
}

func (a *AgentStorage) EnrichAgent(ctx context.Context, agentID string, systemOS, cpuModel string, totalMemory, totalDisk uint64, enrollmentToken string, partitions []*models.AgentPartition) error {
	tx, err := a.conn.Pool().Begin(ctx)
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	sql := `INSERT INTO agent_partition_scheme (agent, mount_point, total_size, is_raid, raid_level) VALUES ($1, $2, $3, $4, $5)`

	for _, partition := range partitions {
		if _, err := tx.Exec(ctx, sql, agentID, partition.MountPoint, partition.TotalSize, partition.Raid, partition.RaidLevel); err != nil {
			return fmt.Errorf("unable to add partition to agent: %w", err)
		}
	}

	sql = `UPDATE agents SET 
                  system_os = $1, cpu_model = $2, 
                  total_memory = $3, total_disk = $4, 
                  enrollment_token = $5, status = 'ACTIVE', 
                  last_handshake = CURRENT_TIMESTAMP, last_update = CURRENT_TIMESTAMP
              WHERE id = $6`

	slog.Info("enrollment", "agent_id", agentID, "system_os", systemOS, "cpu_model", cpuModel, "total_memory", totalMemory, "total_disk", totalDisk, "enrollment_token", enrollmentToken)

	if _, err := tx.Exec(ctx, sql, systemOS, cpuModel, totalMemory, totalDisk, enrollmentToken, agentID); err != nil {
		return fmt.Errorf("unable to enrich agent: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}

	return nil
}

func (a *AgentStorage) DeleteAgent(ctx context.Context, agentID string) error {
	sql := `DELETE FROM agents WHERE id = $1`
	if _, err := a.conn.Pool().Exec(ctx, sql, agentID); err != nil {
		return fmt.Errorf("unable to delete agent: %w", err)
	}

	return nil
}

// --- metrics might be handled here. But right now is not the correct approach.
