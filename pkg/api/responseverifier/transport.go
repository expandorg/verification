package responseverifier

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gemsorg/verification/pkg/apierror"
	"github.com/gemsorg/verification/pkg/service"
	"github.com/gemsorg/verification/pkg/verification"
	kithttp "github.com/go-kit/kit/transport/http"
)

func MakeAutomaticHandler(s service.VerificationService) http.Handler {
	return kithttp.NewServer(
		makeAutomaticEndpoint(s),
		decodeAutomaticRequest,
		encodeResponse,
	)
}

func MakeManualHandler(s service.VerificationService) http.Handler {
	return kithttp.NewServer(
		makeManualEndpoint(s),
		decodeManualRequest,
		encodeResponse,
	)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeAutomaticRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var tr verification.TaskResponse
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&tr)
	if err != nil {
		return nil, apierror.New(500, err.Error(), err)
	}
	return tr, nil
}

func decodeManualRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var v verification.NewVerificationResponse
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&v)
	if err != nil {
		return nil, apierror.New(500, err.Error(), err)
	}
	return v, nil
}
