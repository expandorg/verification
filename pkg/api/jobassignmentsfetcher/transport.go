package jobassignmentsfetcher

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/expandorg/verification/pkg/service"
	kithttp "github.com/go-kit/kit/transport/http"
)

func MakeHandler(s service.VerificationService) http.Handler {
	return kithttp.NewServer(
		makeJobEmptyAssignmentsFetcherEndpoint(s),
		decodeRequest,
		encodeResponse,
	)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}
