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
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type UserStorage struct {
	conn *datasource.TimescaleConnection
}

func NewUserStorage(persistence *datasource.TimescaleConnection) *UserStorage {
	return &UserStorage{
		conn: persistence,
	}
}

func (u *UserStorage) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, email, password, name, surname) VALUES ($1, $2, $3, $4, $5)`

	if _, err := u.conn.Pool().Exec(ctx, query,
		user.ID, user.Email,
		user.Password,
		user.Name, user.Surname); err != nil {

		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23505" { // Unique constraint violation
			return ErrUserAlreadyExists
		}

		return fmt.Errorf("unable to create user: %w", err)
	}

	return nil
}

func (u *UserStorage) GetIdFromEmail(ctx context.Context, email string) (string, error) {
	sql := `SELECT id FROM users WHERE email = $1`
	row := u.conn.Pool().QueryRow(ctx, sql, email)
	var userId string
	if err := row.Scan(&userId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrUserNotFound
		}

		return "", fmt.Errorf("unable to get user: %w", err)
	}
	return userId, nil
}

func (u *UserStorage) GetUser(ctx context.Context, email string) (*models.User, error) {
	sql := `SELECT id, email, password, name, surname, enabled, last_login, created_at FROM users WHERE email = $1`

	row := u.conn.Pool().QueryRow(ctx, sql, email)

	var user models.User
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Surname, &user.Enabled, &user.LastLogin, &user.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("unable to get user: %w", err)
	}

	return &user, nil
}

func (u *UserStorage) GetUserById(ctx context.Context, userId string) (*models.User, error) {
	sql := `SELECT id, email, password, name, surname, enabled, last_login, created_at FROM users WHERE id = $1`
	row := u.conn.Pool().QueryRow(ctx, sql, userId)
	var user models.User
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Surname, &user.Enabled, &user.LastLogin, &user.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("unable to get user: %w", err)
	}

	return &user, nil
}

func (u *UserStorage) GetUserWithMemberships(ctx context.Context, email string) (*models.User, error) {
	user, err := u.GetUser(ctx, email)
	if err != nil {
		//propagate error
		return nil, err
	}

	sql := `SELECT company_id, permission_bitwise FROM company_members WHERE user_id = $1`
	rows, err := u.conn.Pool().Query(ctx, sql, user.ID)

	if err != nil {
		return nil, fmt.Errorf("unable to get company <-> user memberships: %w", err)
	}

	defer rows.Close()
	memberOf := make(map[string]*models.CompanyMember)

	for rows.Next() {
		var companyID string
		var permission int
		if err := rows.Scan(&companyID, &permission); err != nil {
			return nil, fmt.Errorf("unable to get user: %w", err)
		}
		memberOf[companyID] = &models.CompanyMember{
			Permission: models.CompanyPermission(permission),
		}
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("unable to get user: %w", rows.Err())
	}

	user.MemberOf = memberOf

	return user, nil
}

func (u *UserStorage) GetUsers(ctx context.Context) ([]*models.User, error) {
	sql := `SELECT id, email, name, surname, enabled, last_login, created_at FROM users`
	rows, err := u.conn.Pool().Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("unable to get users: %w", err)
	}

	var users []*models.User

	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Surname, &user.Enabled, &user.LastLogin, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("unable to get users: %w", err)
		}

		users = append(users, &user)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("unable to get users: %w", rows.Err())
	}

	return users, nil
}

func (u *UserStorage) UpdatePassword(ctx context.Context, userId, newPassword string) error {
	sql := `UPDATE users SET password = $1, last_updated = CURRENT_TIMESTAMP WHERE id = $2`
	if _, err := u.conn.Pool().Exec(ctx, sql, newPassword, userId); err != nil {
		return fmt.Errorf("unable to update user password: %w", err)
	}

	return nil
}

func (u *UserStorage) UpdateState(ctx context.Context, userId string, enabled bool) error {
	sql := `UPDATE users SET enabled = $1, last_updated = CURRENT_TIMESTAMP WHERE id = $2`
	if _, err := u.conn.Pool().Exec(ctx, sql, enabled, userId); err != nil {
		return fmt.Errorf("unable to update user state: %w", err)
	}

	return nil
}

func (u *UserStorage) UpdateLastLogin(ctx context.Context, userId string) error {
	sql := `UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = $1`
	if _, err := u.conn.Pool().Exec(ctx, sql, userId); err != nil {
		return fmt.Errorf("unable to update user last login: %w", err)
	}

	return nil
}

func (u *UserStorage) UpdateNameAndLastName(ctx context.Context, userId, name, surname string) error {
	sql := `UPDATE users SET name = $1, surname = $2, last_updated = CURRENT_TIMESTAMP WHERE id = $3`
	if _, err := u.conn.Pool().Exec(ctx, sql, name, surname, userId); err != nil {
		return fmt.Errorf("unable to update user name: %w", err)
	}

	return nil
}

func (u *UserStorage) DeleteUser(ctx context.Context, userId string) error {
	sql := `DELETE FROM users WHERE id = $1`
	if _, err := u.conn.Pool().Exec(ctx, sql, userId); err != nil {
		return fmt.Errorf("unable to delete user: %w", err)
	}

	return nil
}
