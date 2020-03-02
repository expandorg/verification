package assignmentfetcher

import (
	"context"

	"github.com/expandorg/verification/pkg/apierror"
	"github.com/expandorg/verification/pkg/authentication"
	"github.com/expandorg/verification/pkg/service"
	"github.com/expandorg/verification/pkg/verification"
	"github.com/go-kit/kit/endpoint"
)

func makeAssignmentsFetcherEndpoint(svc service.VerificationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data, _ := authentication.ParseAuthData(ctx)
		svc.SetAuthData(data)
		params := request.(verification.Params)
		assignments, err := svc.GetAssignments(params)
		if err != nil {
			return nil, errorResponse(err)
		}
		return assignments, nil
	}
}

func makeAssignmentFetcherEndpoint(svc service.VerificationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data, _ := authentication.ParseAuthData(ctx)
		svc.SetAuthData(data)
		req := request.(AssignmentRequest)
		p, err := svc.GetAssignment(req.AssignmentID)
		if err != nil {
			return p, errorResponse(err)
		}
		return p, nil
	}
}

func errorResponse(err error) *apierror.APIError {
	return apierror.New(500, err.Error(), err)
}

type AssignmentRequest struct {
	AssignmentID string
}
