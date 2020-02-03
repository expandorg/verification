package responsecreator

import (
	"context"

	"github.com/gemsorg/verification/pkg/apierror"
	"github.com/gemsorg/verification/pkg/authentication"
	"github.com/gemsorg/verification/pkg/service"
	"github.com/gemsorg/verification/pkg/verification"
	"github.com/go-kit/kit/endpoint"
)

func makeResponseCreatorEndpoint(svc service.VerificationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data, _ := authentication.ParseAuthData(ctx)
		svc.SetAuthData(data)
		req := request.(verification.NewResponse)
		r, err := svc.CreateResponse(req)
		if err != nil {
			return nil, errorResponse(err)
		}
		return r, nil
	}
}

func errorResponse(err error) *apierror.APIError {
	return apierror.New(500, err.Error(), err)
}
