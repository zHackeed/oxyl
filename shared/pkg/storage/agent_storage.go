package storage

import (
	"context"
	sql2 "database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"zhacked.me/oxyl/shared/pkg/datasource"
	"zhacked.me/oxyl/shared/pkg/models"
)

var (
	ErrAgentNotFound = errors.New("agent not found")
	ErrNoAgents      = errors.New("no agents found for this company")
)

type AgentStorage struct {
	conn *datasource.TimescaleConnection
}

func NewAgentStorage(persistence *datasource.TimescaleConnection) *AgentStorage {
	return new(AgentStorage{
		conn: persistence,
	})
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

	if len(agent.Partitions) > 0 {
		for _, partition := range agent.Partitions {
			sql := `INSERT INTO agent_partition_scheme (agent, mount_point, total_size, is_raid, raid_level) VALUES ($1, $2, $3, $4, $5)`
			if _, err := tx.Exec(ctx, sql, agent.ID, partition.MountPoint, partition.TotalSize, partition.Raid, partition.RaidLevel); err != nil {
				return fmt.Errorf("unable to create agent partition scheme: %w", err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}

	return nil
}

func (a *AgentStorage) GetAgent(ctx context.Context, agentID string) (*models.Agent, error) {
	sql := `SELECT id, holder, display_name, registered_ip, status, system_os, cpu_model, total_memory, total_disk, last_handshake, created_at FROM agents WHERE id = $1`
	row := a.conn.Pool().QueryRow(ctx, sql, agentID)
	var agent models.Agent

	// systemOs, cpuModel, totalMemory, totalDisk are nullable
	var systemOS, cpuModel sql2.NullString
	var totalMemory, totalDisk sql2.NullInt64
	var lastHandshake sql2.NullTime

	if err := row.Scan(
		&agent.ID,
		&agent.Holder, &agent.DisplayName,
		&agent.RegisteredIP, &agent.Status,
		&systemOS, &cpuModel, &totalMemory, &totalDisk,
		&lastHandshake, &agent.CreatedAt); err != nil {
		return nil, fmt.Errorf("unable to get agent: %w", err)
	}

	if agent.Status == models.AgentStatusEnrolling {
		// We do not know the data of the server
		return &agent, nil
	}

	agent.SystemOS = systemOS.String
	agent.CPUModel = cpuModel.String
	agent.TotalMemory = totalMemory.Int64
	agent.TotalDisk = totalDisk.Int64

	// the last connection might not have been ever done
	if lastHandshake.Valid {
		agent.LastHandshake = lastHandshake.Time
	}

	sql = `SELECT mount_point, total_size, is_raid, raid_level FROM agent_partition_scheme WHERE agent = $1`
	rows, err := a.conn.Pool().Query(ctx, sql, agentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAgentNotFound
		}
		return nil, fmt.Errorf("unable to get agent partitions: %w", err)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("unable to get agent partitions: %w", rows.Err())
	}

	agent.Partitions = make([]*models.AgentPartition, 0)

	defer rows.Close()
	for rows.Next() {
		var partition models.AgentPartition
		if err := rows.Scan(&partition.MountPoint, &partition.TotalSize, &partition.Raid, &partition.RaidLevel); err != nil {
			return nil, fmt.Errorf("unable to get agent partitions: %w", err)
		}
		agent.Partitions = append(agent.Partitions, &partition)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("unable to get agent partitions: %w", rows.Err())
	}

	return &agent, nil
}

// This is the shittiest way to do this.
func (a *AgentStorage) GetAgents(ctx context.Context) ([]*models.Agent, error) {
	sql := `
        SELECT 
            ag.id, ag.holder, ag.display_name, ag.registered_ip, 
            ag.status, ag.system_os, ag.cpu_model, ag.total_memory, 
            ag.total_disk, ag.last_handshake, ag.created_at,
            ap.mount_point, ap.total_size, ap.is_raid, ap.raid_level
        FROM agents ag
        LEFT JOIN agent_partition_scheme ap 
            ON ag.id = ap.agent AND ag.status != 'ENROLLING'
        ORDER BY ag.id
    `

	rows, err := a.conn.Pool().Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("unable to get agents: %w", err)
	}
	defer rows.Close()

	agentMap := make(map[string]*models.Agent)
	agentOrder := make([]string, 0)

	for rows.Next() {
		var agentID string
		var agent models.Agent

		// Nullable values
		var mountPoint *string
		var totalSize *int64
		var isRaid *bool
		var raidLevel *int

		if err := rows.Scan(
			&agentID, &agent.Holder, &agent.DisplayName, &agent.RegisteredIP,
			&agent.Status, &agent.SystemOS, &agent.CPUModel, &agent.TotalMemory,
			&agent.TotalDisk, &agent.LastHandshake, &agent.CreatedAt,
			&mountPoint, &totalSize, &isRaid, &raidLevel,
		); err != nil {
			return nil, fmt.Errorf("unable to scan agent row: %w", err)
		}

		// If the agent is not in the map, add it and add it to the order. (Issues with the left join)
		if _, exists := agentMap[agentID]; !exists {
			agent.ID = agentID
			agent.Partitions = make([]*models.AgentPartition, 0)
			agentMap[agentID] = &agent
			agentOrder = append(agentOrder, agentID)
		}

		// If the mount point is not null, add the partition to the agent active agent
		if mountPoint != nil {
			agentMap[agentID].Partitions = append(agentMap[agentID].Partitions, &models.AgentPartition{
				MountPoint: *mountPoint,
				TotalSize:  *totalSize,
				Raid:       *isRaid,
				RaidLevel:  *raidLevel,
			})
		}
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("unable to get agents: %w", rows.Err())
	}

	agents := make([]*models.Agent, 0, len(agentOrder))
	for _, id := range agentOrder {
		agents = append(agents, agentMap[id])
	}

	return agents, nil

}

func (a *AgentStorage) GetAgentsOfCompany(ctx context.Context, companyID string) ([]*models.Agent, error) {
	sql := `
        SELECT 
            ag.id, ag.holder, ag.display_name, ag.registered_ip, 
            ag.status, ag.system_os, ag.cpu_model, ag.total_memory, 
            ag.total_disk, ag.last_handshake, ag.created_at,
            ap.mount_point, ap.total_size, ap.is_raid, ap.raid_level
        FROM agents ag
        LEFT JOIN agent_partition_scheme ap 
            ON ag.id = ap.agent AND ag.status != 'ENROLLING'
        WHERE ag.holder = $1
        ORDER BY ag.id
    `

	rows, err := a.conn.Pool().Query(ctx, sql, companyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoAgents
		}
		return nil, fmt.Errorf("unable to get agents of company: %w", err)
	}

	agentMap := make(map[string]*models.Agent)
	agentOrder := make([]string, 0) // preserve order
	defer rows.Close()

	for rows.Next() {
		var agentID string
		var agent models.Agent
		var mountPoint *string

		// Nullable values
		var totalSize *int64
		var isRaid *bool
		var raidLevel *int

		if err := rows.Scan(
			&agentID, &agent.Holder, &agent.DisplayName, &agent.RegisteredIP,
			&agent.Status, &agent.SystemOS, &agent.CPUModel, &agent.TotalMemory,
			&agent.TotalDisk, &agent.LastHandshake, &agent.CreatedAt,
			&mountPoint, &totalSize, &isRaid, &raidLevel,
		); err != nil {
			return nil, fmt.Errorf("unable to scan agent row: %w", err)
		}

		if _, exists := agentMap[agentID]; !exists {
			agent.ID = agentID
			agent.Partitions = make([]*models.AgentPartition, 0)
			agentMap[agentID] = &agent
			agentOrder = append(agentOrder, agentID)
		}
		if mountPoint != nil {
			agentMap[agentID].Partitions = append(agentMap[agentID].Partitions, &models.AgentPartition{
				MountPoint: *mountPoint,
				TotalSize:  *totalSize,
				Raid:       *isRaid,
				RaidLevel:  *raidLevel,
			})
		}
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("unable to get agents of company: %w", rows.Err())
	}

	agents := make([]*models.Agent, 0, len(agentOrder))
	for _, id := range agentOrder {
		agents = append(agents, agentMap[id])
	}

	return agents, nil
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

func (a *AgentStorage) EnrichAgent(ctx context.Context, agentID string, systemOS, cpuModel string, totalMemory, totalDisk int64, enrollmentToken string, partitions []*models.AgentPartition) error {
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
              WHERE id = $7`

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
