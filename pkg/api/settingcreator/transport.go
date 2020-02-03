package settingcreator

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
		makeSettingsCreatorEndpoint(s),
		decodeNewSettingsRequest,
		encodeResponse,
	)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeNewSettingsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var s verification.Settings
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&s)
	if err != nil {
		return nil, apierror.New(500, err.Error(), err)
	}
	return s, nil
}
