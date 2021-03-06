package service

import (
	"testing"

	"github.com/expandorg/verification/pkg/authorization"
	"github.com/expandorg/verification/pkg/automatic"
	"github.com/expandorg/verification/pkg/datastore"
	"github.com/expandorg/verification/pkg/externalsvc"
	"github.com/expandorg/verification/pkg/registrysvc"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	authorizer := authorization.NewAuthorizer()
	ctrl := gomock.NewController(t)
	ds := &datastore.VerificationStore{}
	registry := registrysvc.NewMockRegistrySVC(ctrl)
	external := externalsvc.NewMockExternal(ctrl)
	consensus := automatic.NewMockConsensus(ctrl)

	type args struct {
		s *datastore.VerificationStore
	}
	tests := []struct {
		name string
		args args
		want *service
	}{
		{
			"it creates a new service",
			args{s: ds},
			&service{ds, authorizer, registry, external, consensus},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.s, authorizer, registry, external, consensus)
			assert.Equal(t, got, tt.want, tt.name)
		})
	}
}

func TestHealthy(t *testing.T) {
	ds := &datastore.VerificationStore{}
	type fields struct {
		store *datastore.VerificationStore
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"it returns true if healthy",
			fields{store: ds},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				store: tt.fields.store,
			}
			got := s.Healthy()
			assert.Equal(t, got, tt.want, tt.name)
		})
	}
}
