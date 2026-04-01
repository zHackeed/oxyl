package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
	"zhacked.me/oxyl/shared/pkg/models"
	"zhacked.me/oxyl/shared/pkg/storage"
	"zhacked.me/oxyl/shared/pkg/utils"
)

// Todo: mental way of handling a lot more stuff for the api on the service layer.
// Todo: consolidate errors on a same model class for consistency.

type UserService struct {
	userStorage *storage.UserStorage
}

func NewUserService(userStorage *storage.UserStorage) *UserService {
	return new(UserService{
		userStorage: userStorage,
	})
}

func (u *UserService) Authenticate(ctx context.Context, email, password string) (*models.User, error) {
	user, err := u.userStorage.GetUser(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("unable to find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := u.userStorage.UpdateLastLogin(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("unable to update last login: %w", err)
	}

	return user, nil
}

func (u *UserService) Register(ctx context.Context, name, surname, email, password string) (*models.User, error) {
	user, err := models.NewUser(name, surname, email, password)
	if err != nil {
		return nil, fmt.Errorf("unable to create user: %w", err)
	}

	if err := u.userStorage.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("unable to create user: %w", err)
	}

	slog.Info("[UserService] user created", "user_id", user.ID)

	return user, nil
}

// UpdatePassword should force a logout and invalidation of all sessions.
func (u *UserService) UpdatePassword(ctx context.Context, newPassword string) error {
	userId, found := utils.GetValueFromContext[string](ctx, models.ContextKeyUser)
	if !found {
		return errors.New("user not found in context")
	}

	user, err := u.userStorage.GetUserById(ctx, userId)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPassword)); err == nil {
		return errors.New("passwords are the same")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("unable to hash password: %w", err)
	}

	slog.Info("[UserService] password updated", "user_id", userId)

	return u.userStorage.UpdatePassword(ctx, userId, string(hashed))

}

func (u *UserService) GetUser(ctx context.Context) (*models.User, error) {
	userId, found := utils.GetValueFromContext[string](ctx, models.ContextKeyUser)
	if !found {
		return nil, errors.New("user not found in context")
	}

	slog.Info("[UserService] retrieving user", "user_id", userId)

	return u.userStorage.GetUserById(ctx, userId)
}
