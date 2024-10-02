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

const createUser = `
	INSERT INTO users (
		"first_name",
		"last_name",
		"document",
		"email",
		"password_hash",
		"role"
	) VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING "id";
`

func (r *UserRepository) Create(
	ctx context.Context,
	user models.User,
) (uuid.UUID, error) {
	var id uuid.UUID

	err := r.db.QueryRow(
		ctx,
		createUser,
		user.FirstName,
		user.LastName,
		user.Document,
		user.Email,
		user.PasswordHash,
		user.Role,
	).Scan(&id)

	return id, err
}

const findByEmail = "SELECT * FROM users WHERE email = $1"

func (r *UserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (models.User, error) {
	row := r.db.QueryRow(ctx, findByEmail, email)
	return scanUser(row)
}

const updateBalance = "UPDATE users SET balance = $2 WHERE id = $1"

func (r *UserRepository) UpdateBalance(
	ctx context.Context,
	id uuid.UUID,
	amount float64,
) error {
	_, err := r.db.Exec(ctx, updateBalance, id, amount)
	return err
}
