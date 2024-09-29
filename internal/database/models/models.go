package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Role string

const (
	RoleCommon   Role = "COMMON"
	RoleMerchant Role = "MERCHANT"
)

type User struct {
	ID           uuid.UUID
	FirstName    string
	LastName     string
	Document     string
	Email        string
	PasswordHash string
	Balance      float64
	Role         Role
	CreatedAt    pgtype.Timestamp
	UpdatedAt    pgtype.Timestamp
}

type Transaction struct {
	ID        uuid.UUID
	Amount    float64
	Payer     uuid.UUID
	Payee     uuid.UUID
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}
