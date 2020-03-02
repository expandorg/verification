package externalsvc

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/expandorg/verification/pkg/registrysvc"
	"github.com/expandorg/verification/pkg/verification"
)

type External interface {
	Verify(reg *registrysvc.Registration, resp verification.TaskResponse) (VerificationResults, error)
}

type VerificationResult struct {
	JobID      uint64 `json:"job_id"`
	TaskID     uint64 `json:"task_id"`
	WorkerID   uint64 `json:"worker_id"`
	ResponseID uint64 `json:"response_id"`
	Accepted   bool   `json:"accepted"`
}

type VerificationResults []VerificationResult

type external struct {
	AuthToken string
}

func New(authToken string) *external {
	return &external{
		AuthToken: authToken,
	}
}

func (e *external) Verify(reg *registrysvc.Registration, resp verification.TaskResponse) (VerificationResults, error) {
	url := reg.Services[registrysvc.ResponseVerifier].URL
	requestByte, _ := json.Marshal(resp)
	reqBody := bytes.NewReader(requestByte)
	res, err := serviceRequest("POST", url, reg.APIKeyID, reg.RequesterID, reqBody)
	if err != nil {
		return nil, err
	}
	vresp := VerificationResults{}
	err = json.Unmarshal(res, &vresp)
	if err != nil {
		return nil, err
	}
	return vresp, nil
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

func (vr *VerificationResult) ToVerificationResponse() verification.VerificationResponse {
	return verification.VerificationResponse{
		JobID:      vr.JobID,
		TaskID:     vr.TaskID,
		ResponseID: vr.ResponseID,
		WorkerID:   vr.WorkerID,
		Accepted:   vr.Accepted,
	}
}

func (vrs VerificationResults) ToVerificationResponses() verification.VerificationResponses {
	results := make(verification.VerificationResponses, len(vrs))
	for i, vr := range vrs {
		results[i] = vr.ToVerificationResponse()
	}
	return results
}
