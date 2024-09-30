package repo

import (
	"context"
	"errors"
	"time"

	"github.com/edulustosa/go-pay/internal/database/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type InMemoryUserRepository struct {
	Users []models.User
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInsertionOnUnique = errors.New("insertion on unique field")
)

func (r *InMemoryUserRepository) FindByDocument(
	_ context.Context,
	document string,
) (models.User, error) {
	for _, user := range r.Users {
		if user.Document == document {
			return user, nil
		}
	}

	return models.User{}, ErrUserNotFound
}

func (r *InMemoryUserRepository) FindByID(
	_ context.Context,
	id uuid.UUID,
) (models.User, error) {
	for _, user := range r.Users {
		if user.ID == id {
			return user, nil
		}
	}

	return models.User{}, ErrUserNotFound
}

func (r *InMemoryUserRepository) FindByEmail(
	_ context.Context,
	email string,
) (models.User, error) {
	for _, user := range r.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return models.User{}, ErrUserNotFound
}

func (r *InMemoryUserRepository) Create(
	ctx context.Context,
	user models.User,
) (uuid.UUID, error) {
	_, err := r.FindByDocument(ctx, user.Document)
	if err == nil {
		return uuid.Nil, ErrInsertionOnUnique
	}

	_, err = r.FindByEmail(ctx, user.Email)
	if err == nil {
		return uuid.Nil, ErrInsertionOnUnique
	}

	user.ID = uuid.New()
	if user.Role == "" {
		user.Role = models.RoleCommon
	}
	user.CreatedAt = pgtype.Timestamp{Time: time.Now()}
	user.UpdatedAt = pgtype.Timestamp{Time: time.Now()}

	r.Users = append(r.Users, user)
	return user.ID, nil
}

func (r *InMemoryUserRepository) UpdateBalance(
	_ context.Context,
	id uuid.UUID,
	balance float64,
) error {
	for i, user := range r.Users {
		if user.ID == id {
			r.Users[i].Balance = balance
			return nil
		}
	}

	return ErrUserNotFound
}
