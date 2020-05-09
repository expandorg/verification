package jobassignmentsfetcher

import (
	"context"

	"github.com/expandorg/verification/pkg/apierror"
	"github.com/expandorg/verification/pkg/authentication"
	"github.com/expandorg/verification/pkg/service"
	"github.com/go-kit/kit/endpoint"
)

func makeEligibleJobsFetcherEndpoint(svc service.VerificationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data, _ := authentication.ParseAuthData(ctx)
		svc.SetAuthData(data)

		params := request.(VerifierRequest)

		j, err := svc.GetEligibleJobs(params.VerifierID)
		if err != nil {
			return nil, errorResponse(err)
		}
		return j, nil
	}
}

func errorResponse(err error) *apierror.APIError {
	return apierror.New(500, err.Error(), err)
}

type VerifierRequest struct {
	VerifierID uint64
}
