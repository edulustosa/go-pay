package transfer

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Rhymond/go-money"
	"github.com/edulustosa/go-pay/internal/database/models"
	"github.com/edulustosa/go-pay/internal/dtos"
	"github.com/edulustosa/go-pay/internal/services/notification"
	"github.com/google/uuid"
)

type transactionsRepository interface {
	Create(
		ctx context.Context,
		transaction models.Transaction,
	) (uuid.UUID, error)
}

type userService interface {
	FindByID(ctx context.Context, id uuid.UUID) (models.User, error)
	UpdateBalance(ctx context.Context, id uuid.UUID, balance float64) error
}

type Service struct {
	repo transactionsRepository
	user userService
}

func NewService(
	repo transactionsRepository,
	user userService,
) *Service {
	return &Service{
		repo,
		user,
	}
}

var (
	ErrMerchantNotAllowed       = errors.New("merchant not allowed to make transactions")
	ErrInsufficientFunds        = errors.New("insufficient funds")
	ErrTransactionNotAuthorized = errors.New("transaction not authorized")
	ErrUserNotFound             = errors.New("user not found")
)

func validateTransaction(payer *models.User, amount *money.Money) error {
	if payer.Role == models.RoleMerchant {
		return ErrMerchantNotAllowed
	}

	balance := money.NewFromFloat(payer.Balance, money.BRL)
	ok, err := balance.LessThan(amount)
	if err != nil {
		return err
	}
	if ok {
		return ErrInsufficientFunds
	}

	return nil
}

func (s *Service) NewTransaction(
	ctx context.Context,
	transactionDTO dtos.TransactionDTO,
) (uuid.UUID, error) {
	payer, err := s.user.FindByID(ctx, transactionDTO.Payer)
	if err != nil {
		return uuid.Nil, ErrUserNotFound
	}

	payee, err := s.user.FindByID(ctx, transactionDTO.Payee)
	if err != nil {
		return uuid.Nil, ErrUserNotFound
	}

	amount := money.NewFromFloat(transactionDTO.Value, money.BRL)
	if err = validateTransaction(&payer, amount); err != nil {
		return uuid.Nil, err
	}

	if err = authorizeTransaction(); err != nil {
		return uuid.Nil, ErrTransactionNotAuthorized
	}

	// The amount is negative because it is being subtracted from the payer
	err = s.updateBalance(ctx, &payer, amount.Negative())
	if err != nil {
		return uuid.Nil, err
	}

	err = s.updateBalance(ctx, &payee, amount)
	if err != nil {
		return uuid.Nil, err
	}

	transaction := models.Transaction{
		Amount: amount.AsMajorUnits(),
		Payer:  payer.ID,
		Payee:  payee.ID,
	}

	id, err := s.repo.Create(ctx, transaction)
	if err != nil {
		return uuid.Nil, err
	}

	go func() {
		err := notification.Send(&payer, "Transaction completed successfully")
		if err != nil {
			slog.Error("failed to send notification", "error", err)
		}
		err = notification.Send(&payee, "Transaction received successfully")
		if err != nil {
			slog.Error("failed to send notification", "error", err)
		}
	}()

	return id, nil
}

// Updates the user balance by the given amount
func (s *Service) updateBalance(
	ctx context.Context,
	user *models.User,
	amount *money.Money,
) error {
	balance := money.NewFromFloat(user.Balance, money.BRL)

	balance, err := balance.Add(amount)
	if err != nil {
		return err
	}

	return s.user.UpdateBalance(ctx, user.ID, balance.AsMajorUnits())
}

type Authorizer struct {
	Status string `json:"status"`
	Data   struct {
		Authorization bool `json:"authorization"`
	}
}

func authorizeTransaction() error {
	resp, err := http.Get("https://util.devi.tools/api/v2/authorize")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrTransactionNotAuthorized
	}

	var authorizer Authorizer
	if err := json.NewDecoder(resp.Body).Decode(&authorizer); err != nil {
		return err
	}

	if authorizer.Status != "success" || !authorizer.Data.Authorization {
		return ErrTransactionNotAuthorized
	}

	return nil
}
