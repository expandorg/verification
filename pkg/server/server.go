package server

import (
	"net/http"

	"github.com/gemsorg/verification/pkg/api/responsefetcher"
	"github.com/gemsorg/verification/pkg/api/responseverifier"
	"github.com/gemsorg/verification/pkg/authentication"

	"github.com/gemsorg/verification/pkg/api/healthchecker"
	"github.com/gemsorg/verification/pkg/api/responsecreator"
	"github.com/gemsorg/verification/pkg/api/settingcreator"
	"github.com/gemsorg/verification/pkg/api/settingfetcher"
	"github.com/gemsorg/verification/pkg/service"
	"github.com/gorilla/mux"
)

func New(s service.VerificationService) http.Handler {
	r := mux.NewRouter()

	r.Handle("/_health", healthchecker.MakeHandler(s)).Methods("GET")

	r.Handle("/verify/manual", responseverifier.MakeManualHandler(s)).Methods("POST")
	r.Handle("/verify/automatic", responseverifier.MakeAutomaticHandler(s)).Methods("POST")

	r.Handle("/response", responsecreator.MakeHandler(s)).Methods("POST")
	r.Handle("/response", responsefetcher.MakeResponsesFetcherHandler(s)).Methods("GET")
	r.Handle("/response/{response_id}", responsefetcher.MakeResponseFetcherHandler(s)).Methods("GET")
	r.Handle("/settings/{job_id}", settingfetcher.MakeHandler(s)).Methods("GET")
	r.Handle("/settings", settingcreator.MakeHandler(s)).Methods("PUT")

	r.Use(authentication.AuthMiddleware)
	return withHandlers(r)
}
