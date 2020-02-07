package externalsvc

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gemsorg/verification/pkg/registrysvc"
	"github.com/gemsorg/verification/pkg/verification"
)

type External interface {
	Verify(reg registrysvc.Registration, resp verification.NewResponse) (*VerifyResponse, error)
}

type VerifyResponse struct {
	Value bool `json:"value"`
}

type external struct {
	AuthToken string
}

func New(authToken string) *external {
	return &external{
		AuthToken: authToken,
	}
}

func (e *external) Verify(reg registrysvc.Registration, resp verification.NewResponse) (*VerifyResponse, error) {
	url := reg.Services[registrysvc.ResponseVerifier].URL
	requestByte, _ := json.Marshal(resp)
	reqBody := bytes.NewReader(requestByte)
	res, err := serviceRequest("POST", url, reg.APIKeyID, reg.RequesterID, reqBody)
	if err != nil {
		return nil, err
	}
	vresp := VerifyResponse{}
	err = json.Unmarshal(res, &vresp)
	if err != nil {
		return nil, err
	}
	return &vresp, nil
}

func serviceRequest(action, url, key string, userID uint64, reqBody io.Reader) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest(action, url, reqBody)
	if err != nil {
		return nil, err
	}
	apiKey, err := GenerateAPIKeyJWT(userID, key)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", apiKey)
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
