package factories

import (
	"github.com/edulustosa/go-pay/internal/database/repo"
	"github.com/edulustosa/go-pay/internal/services/transfer"
	"github.com/edulustosa/go-pay/internal/services/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

func MakeTransferService(pool *pgxpool.Pool) *transfer.Service {
	transactionRepository := repo.NewTransactionsRepository(pool)
	usersRepository := repo.NewUserRepository(pool)
	userService := user.NewService(usersRepository)
	transferService := transfer.NewService(
		transactionRepository,
		userService,
	)

	return transferService
}
