package jobassignmentsfetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/expandorg/verification/pkg/service"
	kithttp "github.com/go-kit/kit/transport/http"
)

func MakeAvailableHandler(s service.VerificationService) http.Handler {
	return kithttp.NewServer(
		makeJobEmptyAssignmentsFetcherEndpoint(s),
		decodeEmptyRequest,
		encodeResponse,
	)
}

func MakeEligibleHandler(s service.VerificationService) http.Handler {
	return kithttp.NewServer(
		makeEligibleJobsFetcherEndpoint(s),
		decodeEligibleRequest,
		encodeResponse,
	)
}

func decodeEmptyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeEligibleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	params := r.URL.Query()
	param, ok := params["verifier_id"]

	if !ok && len(param) == 0 {
		return nil, fmt.Errorf("missing verifier_id parameter")
	}

	verifierID, err := strconv.ParseUint(param[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("missing verifier_id parameter")
	}

	return VerifierRequest{VerifierID: verifierID}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
