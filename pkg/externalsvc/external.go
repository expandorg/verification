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
	Verify(reg *registrysvc.Registration, resp verification.TaskResponse) (*VerifyResponse, error)
}

type VerifyResponse struct {
	Results []VerificationResult `json:"results"`
}

type VerificationResult struct {
	JobID      uint64 `json:"job_id"`
	TaskID     uint64 `json:"task_id"`
	ResponseID uint64 `json:"response_id"`
	Accepted   bool   `json:"accepted`
}

type external struct {
	AuthToken string
}

func New(authToken string) *external {
	return &external{
		AuthToken: authToken,
	}
}

func (e *external) Verify(reg *registrysvc.Registration, resp verification.TaskResponse) (*VerifyResponse, error) {
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
