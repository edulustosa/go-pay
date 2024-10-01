package user_test

import (
	"context"
	"testing"

	"github.com/edulustosa/go-pay/internal/database/models"
	"github.com/edulustosa/go-pay/internal/database/repo"
	"github.com/edulustosa/go-pay/internal/services/user"
)

func TestUserService_Create(t *testing.T) {
	userRepository := repo.InMemoryUserRepository{}
	sut := user.NewUserService(&userRepository)

	ctx := context.Background()
	t.Run("should be able to create a new user", func(t *testing.T) {
		user := models.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@email.com",
			Document:  "12345678900",
		}

		userID, err := sut.Create(ctx, user)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		createdUser, err := userRepository.FindByID(ctx, userID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if createdUser.FirstName != user.FirstName {
			t.Errorf("expected %s, got %s", user.FirstName, createdUser.FirstName)
		}

		t.Logf("created user: %+v", createdUser)
	})

	t.Run("should not be able to create a user with the same document", func(t *testing.T) {
		userRepository.Users = []models.User{}

		user1 := models.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@email.com",
			Document:  "12345678900",
		}

		_, err := sut.Create(ctx, user1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		user2 := models.User{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "janedoe@email.com",
			Document:  "12345678900",
		}

		_, err = sut.Create(ctx, user2)
		if err != user.ErrUserAlreadyExists {
			t.Errorf("expected %v, got %v", user.ErrUserAlreadyExists, err)
		}

		t.Logf("error: %v", err)
	})

	t.Run("should not be able to create a user with the same email", func(t *testing.T) {
		userRepository.Users = []models.User{}

		user1 := models.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@email.com",
			Document:  "12345678900",
		}

		_, err := sut.Create(ctx, user1)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		user2 := models.User{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "johndoe@email.com",
			Document:  "12345678901",
		}

		_, err = sut.Create(ctx, user2)
		if err != user.ErrUserAlreadyExists {
			t.Errorf("expected %v, got %v", user.ErrUserAlreadyExists, err)
		}

		t.Logf("error: %v", err)
	})
}
