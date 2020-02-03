package settingfetcher

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/gemsorg/verification/pkg/apierror"
	"github.com/gemsorg/verification/pkg/authentication"
	"github.com/gemsorg/verification/pkg/service"
	"github.com/go-kit/kit/endpoint"
)

func makeSettingFetcherEndpoint(svc service.VerificationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data, _ := authentication.ParseAuthData(ctx)
		svc.SetAuthData(data)
		req := request.(SettingRequest)

		jobID, err := strconv.ParseUint(req.JobID, 10, 64)
		if err != nil {
			return nil, errorResponse(err)
		}

		s, err := svc.GetSettings(jobID)
		if err != nil {
			return nil, errorResponse(err)
		}
		if s == nil {
			return json.RawMessage("{}"), nil
		}
		return s, nil
	}
}

func errorResponse(err error) *apierror.APIError {
	return apierror.New(500, err.Error(), err)
}

type SettingRequest struct {
	JobID string
}
