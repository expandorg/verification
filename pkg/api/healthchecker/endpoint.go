package healthchecker

import (
	"context"

	service "github.com/expandorg/verification/pkg/service"

	"github.com/go-kit/kit/endpoint"
)

func makeHealthyEndpoint(svc service.VerificationService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		healthy := svc.Healthy()
		return HealthyResponse{healthy}, nil
	}
}

type HealthyResponse struct {
	Healthy bool `json:"healthy"`
}
