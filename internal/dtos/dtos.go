package dtos

import (
	"fmt"
	"net/mail"

	"github.com/edulustosa/go-pay/helpers"
	"github.com/edulustosa/go-pay/internal/database/models"
	"github.com/google/uuid"
)

func validLength(field string, min, max int) bool {
	return len(field) >= min && len(field) <= max
}

type TransactionDTO struct {
	Value float64   `json:"amount"`
	Payer uuid.UUID `json:"payer"`
	Payee uuid.UUID `json:"payee"`
}

type NotificationDTO struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

type UserDTO struct {
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	Document  string      `json:"document"`
	Email     string      `json:"email"`
	Password  string      `json:"password"`
	Balance   float64     `json:"balance"`
	Role      models.Role `json:"role,omitempty"`
}

func (u UserDTO) Valid() (problems map[string]string) {
	problems = make(map[string]string)

	if !validLength(u.FirstName, 3, 255) {
		problems["firstName"] = "must be between 3 and 255 characters"
	}

	if !validLength(u.LastName, 3, 255) {
		problems["lastName"] = "must be between 3 and 255 characters"
	}

	if err := helpers.ParseDocument(u.Document); err != nil {
		problems["document"] = err.Error()
	}

	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		problems["email"] = fmt.Sprintf("%s is not a valid email", u.Email)
	}

	if !validLength(u.Password, 6, 255) {
		problems["password"] = "must be between 6 and 255 characters"
	}

	if u.Balance < 0 {
		problems["balance"] = "must be greater than or equal to 0"
	}

	return problems
}
