package authorization

import (
	"github.com/gemsorg/verification/pkg/authentication"
)

type Authorizer interface {
	SetAuthData(data authentication.AuthData)
	GetUserID() uint64
}

type authorizor struct {
	authData authentication.AuthData
}

func NewAuthorizer() Authorizer {
	return &authorizor{
		authentication.AuthData{},
	}
}

func (a *authorizor) SetAuthData(data authentication.AuthData) {
	a.authData = data
}

func (a *authorizor) GetUserID() uint64 {
	return a.authData.UserID
}
