package server

import (
	"net/http"

	"github.com/gemsorg/verification/pkg/authentication"

	"github.com/jmoiron/sqlx"

	"github.com/gemsorg/verification/pkg/api/healthchecker"
	"github.com/gemsorg/verification/pkg/service"
	"github.com/gorilla/mux"
)

func New(
	db *sqlx.DB,
	s service.VerificationService,
) http.Handler {
	r := mux.NewRouter()

	r.Handle("/_health", healthchecker.MakeHandler(s)).Methods("GET")
	r.Use(authentication.AuthMiddleware)
	return withHandlers(r)
}
