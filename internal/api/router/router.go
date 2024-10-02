package router

import (
	"net/http"

	"github.com/edulustosa/go-pay/internal/api/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewServer(pool *pgxpool.Pool) http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("POST /users", handlers.HandleCreateUser(pool))

	return r
}
