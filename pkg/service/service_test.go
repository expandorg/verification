package service

import (
	"reflect"
	"testing"

	"github.com/gemsorg/verification/pkg/authentication"

	"github.com/gemsorg/verification/pkg/authorization"
	"github.com/gemsorg/verification/pkg/datastore"
	"github.com/gemsorg/verification/pkg/filter"
	"github.com/gemsorg/verification/pkg/mock"
	"github.com/gemsorg/verification/pkg/workerprofile"
	"github.com/stretchr/testify/assert"
)


func TestNew(t *testing.T) {
	authorizer := authorization.NewAuthorizer()
	ds := &datastore.VerificationStore{}
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
			&service{ds, authorizer},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.s, authorizer)
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