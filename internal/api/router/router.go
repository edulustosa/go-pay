package router

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewServer(pool *pgxpool.Pool) http.Handler {
	r := http.NewServeMux()

	return r
}
