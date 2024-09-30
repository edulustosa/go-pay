package dtos

import "github.com/google/uuid"

type TransactionDTO struct {
	Value float64   `json:"amount"`
	Payer uuid.UUID `json:"payer"`
	Payee uuid.UUID `json:"payee"`
}
