package assignmentcreator

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gemsorg/verification/pkg/apierror"
	"github.com/gemsorg/verification/pkg/service"
	"github.com/gemsorg/verification/pkg/verification"
	kithttp "github.com/go-kit/kit/transport/http"
)

func MakeHandler(s service.VerificationService) http.Handler {
	return kithttp.NewServer(
		makeAssignmentCreatorEndpoint(s),
		decodeRequest,
		encodeResponse,
	)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var a verification.NewAssignment
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&a)
	if err != nil {
		return nil, apierror.New(500, err.Error(), err)
	}
	return a, nil
}
