package assignmentdestroyer

import (
	"context"

	"github.com/expandorg/verification/pkg/apierror"
	"github.com/expandorg/verification/pkg/authentication"
	"github.com/expandorg/verification/pkg/service"
	"github.com/go-kit/kit/endpoint"
)

func makeAssignmentDestroyerEndpoint(svc service.VerificationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data, _ := authentication.ParseAuthData(ctx)
		svc.SetAuthData(data)
		req := request.(AssignmentRequest)
		p, err := svc.DeleteAssignment(req.AssignmentID)
		if err != nil {
			return AssignmentResponse{p}, errorResponse(err)
		}
		return AssignmentResponse{p}, nil
	}
}

func errorResponse(err error) *apierror.APIError {
	return apierror.New(500, err.Error(), err)
}

type AssignmentRequest struct {
	AssignmentID string `json:"assignment_id"`
}

type AssignmentResponse struct {
	Destroyed bool `json:"destroyed"`
}
