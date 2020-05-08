package server

import (
	"net/http"

	"github.com/expandorg/verification/pkg/api/assignmentcreator"
	"github.com/expandorg/verification/pkg/api/assignmentdestroyer"
	"github.com/expandorg/verification/pkg/api/assignmentfetcher"
	"github.com/expandorg/verification/pkg/api/assignmentmaker"
	"github.com/expandorg/verification/pkg/api/assignmentupdater"
	"github.com/expandorg/verification/pkg/api/jobassignmentsfetcher"
	"github.com/expandorg/verification/pkg/api/responsefetcher"
	"github.com/expandorg/verification/pkg/api/responseverifier"
	"github.com/expandorg/verification/pkg/authentication"

	"github.com/expandorg/verification/pkg/api/healthchecker"
	"github.com/expandorg/verification/pkg/api/settingcreator"
	"github.com/expandorg/verification/pkg/api/settingfetcher"
	"github.com/expandorg/verification/pkg/service"
	"github.com/gorilla/mux"
)

func New(s service.VerificationService) http.Handler {
	r := mux.NewRouter()

	r.Handle("/_health", healthchecker.MakeHandler(s)).Methods("GET")

	r.Handle("/assign", assignmentcreator.MakeHandler(s)).Methods("POST")
	r.Handle("/assignments", assignmentmaker.MakeHandler(s)).Methods("POST")
	r.Handle("/assignments/{assignment_id}", assignmentdestroyer.MakeHandler(s)).Methods("DELETE")
	r.Handle("/verify/manual", responseverifier.MakeManualHandler(s)).Methods("POST")
	r.Handle("/verify/automatic", responseverifier.MakeAutomaticHandler(s)).Methods("POST")

	r.Handle("/assignments", assignmentfetcher.MakeAssignmentsFetcherHandler(s)).Methods("GET")
	r.Handle("/assignments/{assignment_id}", assignmentfetcher.MakeAssignmentFetcherHandler(s)).Methods("GET")
	r.Handle("/assignments", assignmentupdater.MakeHandler(s)).Methods("PATCH")
	r.Handle("/response", responsefetcher.MakeResponsesFetcherHandler(s)).Methods("GET")
	r.Handle("/response/{response_id}", responsefetcher.MakeResponseFetcherHandler(s)).Methods("GET")
	r.Handle("/settings/{job_id}", settingfetcher.MakeHandler(s)).Methods("GET")
	r.Handle("/settings", settingcreator.MakeHandler(s)).Methods("PUT")

	r.Handle("/jobs/available", jobassignmentsfetcher.MakeAvailableHandler(s)).Methods("GET")
	r.Handle("/jobs/eligible", jobassignmentsfetcher.MakeEligibleHandler(s)).Methods("GET")

	r.Use(authentication.AuthMiddleware)
	return withHandlers(r)
}
