package responsesvc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/expandorg/verification/pkg/apierror"
	"github.com/expandorg/verification/pkg/verification"
)

type responsesvc struct {
}

func New(token string) ResponseSVC {
	return &responsesvc{}
}

func (rs *responsesvc) GetPending(authToken string, taskID uint64) (verification.TaskResponses, error) {
	r := PendingResult{}
	route := fmt.Sprintf("tasks/%d/responses/pending", taskID)
	result, err := rs.serviceRequest("GET", route, authToken, nil)
	if err != nil {
		return r.Responses, err
	}

	err = json.Unmarshal(result, &r)
	if err != nil {
		return r.Responses, err
	}
	return r.Responses, nil
}

func (rs *responsesvc) serviceRequest(action string, route string, authToken string, reqBody io.Reader) ([]byte, error) {
	client := &http.Client{}
	serviceURL := fmt.Sprintf("%s/%s", os.Getenv("RESPONSES_SVC_URL"), route)

	req, err := http.NewRequest(action, serviceURL, reqBody)
	if err != nil {
		return nil, errorResponse(err)
	}

	req.AddCookie(&http.Cookie{Name: "JWT", Value: authToken})

	r, err := client.Do(req)
	if err != nil {
		return nil, errorResponse(err)
	}

	if r.StatusCode != 200 {
		return nil, errors.New("bad request")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errorResponse(err)
	}
	return body, nil
}

func errorResponse(err error) *apierror.APIError {
	return apierror.New(500, err.Error(), err)
}
