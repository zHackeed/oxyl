package storage

import (
	"context"
	"fmt"

	"zhacked.me/oxyl/shared/pkg/datasource"
	comm "zhacked.me/oxyl/shared/pkg/models"
)

type NotificationStorage struct {
	conn *datasource.TimescaleConnection
}

func NewNotificationStorage(conn *datasource.TimescaleConnection) *NotificationStorage {
	return &NotificationStorage{conn: conn}
}

func (s *NotificationStorage) Insert(ctx context.Context, identifier, agentID string, reason comm.NotificationType, value string) error {
	_, err := s.conn.Pool().Exec(ctx,
		`INSERT INTO agent_notification_logs (identifier, agent, trigger_reason, trigger_value)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (identifier, sent_at) DO NOTHING`,
		identifier, agentID, reason, value,
	)
	if err != nil {
		return fmt.Errorf("unable to insert notification log: %w", err)
	}

	return nil
}

func (s *NotificationStorage) Ack(ctx context.Context, identifier string) error {
	_, err := s.conn.Pool().Exec(ctx,
		`UPDATE agent_notification_logs SET ack = true WHERE identifier = $1`,
		identifier,
	)
	if err != nil {
		return fmt.Errorf("unable to ack notification log: %w", err)
	}

	return nil
}

func (s *NotificationStorage) MarkFailed(ctx context.Context, identifier string) error {
	_, err := s.conn.Pool().Exec(ctx,
		`UPDATE agent_notification_logs SET failed = true WHERE identifier = $1`,
		identifier,
	)
	if err != nil {
		return fmt.Errorf("unable to mark notification log as failed: %w", err)
	}

	return nil
}
