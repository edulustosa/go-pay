package transfer_test

import (
	"context"
	"testing"

	"github.com/edulustosa/go-pay/helpers"
	"github.com/edulustosa/go-pay/internal/database/models"
	"github.com/edulustosa/go-pay/internal/database/repo"
	"github.com/edulustosa/go-pay/internal/dtos"
	"github.com/edulustosa/go-pay/internal/services/transfer"
	"github.com/edulustosa/go-pay/internal/services/user"
)

func TestTransferService(t *testing.T) {
	transactionsRepository := &repo.InMemoryTransactionsRepository{}
	userRepository := &repo.InMemoryUserRepository{}
	userService := user.NewUserService(userRepository)
	sut := transfer.NewTransferService(transactionsRepository, userService)

	ctx := context.Background()
	t.Run("should be able to make a transfer between users", func(t *testing.T) {
		user1, _ := userRepository.Create(ctx, models.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@email.com",
			Document:  "12345678900",
			Balance:   1000,
		})
		user2, _ := userRepository.Create(ctx, models.User{
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "janedoe@email.com",
			Document:  "09876543211",
			Balance:   500,
		})

		transaction := dtos.TransactionDTO{
			Value: 100,
			Payer: user1,
			Payee: user2,
		}

		transactionID, err := sut.NewTransaction(ctx, transaction)
		if err != nil {
			if err == transfer.ErrTransactionNotAuthorized {
				t.Fatal("transaction not authorized")
			}

			t.Fatalf("expected no error, got %v", err)
		}

		t.Logf("transactionID: %v", transactionID)

		user1Model, _ := userRepository.FindByID(ctx, user1)
		user2Model, _ := userRepository.FindByID(ctx, user2)

		if user1Model.Balance != 900 {
			t.Errorf("expected user1 balance to be 900, got %v", user1Model.Balance)
		}

		if user2Model.Balance != 600 {
			t.Errorf("expected user2 balance to be 600, got %v", user2Model.Balance)
		}

		helpers.PrettyPrint(user1Model, user2Model)
	})
}
