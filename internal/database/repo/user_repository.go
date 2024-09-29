package repo

import (
	"context"

	"github.com/edulustosa/go-pay/internal/database/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db,
	}
}

func scanUser(row pgx.Row) (models.User, error) {
	var user models.User
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Document,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}

const findByDocument = "SELECT * FROM users WHERE document = $1"

func (r *UserRepository) FindByDocument(
	ctx context.Context,
	document string,
) (models.User, error) {
	row := r.db.QueryRow(ctx, findByDocument, document)
	return scanUser(row)
}

const findByID = "SELECT * FROM users WHERE id = $1"

func (r *UserRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (models.User, error) {
	row := r.db.QueryRow(ctx, findByID, id)
	return scanUser(row)
}
