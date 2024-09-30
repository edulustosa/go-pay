package repo

import (
	"context"
	"time"

	"github.com/edulustosa/go-pay/internal/database/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type InMemoryTransactionsRepository struct {
	Transaction []models.Transaction
}

func (r *InMemoryTransactionsRepository) Create(
	_ context.Context,
	transaction models.Transaction,
) (uuid.UUID, error) {
	transaction.ID = uuid.New()
	transaction.CreatedAt = pgtype.Timestamp{Time: time.Now()}
	transaction.UpdatedAt = pgtype.Timestamp{Time: time.Now()}

	r.Transaction = append(r.Transaction, transaction)
	return transaction.ID, nil
}
