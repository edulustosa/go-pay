package user

import (
	"context"
	"errors"
	"strings"

	"github.com/edulustosa/go-pay/internal/database/models"
	"github.com/edulustosa/go-pay/internal/dtos"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (models.User, error)
	FindByDocument(ctx context.Context, document string) (models.User, error)
	FindByEmail(ctx context.Context, email string) (models.User, error)
	Create(ctx context.Context, user models.User) (uuid.UUID, error)
	UpdateBalance(ctx context.Context, id uuid.UUID, balance float64) error
	FindMany(ctx context.Context, page int) ([]models.User, error)
}

type Service struct {
	repo userRepository
}

func NewService(repo userRepository) *Service {
	return &Service{
		repo,
	}
}

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

func (s *Service) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (models.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Create(
	ctx context.Context,
	userDTO dtos.UserDTO,
) (uuid.UUID, error) {
	userDTO.Document = normalizeDocument(userDTO.Document)

	_, err := s.repo.FindByDocument(ctx, userDTO.Document)
	if err == nil {
		return uuid.Nil, ErrUserAlreadyExists
	}

	_, err = s.repo.FindByEmail(ctx, userDTO.Email)
	if err == nil {
		return uuid.Nil, ErrUserAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(userDTO.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return uuid.Nil, err
	}

	if userDTO.Role == "" {
		userDTO.Role = models.RoleCommon
	}

	user := models.User{
		FirstName:    userDTO.FirstName,
		LastName:     userDTO.LastName,
		Document:     userDTO.Document,
		Email:        userDTO.Email,
		PasswordHash: string(passwordHash),
		Balance:      userDTO.Balance,
		Role:         userDTO.Role,
	}

	return s.repo.Create(ctx, user)
}

func normalizeDocument(document string) string {
	document = strings.ReplaceAll(document, ".", "")
	document = strings.ReplaceAll(document, "-", "")
	return document
}

func (s *Service) UpdateBalance(
	ctx context.Context,
	id uuid.UUID,
	balance float64,
) error {
	return s.repo.UpdateBalance(ctx, id, balance)
}

func (s *Service) FindMany(
	ctx context.Context,
	page int,
) ([]dtos.UserResponseDTO, error) {
	if page < 1 {
		page = 1
	}

	users, err := s.repo.FindMany(ctx, page)
	if err != nil {
		return nil, err
	}

	usersDTO := make([]dtos.UserResponseDTO, len(users))
	for i, user := range users {
		usersDTO[i] = dtos.UserResponseDTO{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Document:  user.Document,
			Email:     user.Email,
			Balance:   user.Balance,
			Role:      user.Role,
		}
	}

	return usersDTO, nil
}
