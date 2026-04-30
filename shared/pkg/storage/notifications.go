package storage

import (
	"context"

	"zhacked.me/oxyl/shared/pkg/datasource"
	"zhacked.me/oxyl/shared/pkg/models"
)

type NotificationStorage struct {
	conn *datasource.TimescaleConnection
}

func NewNotificationStorage(conn *datasource.TimescaleConnection) *NotificationStorage {
	return &NotificationStorage{conn: conn}
}

func (s *NotificationStorage) Insert(ctx context.Context, identifier, agentID string, reason models.NotificationType, value string) error {
	_, err := s.conn.Pool().Exec(ctx,
		`INSERT INTO agent_notification_logs (identifier, agent, trigger_reason, trigger_value)
         VALUES ($1, $2, $3, $4)
         ON CONFLICT (identifier, sent_at) DO NOTHING`,
		identifier, agentID, reason, value,
	)
	return err
}

func (s *NotificationStorage) GetByAgent(ctx context.Context, agentID string) ([]*models.NotificationLog, error) {
	rows, err := s.conn.Pool().Query(ctx,
		`SELECT identifier, agent, trigger_reason, trigger_value, ack, failed, sent_at
         FROM agent_notification_logs
         WHERE agent = $1
         ORDER BY sent_at DESC`,
		agentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.NotificationLog
	for rows.Next() {
		l := new(models.NotificationLog)
		if err := rows.Scan(&l.Identifier, &l.Agent, &l.TriggerReason, &l.TriggerValue, &l.Ack, &l.Failed, &l.SentAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}

	return logs, rows.Err()
}

func (s *NotificationStorage) Ack(ctx context.Context, identifier string) error {
	_, err := s.conn.Pool().Exec(ctx,
		`UPDATE agent_notification_logs SET ack = true WHERE identifier = $1`,
		identifier,
	)
	return err
}

func (s *NotificationStorage) MarkFailed(ctx context.Context, identifier string) error {
	_, err := s.conn.Pool().Exec(ctx,
		`UPDATE agent_notification_logs SET failed = true WHERE identifier = $1`,
		identifier,
	)
	return err
}
