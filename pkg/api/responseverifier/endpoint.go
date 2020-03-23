package responseverifier

import (
	"context"

	"github.com/expandorg/verification/pkg/apierror"
	"github.com/expandorg/verification/pkg/authentication"
	"github.com/expandorg/verification/pkg/service"
	"github.com/expandorg/verification/pkg/verification"
	"github.com/go-kit/kit/endpoint"
)

func makeAutomaticEndpoint(svc service.VerificationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data, _ := authentication.ParseAuthData(ctx)
		svc.SetAuthData(data)

		req := request.(verification.TaskResponse)
		settings, err := svc.GetSettings(req.JobID)
		if err != nil {
			return nil, err
		}
		r, err := svc.VerifyAutomatic(req, settings)
		if err != nil {
			return nil, errorResponse(err)
		}
		return r, nil
	}
}

func makeManualEndpoint(svc service.VerificationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data, _ := authentication.ParseAuthData(ctx)
		svc.SetAuthData(data)

		req := request.(verification.NewVerificationResponse)
		settings, err := svc.GetSettings(req.JobID)
		if err != nil {
			return nil, err
		}
		r, err := svc.VerifyManual(req, settings)
		if err != nil {
			return nil, errorResponse(err)
		}
		return r, nil
	}
}

func errorResponse(err error) *apierror.APIError {
	return apierror.New(500, err.Error(), err)
}
