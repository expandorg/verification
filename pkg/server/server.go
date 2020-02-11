package server

import (
	"net/http"

	"github.com/gemsorg/verification/pkg/api/assignmentcreator"
	"github.com/gemsorg/verification/pkg/api/assignmentfetcher"
	"github.com/gemsorg/verification/pkg/api/responsefetcher"
	"github.com/gemsorg/verification/pkg/api/responseverifier"
	"github.com/gemsorg/verification/pkg/authentication"

	"github.com/gemsorg/verification/pkg/api/healthchecker"
	"github.com/gemsorg/verification/pkg/api/settingcreator"
	"github.com/gemsorg/verification/pkg/api/settingfetcher"
	"github.com/gemsorg/verification/pkg/service"
	"github.com/gorilla/mux"
)

func New(s service.VerificationService) http.Handler {
	r := mux.NewRouter()

	r.Handle("/_health", healthchecker.MakeHandler(s)).Methods("GET")

	r.Handle("/assign", assignmentcreator.MakeHandler(s)).Methods("POST")
	r.Handle("/verify/manual", responseverifier.MakeManualHandler(s)).Methods("POST")
	r.Handle("/verify/automatic", responseverifier.MakeAutomaticHandler(s)).Methods("POST")

<<<<<<< HEAD
	r.Handle("/assignments", assignmentfetcher.MakeAssignmentsFetcherHandler(s)).Methods("GET")
	r.Handle("/assignments/{assignment_id}", assignmentfetcher.MakeAssignmentFetcherHandler(s)).Methods("GET")

	r.Handle("/response", responsecreator.MakeHandler(s)).Methods("POST")
=======
>>>>>>> eda6324... consensus verification
	r.Handle("/response", responsefetcher.MakeResponsesFetcherHandler(s)).Methods("GET")
	r.Handle("/response/{response_id}", responsefetcher.MakeResponseFetcherHandler(s)).Methods("GET")
	r.Handle("/settings/{job_id}", settingfetcher.MakeHandler(s)).Methods("GET")
	r.Handle("/settings", settingcreator.MakeHandler(s)).Methods("PUT")

	r.Use(authentication.AuthMiddleware)
	return withHandlers(r)
}
