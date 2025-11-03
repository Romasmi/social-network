package routes

import (
	"github.com/Romasmi/social-network/internal/config"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserRoutes(router *mux.Router, db *pgxpool.Pool, cng *config.Config) {

}
