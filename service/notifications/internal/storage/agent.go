package storage

import (
	"context"

	"zhacked.me/oxyl/shared/pkg/datasource"
)

type AgentToCompanyMapperStorage struct {
	conn *datasource.TimescaleConnection
}

func NewAgentToCompanyMapperStorage(conn *datasource.TimescaleConnection) *AgentToCompanyMapperStorage {
	return &AgentToCompanyMapperStorage{
		conn: conn,
	}
}

func (s *AgentToCompanyMapperStorage) GetDisplayName(ctx context.Context, agentID string) (string, error) {
	row := s.conn.Pool().QueryRow(ctx,
		`SELECT display_name FROM agents WHERE id = $1`,
		agentID,
	)

	var displayName string

	if err := row.Scan(&displayName); err != nil {
		return "", err
	}

	return displayName, nil
}

func (s *AgentToCompanyMapperStorage) GetAll(ctx context.Context) (map[string]string, error) {
	sql := `SELECT id, holder FROM agents`

	rows, err := s.conn.Pool().Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mapping := make(map[string]string)

	for rows.Next() {
		var agentId string
		var holder string

		if err := rows.Scan(&agentId, &holder); err != nil {
			return nil, err
		}

		mapping[agentId] = holder
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return mapping, nil
}

func (s *AgentToCompanyMapperStorage) GetByAgent(ctx context.Context, agentID string) (string, error) {
	sql := `SELECT holder FROM agents WHERE id = $1`

	row := s.conn.Pool().QueryRow(ctx, sql, agentID)

	var holder string

	if err := row.Scan(&holder); err != nil {
		return "", err
	}

	return holder, nil
}
