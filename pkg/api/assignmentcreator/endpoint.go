package assignmentcreator

import (
	"context"

	"github.com/expandorg/verification/pkg/apierror"
	"github.com/expandorg/verification/pkg/authentication"
	"github.com/expandorg/verification/pkg/service"
	"github.com/expandorg/verification/pkg/verification"
	"github.com/go-kit/kit/endpoint"
)

func makeAssignmentCreatorEndpoint(svc service.VerificationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data, _ := authentication.ParseAuthData(ctx)
		svc.SetAuthData(data)

		req := request.(verification.NewAssignment)
		s, err := svc.GetSettings(req.JobID)
		if err != nil {
			return nil, err
		}
		a, err := svc.Assign(req, s)
		if err != nil {
			return nil, errorResponse(err)
		}
		return a, nil
	}
}

func errorResponse(err error) *apierror.APIError {
	return apierror.New(500, err.Error(), err)
}
