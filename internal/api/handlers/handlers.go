package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/edulustosa/go-pay/internal/database/repo"
	"github.com/edulustosa/go-pay/internal/dtos"
	"github.com/edulustosa/go-pay/internal/factories"
	"github.com/edulustosa/go-pay/internal/services/transfer"
	"github.com/edulustosa/go-pay/internal/services/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JSON map[string]any

type Validator interface {
	Valid() (problems map[string]string)
}

func decodeValid[T Validator](r *http.Request) (T, map[string]string, error) {
	var v T

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, nil, fmt.Errorf("decode json: %w", err)
	}

	if problems := v.Valid(); len(problems) > 0 {
		return v, problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
	}

	return v, nil, nil
}

func encode[T any](w http.ResponseWriter, status int, v T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type ErrorList struct {
	Errors []Error `json:"errors"`
}

type Error struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

var InternalServerErr = Error{Message: "something went wrong, please try again later"}

func handleError(w http.ResponseWriter, status int, errors ...Error) {
	encode(w, status, ErrorList{errors})
}

func handleInvalidRequest(w http.ResponseWriter, problems map[string]string) {
	var errors []Error

	if len(problems) > 0 {
		errors = make([]Error, 0, len(problems))
		for field, problem := range problems {
			err := Error{
				Message: fmt.Sprintf("invalid %s", field),
				Details: problem,
			}
			errors = append(errors, err)
		}
	} else {
		errors = make([]Error, 0, 1)
		errors = append(errors, Error{
			Message: "invalid input",
			Details: "failed to parse request",
		})
	}

	handleError(w, http.StatusBadRequest, errors...)
}

func HandleCreateUser(pool *pgxpool.Pool) http.HandlerFunc {
	usersRepository := repo.NewUserRepository(pool)
	userService := user.NewService(usersRepository)

	return func(w http.ResponseWriter, r *http.Request) {
		req, problems, err := decodeValid[dtos.UserDTO](r)
		if err != nil {
			handleInvalidRequest(w, problems)
			return
		}

		userID, err := userService.Create(r.Context(), req)
		if err != nil {
			if errors.Is(err, user.ErrUserAlreadyExists) {
				handleError(w, http.StatusConflict, Error{
					Message: err.Error(),
					Details: "an user with the same email or document already exists",
				})
				return
			}

			slog.Error("failed to create user", "error", err, "user", req)
			handleError(w, http.StatusInternalServerError, InternalServerErr)
			return
		}

		encode(w, http.StatusCreated, JSON{"id": userID})
	}
}

func HandleTransfer(pool *pgxpool.Pool) http.HandlerFunc {
	transferService := factories.MakeTransferService(pool)

	return func(w http.ResponseWriter, r *http.Request) {
		req, problems, err := decodeValid[dtos.TransactionDTO](r)
		if err != nil {
			handleInvalidRequest(w, problems)
			return
		}

		_, err = transferService.NewTransaction(r.Context(), req)
		if err != nil {
			if errors.Is(err, transfer.ErrUserNotFound) {
				handleError(w, http.StatusNotFound, Error{
					Message: err.Error(),
				})
				return
			}

			if errors.Is(err, transfer.ErrMerchantNotAllowed) {
				handleError(w, http.StatusForbidden, Error{
					Message: err.Error(),
				})
				return
			}

			if errors.Is(err, transfer.ErrInsufficientFunds) {
				handleError(w, http.StatusUnprocessableEntity, Error{
					Message: err.Error(),
					Details: "payer has insufficient funds",
				})
				return
			}

			if errors.Is(err, transfer.ErrTransactionNotAuthorized) {
				handleError(w, http.StatusUnauthorized, Error{
					Message: err.Error(),
				})
				return
			}

			slog.Error("failed to make transfer", "error", err, "transfer", req)
			handleError(w, http.StatusInternalServerError, InternalServerErr)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
