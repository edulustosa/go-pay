package repo

import (
	"context"

	"github.com/edulustosa/go-pay/internal/database/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionsRepository struct {
	db *pgxpool.Pool
}

func NewTransactionsRepository(db *pgxpool.Pool) *TransactionsRepository {
	return &TransactionsRepository{
		db,
	}
}

const create = `
	INSERT INTO transactions (
		payer,
		payee,
		amount
	) VALUES ($1, $2, $3)
	RETURNING id;
`

func (r *TransactionsRepository) Create(
	ctx context.Context,
	transaction models.Transaction,
) (uuid.UUID, error) {
	var id uuid.UUID

	err := r.db.QueryRow(
		ctx,
		create,
		transaction.Payer,
		transaction.Payee,
		transaction.Amount,
	).Scan(&id)

	return id, err
}
