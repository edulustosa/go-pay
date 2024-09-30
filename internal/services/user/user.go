package user

import (
	"context"
	"errors"

	"github.com/edulustosa/go-pay/internal/database/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (models.User, error)
	FindByDocument(ctx context.Context, document string) (models.User, error)
	FindByEmail(ctx context.Context, email string) (models.User, error)
	Create(ctx context.Context, user models.User) (uuid.UUID, error)
	UpdateBalance(ctx context.Context, id uuid.UUID, balance float64) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo,
	}
}

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

func (s *UserService) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (models.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) Create(
	ctx context.Context,
	user models.User,
) (uuid.UUID, error) {
	errChan := make(chan error, 2)

	go func() {
		_, err := s.repo.FindByDocument(ctx, user.Document)
		errChan <- err
	}()

	go func() {
		_, err := s.repo.FindByEmail(ctx, user.Email)
		errChan <- err
	}()

	errEmail, errDocument := <-errChan, <-errChan
	if errEmail == nil || errDocument == nil {
		return uuid.Nil, ErrUserAlreadyExists
	}

	return s.repo.Create(ctx, user)
}

func (s *UserService) UpdateBalance(
	ctx context.Context,
	id uuid.UUID,
	balance float64,
) error {
	return s.repo.UpdateBalance(ctx, id, balance)
}
