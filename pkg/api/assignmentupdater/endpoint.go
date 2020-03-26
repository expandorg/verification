package assignmentupdater

import (
	"context"

	"github.com/expandorg/verification/pkg/apierror"
	"github.com/expandorg/verification/pkg/authentication"
	"github.com/expandorg/verification/pkg/service"
	"github.com/expandorg/verification/pkg/verification"
	"github.com/go-kit/kit/endpoint"
)

func makeAssignmentUpdaterEndpoint(svc service.VerificationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data, _ := authentication.ParseAuthData(ctx)
		svc.SetAuthData(data)
		p, err := svc.UpdateAssignment(request.(verification.Assignment))
		if err != nil {
			return AssignmentResponse{p != nil}, errorResponse(err)
		}
		return AssignmentResponse{p != nil}, nil
	}
}

func errorResponse(err error) *apierror.APIError {
	return apierror.New(500, err.Error(), err)
}

type AssignmentResponse struct {
	Updated bool `json:"updated"`
}
