package settingcreator

import (
	"context"

	"github.com/expandorg/verification/pkg/apierror"
	"github.com/expandorg/verification/pkg/authentication"
	"github.com/expandorg/verification/pkg/service"
	"github.com/expandorg/verification/pkg/verification"
	"github.com/go-kit/kit/endpoint"
)

func makeSettingsCreatorEndpoint(svc service.VerificationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data, _ := authentication.ParseAuthData(ctx)
		svc.SetAuthData(data)
		req := request.(verification.Settings)
		settings, err := svc.CreateSettings(req)
		if err != nil {
			return nil, errorResponse(err)
		}
		return settings, nil
	}
}

func errorResponse(err error) *apierror.APIError {
	return apierror.New(500, err.Error(), err)
}
