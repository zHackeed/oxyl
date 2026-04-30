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
	ErrCompanyNotFound              = errors.New("company not found")
	ErrMemberAlreadyExists          = errors.New("member already exists")
	ErrMemberNotFound               = errors.New("member not found")
	ErrNotificationEndpointNotFound = errors.New("notification endpoint not found")
	ErrDiscordChannelRequired       = errors.New("channel is required for Discord")
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

	//TODO: Define permission logic adapt
	if _, err := tx.Exec(ctx, sql, company.Holder, company.ID, models.CompanyPermissionOwner); err != nil {
		return fmt.Errorf("unable to add holder as member: %w", err)
	}

	sql = `INSERT INTO company_notification_thresholds (holder, notification_type, value) VALUES ($1, $2, $3)`
	for notificationType, threshold := range company.NotificationThresholds {
		if _, err := tx.Exec(ctx, sql, company.ID, notificationType, threshold); err != nil {
			return fmt.Errorf("unable to create company notification thresholds: %w", err)
		}
	}

	sql = `INSERT INTO company_notification_settings (id, holder, webhook_type, endpoint, channel) VALUES ($1, $2, $3, $4, $5)`
	for _, endpoint := range company.NotificationEndpoints {
		if _, err := tx.Exec(ctx, sql, endpoint.ID, company.ID, endpoint.WebhookType, endpoint.Endpoint, endpoint.Channel); err != nil {
			return fmt.Errorf("unable to create company notification settings: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}

	return nil
}

func (c *CompanyStorage) GetCountOfAgentsInCompany(ctx context.Context, companyID string) (int, error) {
	sql := `SELECT COUNT(id) FROM AGENTS WHERE holder = $1`

	row := c.conn.Pool().QueryRow(ctx, sql, companyID)

	var count int

	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("unable to get count of agents in company: %w", err)
	}

	return count, nil
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

	nodes, err := c.GetCountOfAgentsInCompany(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("unable to get count of agents in company: %w", err)
	}

	company.Nodes = nodes

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

		nodes, err := c.GetCountOfAgentsInCompany(ctx, company.ID)
		if err != nil {
			return nil, fmt.Errorf("unable to get count of agents in company: %w", err)
		}

		company.Nodes = nodes

		companies = append(companies, &company)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("unable to get companies of user: %w", rows.Err())
	}

	return companies, nil
}

func (c *CompanyStorage) GetCompanyMembership(ctx context.Context, userID, companyID string) (*models.CompanyMember, error) {
	sql := `SELECT permission_bitwise FROM company_members WHERE user_id = $1 AND company_id = $2`
	row := c.conn.Pool().QueryRow(ctx, sql, userID, companyID)

	var permission int
	if err := row.Scan(&permission); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMemberNotFound
		}

		return nil, fmt.Errorf("unable to get company membership: %w", err)
	}

	return &models.CompanyMember{
		UserID:     userID,
		Permission: models.CompanyPermission(permission),
	}, nil
}

func (c *CompanyStorage) AddMemberToCompany(ctx context.Context, userId, companyID string, permission int) error {
	sql := `INSERT INTO company_members (user_id, company_id, permission_bitwise) VALUES ($1, $2, $3)`
	if _, err := c.conn.Pool().Exec(ctx, sql, userId, companyID, permission); err != nil {

		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23505" { // Unique constraint violation
			return ErrMemberAlreadyExists
		}

		return fmt.Errorf("unable to add member to company: %w", err)
	}

	return nil
}

func (c *CompanyStorage) GetMembersOfCompany(ctx context.Context, companyID string) ([]*models.CompanyMemberComposite, error) {
	sql := `SELECT u.name, u.surname, u.email, cm.permission_bitwise, cm.created_at FROM company_members as cm INNER JOIN users u on cm.user_id = u.id WHERE cm.company_id=$1`

	rows, err := c.conn.Pool().Query(ctx, sql, companyID)
	if err != nil {
		return nil, fmt.Errorf("unable to get company members: %w", err)
	}

	members := make([]*models.CompanyMemberComposite, 0)

	for rows.Next() {
		member := new(models.CompanyMemberComposite)

		if err := rows.Scan(&member.User.Name, &member.User.Surname, &member.User.Email, &member.Permissions, &member.CreatedAt); err != nil {
			return nil, errors.New("failed parsing user")
		}

		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read from tables %w", err)
	}

	return members, nil
}

func (c *CompanyStorage) RemoveMemberFromCompany(ctx context.Context, companyID, userID string) error {
	sql := `DELETE FROM company_members WHERE user_id = $1 AND company_id = $2`
	if _, err := c.conn.Pool().Exec(ctx, sql, userID, companyID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrMemberNotFound
		}

		return fmt.Errorf("unable to remove member from company: %w", err)
	}

	return nil
}

// Todo: rethink this

func (c *CompanyStorage) GetNotificationThresholdsOfCompany(ctx context.Context, companyID string) (map[models.NotificationType]int, error) {
	sql := `SELECT notification_type, value FROM company_notification_thresholds WHERE holder = $1`
	rows, err := c.conn.Pool().Query(ctx, sql, companyID)
	if err != nil {
		return nil, fmt.Errorf("unable to get company notification thresholds: %w", err)
	}

	thresholds := make(map[models.NotificationType]int)
	for rows.Next() {
		var notificationType models.NotificationType
		var value int
		if err := rows.Scan(&notificationType, &value); err != nil {
			return nil, fmt.Errorf("unable to get company notification thresholds: %w", err)
		}
		thresholds[notificationType] = value
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("unable to get company notification thresholds: %w", rows.Err())
	}

	return thresholds, nil
}
func (c *CompanyStorage) SetNotificationThresholdForCompany(ctx context.Context, companyID string, notificationType models.NotificationType, threshold int) error {
	sql := `UPDATE company_notification_thresholds SET value = $3 WHERE holder = $1 AND notification_type = $2`
	if _, err := c.conn.Pool().Exec(ctx, sql, companyID, notificationType, threshold); err != nil {
		return fmt.Errorf("unable to set notification threshold for company: %w", err)
	}

	return nil
}
func (c *CompanyStorage) DeleteCompany(ctx context.Context, companyID string) error {
	sql := `DELETE FROM companies WHERE id = $1`
	if _, err := c.conn.Pool().Exec(ctx, sql, companyID); err != nil {
		return fmt.Errorf("unable to delete company: %w", err)
	}

	return nil
}

// -- Notifications --
func (c *CompanyStorage) GetNotificationEndpointsOfCompany(ctx context.Context, companyID string) ([]*models.CompanyNotificationSettings, error) {
	sql := `SELECT id, holder, webhook_type, endpoint, channel FROM company_notification_settings WHERE holder = $1`
	rows, err := c.conn.Pool().Query(ctx, sql, companyID)
	if err != nil {
		return nil, fmt.Errorf("unable to get company notification endpoints: %w", err)
	}

	defer rows.Close()

	var endpoints []*models.CompanyNotificationSettings
	for rows.Next() {
		endpoint := new(models.CompanyNotificationSettings)
		if err := rows.Scan(&endpoint.ID, &endpoint.Holder, &endpoint.WebhookType, &endpoint.Endpoint, &endpoint.Channel); err != nil {
			return nil, fmt.Errorf("unable to scan company notification endpoint: %w", err)
		}
		endpoints = append(endpoints, endpoint)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("unable to get company notification endpoints: %w", err)
	}

	return endpoints, nil
}

func (c *CompanyStorage) GetNotificationEndpoint(ctx context.Context, id, companyID string) (*models.CompanyNotificationSettings, error) {
	sql := `SELECT id, holder, webhook_type, endpoint, channel FROM company_notification_settings WHERE id = $1 AND holder = $2`

	endpoint := new(models.CompanyNotificationSettings)
	err := c.conn.Pool().QueryRow(ctx, sql, id, companyID).Scan(
		&endpoint.ID, &endpoint.Holder, &endpoint.WebhookType, &endpoint.Endpoint, &endpoint.Channel,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotificationEndpointNotFound
		}
		return nil, fmt.Errorf("unable to get notification endpoint: %w", err)
	}

	return endpoint, nil
}

func (c *CompanyStorage) AddNotificationEndpointToCompany(ctx context.Context, companyID string, endpoint *models.CompanyNotificationSettings) error {
	if endpoint.WebhookType == models.WebhookTypeDiscord && (endpoint.Channel == nil || *endpoint.Channel == "") {
		return ErrDiscordChannelRequired
	}

	sql := `INSERT INTO company_notification_settings (id, holder, webhook_type, endpoint, channel) VALUES ($1, $2, $3, $4, $5)`
	if _, err := c.conn.Pool().Exec(ctx, sql, endpoint.ID, companyID, endpoint.WebhookType, endpoint.Endpoint, endpoint.Channel); err != nil {
		return fmt.Errorf("unable to add notification endpoint to company: %w", err)
	}

	return nil
}

func (c *CompanyStorage) UpdateNotificationEndpoint(ctx context.Context, companyID string, endpoint *models.CompanyNotificationSettings) error {
	if endpoint.WebhookType == models.WebhookTypeDiscord && (endpoint.Channel == nil || *endpoint.Channel == "") {
		return ErrDiscordChannelRequired
	}

	sql := `UPDATE company_notification_settings SET webhook_type = $1, endpoint = $2, channel = $3 WHERE id = $4 AND holder = $5`
	if _, err := c.conn.Pool().Exec(ctx, sql, endpoint.WebhookType, endpoint.Endpoint, endpoint.Channel, endpoint.ID, companyID); err != nil {
		return fmt.Errorf("unable to update notification endpoint: %w", err)
	}

	return nil
}

func (c *CompanyStorage) RemoveNotificationEndpointFromCompany(ctx context.Context, companyID, id string) error {
	sql := `DELETE FROM company_notification_settings WHERE id = $1 AND holder = $2`
	if _, err := c.conn.Pool().Exec(ctx, sql, id, companyID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotificationEndpointNotFound
		}
		return fmt.Errorf("unable to remove notification endpoint from company: %w", err)
	}

	return nil
}
