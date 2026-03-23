package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"zhacked.me/oxyl/shared/pkg/datasource"
	"zhacked.me/oxyl/shared/pkg/models"
)

var (
	ErrCompanyNotFound     = errors.New("company not found")
	ErrMemberAlreadyExists = errors.New("member already exists")
)

type CompanyStorage struct {
	conn *datasource.TimescaleConnection
}

func NewCompanyStorage(persistence *datasource.TimescaleConnection) *CompanyStorage {
	return &CompanyStorage{
		conn: persistence,
	}
}

func (c *CompanyStorage) CreateCompany(ctx context.Context, company *models.Company) error {
	tx, err := c.conn.BeginTx(ctx)

	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	sql := `INSERT INTO companies (id, display_name, holder) VALUES ($1, $2, $3)`
	if _, err := tx.Exec(ctx, sql, company.ID, company.DisplayName, company.Holder); err != nil {
		return fmt.Errorf("unable to create company: %w", err)
	}

	sql = `INSERT INTO company_members (user_id, company_id, permission_bitwise) VALUES ($1, $2, $3)`

	//TODO: Define permission logic handler
	if _, err := tx.Exec(ctx, sql, company.Holder, company.ID, models.CompanyPermission(999)); err != nil {
		return fmt.Errorf("unable to add holder as member: %w", err)
	}

	for notificationType, threshold := range company.NotificationThresholds {
		sql := `INSERT INTO company_notification_thresholds (holder, notification_type, value) VALUES ($1, $2, $3)`
		if _, err := tx.Exec(ctx, sql, company.ID, notificationType, threshold); err != nil {
			return fmt.Errorf("unable to create company notification thresholds: %w", err)
		}
	}

	for _, endpoint := range company.NotificationEndpoints {
		sql := `INSERT INTO company_notification_settings (holder, webhook_type, endpoint, metakeys) VALUES ($1, $2, $3, $4)`
		if _, err := tx.Exec(ctx, sql, company.ID, endpoint.WebhookType, endpoint.EndpointUrl, endpoint.MetaKeys); err != nil {
			return fmt.Errorf("unable to create company notification settings: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}

	return nil
}

func (c *CompanyStorage) GetCompany(ctx context.Context, companyID string) (*models.Company, error) {
	sql := `SELECT display_name, holder, limit_nodes, enabled, created_at FROM companies WHERE id = $1`
	row := c.conn.Pool().QueryRow(ctx, sql, companyID)
	var company models.Company

	if err := row.Scan(&company.DisplayName, &company.Holder, &company.LimitNodes, &company.Enabled, &company.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCompanyNotFound
		}
		return nil, fmt.Errorf("unable to get company: %w", err)
	}

	company.ID = companyID

	// Avoid allocating nil maps, panics
	company.Members = make(map[string]*models.CompanyMember)
	company.NotificationEndpoints = make([]*models.CompanyNotificationSettings, 0)
	company.NotificationThresholds = make(map[models.NotificationType]int)

	return &company, nil
}

func (c *CompanyStorage) GetCompanies(ctx context.Context) ([]*models.Company, error) {
	sql := `SELECT id, display_name, holder, limit_nodes, enabled, created_at FROM companies`
	rows, err := c.conn.Pool().Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("unable to get companies: %w", err)
	}

	companies := make([]*models.Company, 0)
	defer rows.Close()

	for rows.Next() {
		var company models.Company
		if err := rows.Scan(&company.ID, &company.DisplayName, &company.Holder, &company.LimitNodes, &company.Enabled, &company.CreatedAt); err != nil {
			return nil, fmt.Errorf("unable to get companies: %w", err)
		}
		companies = append(companies, &company)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("unable to get companies: %w", rows.Err())
	}

	return companies, nil
}

func (c *CompanyStorage) GetMembersOfCompany(ctx context.Context, companyID string) ([]*models.CompanyMember, error) {
	sql := `SELECT user_id, permission_bitwise FROM company_members WHERE company_id = $1`
	rows, err := c.conn.Pool().Query(ctx, sql, companyID)
	if err != nil {
		return nil, fmt.Errorf("unable to get company members: %w", err)
	}

	defer rows.Close()
	var members []*models.CompanyMember
	for rows.Next() {
		var member models.CompanyMember
		if err := rows.Scan(&member.UserID, &member.Permission); err != nil {
			return nil, fmt.Errorf("unable to get company members: %w", err)
		}
		members = append(members, &member)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("unable to get company members: %w", rows.Err())
	}

	return members, nil
}

func (c *CompanyStorage) GetCompaniesOfUser(ctx context.Context, userID string) ([]*models.Company, error) {
	sql := `SELECT c.id, c.display_name, c.holder, c.limit_nodes, c.enabled, c.created_at FROM companies c INNER JOIN company_members cm ON c.id = cm.company_id WHERE cm.user_id = $1`

	rows, err := c.conn.Pool().Query(ctx, sql, userID)

	if err != nil {
		return nil, fmt.Errorf("unable to get companies of user: %w", err)
	}

	defer rows.Close()
	var companies []*models.Company
	for rows.Next() {
		var company models.Company
		if err := rows.Scan(&company.ID, &company.DisplayName, &company.Holder, &company.LimitNodes, &company.Enabled, &company.CreatedAt); err != nil {
			return nil, fmt.Errorf("unable to get companies of user: %w", err)
		}
		companies = append(companies, &company)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("unable to get companies of user: %w", rows.Err())
	}

	return companies, nil
}

func (c *CompanyStorage) GetNotificationEndpointsOfCompany(ctx context.Context, companyID string) ([]*models.CompanyNotificationSettings, error) {
	sql := `SELECT webhook_type, endpoint, metakeys FROM company_notification_settings WHERE holder = $1`
	rows, err := c.conn.Pool().Query(ctx, sql, companyID)
	if err != nil {
		return nil, fmt.Errorf("unable to get company notification endpoints: %w", err)
	}

	defer rows.Close()
	var endpoints []*models.CompanyNotificationSettings
	for rows.Next() {
		var endpoint models.CompanyNotificationSettings
		if err := rows.Scan(&endpoint.WebhookType, &endpoint.EndpointUrl, &endpoint.MetaKeys); err != nil {
			return nil, fmt.Errorf("unable to get company notification endpoints: %w", err)
		}
		endpoints = append(endpoints, &endpoint)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("unable to get company notification endpoints: %w", rows.Err())
	}

	return endpoints, nil
}

func (c *CompanyStorage) AddMemberToCompany(ctx context.Context, companyID, userID string, permission int) error {
	sql := `INSERT INTO company_members (user_id, company_id, permission_bitwise) VALUES ($1, $2, $3)`
	if _, err := c.conn.Pool().Exec(ctx, sql, userID, companyID, permission); err != nil {

		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23505" { // Unique constraint violation
			return ErrMemberAlreadyExists
		}

		return fmt.Errorf("unable to add member to company: %w", err)
	}

	return nil
}

func (c *CompanyStorage) RemoveMemberFromCompany(ctx context.Context, companyID, userID string) error {
	sql := `DELETE FROM company_members WHERE user_id = $1 AND company_id = $2`
	if _, err := c.conn.Pool().Exec(ctx, sql, userID, companyID); err != nil {
		return fmt.Errorf("unable to remove member from company: %w", err)
	}

	return nil
}

// Todo: rethink this

func (c *CompanyStorage) AddNotificationEndpointToCompany(ctx context.Context, companyID string, endpoint *models.CompanyNotificationSettings) error {
	sql := `INSERT INTO company_notification_settings (holder, webhook_type, endpoint, metakeys) VALUES ($1, $2, $3, $4)`
	if _, err := c.conn.Pool().Exec(ctx, sql, companyID, endpoint.WebhookType, endpoint.EndpointUrl, endpoint.MetaKeys); err != nil {
		return fmt.Errorf("unable to add notification endpoint to company: %w", err)
	}

	return nil
}

func (c *CompanyStorage) RemoveNotificationEndpointFromCompany(ctx context.Context, companyID string, endpoint *models.CompanyNotificationSettings) error {
	sql := `DELETE FROM company_notification_settings WHERE holder = $1 AND endpoint = $2`
	if _, err := c.conn.Pool().Exec(ctx, sql, companyID, endpoint.EndpointUrl); err != nil {
		return fmt.Errorf("unable to remove notification endpoint from company: %w", err)
	}

	return nil
}

func (c *CompanyStorage) SetNotificationThresholdForCompany(ctx context.Context, companyID string, notificationType models.NotificationType, threshold int) error {
	sql := `INSERT INTO company_notification_thresholds (holder, notification_type, value) VALUES ($1, $2, $3) ON CONFLICT (holder, notification_type) DO UPDATE SET value = $3`
	if _, err := c.conn.Pool().Exec(ctx, sql, companyID, notificationType, threshold); err != nil {
		return fmt.Errorf("unable to set notification threshold for company: %w", err)
	}

	return nil
}

// Todo: might use enable bool?
func (c *CompanyStorage) DeleteCompany(ctx context.Context, companyID string) error {
	sql := `DELETE FROM companies WHERE id = $1`
	if _, err := c.conn.Pool().Exec(ctx, sql, companyID); err != nil {
		return fmt.Errorf("unable to delete company: %w", err)
	}

	return nil
}
