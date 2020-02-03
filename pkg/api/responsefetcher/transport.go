package responsefetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gemsorg/verification/pkg/service"
	"github.com/gemsorg/verification/pkg/verification"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeResponsesFetcherHandler(s service.VerificationService) http.Handler {
	return kithttp.NewServer(
		makeResponsesFetcherEndpoint(s),
		decodeResponsesFetcherRequest,
		encodeResponse,
	)
}

func MakeResponseFetcherHandler(s service.VerificationService) http.Handler {
	return kithttp.NewServer(
		makeResponseFetcherEndpoint(s),
		decodeResponseFetcherRequest,
		encodeResponse,
	)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeResponsesFetcherRequest(_ context.Context, r *http.Request) (interface{}, error) {
	p := verification.Params{}
	params := r.URL.Query()

	workerID, ok := params["worker_id"]
	if ok && len(workerID) > 0 {
		p.WorkerID = workerID[0]
	}
	jobID, ok := params["job_id"]
	if ok && len(jobID) > 0 {
		p.JobID = jobID[0]
	}
	taskID, ok := params["task_id"]
	if ok && len(taskID) > 0 {
		p.TaskID = taskID[0]
	}
	return p, nil
}

func decodeResponseFetcherRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	var ok bool
	responseID, ok := vars["response_id"]
	if !ok {
		return nil, fmt.Errorf("missing response_id parameter")
	}
	return ResponseRequest{ResponseID: responseID}, nil
}
