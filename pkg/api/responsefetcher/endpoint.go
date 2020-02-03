package responsefetcher

import (
	"context"

	"github.com/gemsorg/verification/pkg/apierror"
	"github.com/gemsorg/verification/pkg/authentication"
	"github.com/gemsorg/verification/pkg/service"
	"github.com/gemsorg/verification/pkg/verification"
	"github.com/go-kit/kit/endpoint"
)

func makeResponsesFetcherEndpoint(svc service.VerificationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data, _ := authentication.ParseAuthData(ctx)
		svc.SetAuthData(data)
		params := request.(verification.Params)
		responses, err := svc.GetResponses(params)
		if err != nil {
			return nil, errorResponse(err)
		}
		return responses, nil
	}
}

func makeResponseFetcherEndpoint(svc service.VerificationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data, _ := authentication.ParseAuthData(ctx)
		svc.SetAuthData(data)
		req := request.(ResponseRequest)
		r, err := svc.GetResponse(req.ResponseID)
		if err != nil {
			return r, errorResponse(err)
		}
		return r, nil
	}
}

func errorResponse(err error) *apierror.APIError {
	return apierror.New(500, err.Error(), err)
}

type ResponseRequest struct {
	ResponseID string
}
