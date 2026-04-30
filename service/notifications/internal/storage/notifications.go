package storage

import (
	"context"

	"zhacked.me/oxyl/shared/pkg/datasource"
	comm "zhacked.me/oxyl/shared/pkg/models"
)

type NotificationSettingStorage struct {
	conn *datasource.TimescaleConnection
}

func NewNotificationSettingStorage(conn *datasource.TimescaleConnection) *NotificationSettingStorage {
	return &NotificationSettingStorage{
		conn: conn,
	}
}

func (s *NotificationSettingStorage) GetAll(ctx context.Context) (map[string][]*comm.CompanyNotificationSettings, error) {
	rows, err := s.conn.Pool().Query(ctx,
		`SELECT id, holder, webhook_type, endpoint, channel FROM company_notification_settings`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settings := make(map[string][]*comm.CompanyNotificationSettings)

	for rows.Next() {
		s := new(comm.CompanyNotificationSettings)
		if err := rows.Scan(&s.ID, &s.Holder, &s.WebhookType, &s.Endpoint, &s.Channel); err != nil {
			return nil, err
		}
		settings[s.Holder] = append(settings[s.Holder], s)
	}

	return settings, rows.Err()
}

func (s *NotificationSettingStorage) GetById(ctx context.Context, id string) (*comm.CompanyNotificationSettings, error) {
	rows := s.conn.Pool().QueryRow(ctx,
		`SELECT id, holder, webhook_type, endpoint, channel FROM company_notification_settings WHERE id = $1`,
	)

	setting := new(comm.CompanyNotificationSettings)

	if err := rows.Scan(&setting.ID, &setting.Holder, &setting.WebhookType, &setting.Endpoint, &setting.Channel); err != nil {
		return nil, err
	}

	return setting, nil
}

func (s *NotificationSettingStorage) GetByCompany(ctx context.Context, companyID string) ([]*comm.CompanyNotificationSettings, error) {
	rows, err := s.conn.Pool().Query(ctx,
		`SELECT id, webhook_type, endpoint, channel FROM company_notification_settings WHERE holder = $1`,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settings := make([]*comm.CompanyNotificationSettings, 0)

	for rows.Next() {
		s := new(comm.CompanyNotificationSettings)
		if err := rows.Scan(&s.ID, &s.WebhookType, &s.Endpoint, &s.Channel); err != nil {
			return nil, err
		}

		s.Holder = companyID
		settings = append(settings, s)
	}

	return settings, rows.Err()
}
