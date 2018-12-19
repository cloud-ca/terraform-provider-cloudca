package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

//Status codes
const (
	OK              = 200
	MUTIPLE_CHOICES = 300
	BAD_REQUEST     = 400
	NOT_FOUND       = 404
)

//An API error
type CcaError struct {
	ErrorCode string                 `json:"errorCode"`
	Message   string                 `json:"message"`
	Context   map[string]interface{} `json:"context"`
}

//An Api Response
type CcaResponse struct {
	TaskId     string
	TaskStatus string
	StatusCode int
	Data       []byte
	Errors     []CcaError
	MetaData   map[string]interface{}
}

//Returns true if API response has errors
func (ccaResponse CcaResponse) IsError() bool {
	return !isInOKRange(ccaResponse.StatusCode)
}

//An Api Response with errors
type CcaErrorResponse CcaResponse

func (errorResponse CcaErrorResponse) Error() string {
	var errorStr string = "[ERROR] Received HTTP status code " + strconv.Itoa(errorResponse.StatusCode) + "\n"
	for _, e := range errorResponse.Errors {
		context, _ := json.Marshal(e.Context)
		errorStr += "[ERROR] Error Code: " + e.ErrorCode + ", Message: " + e.Message + ", Context: " + string(context) + "\n"
	}
	return errorStr
}

func NewCcaResponse(response *http.Response) (*CcaResponse, error) {
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	ccaResponse := CcaResponse{}
	ccaResponse.StatusCode = response.StatusCode
	responseMap := map[string]*json.RawMessage{}
	json.Unmarshal(respBody, &responseMap)

	if val, ok := responseMap["taskId"]; ok {
		json.Unmarshal(*val, &ccaResponse.TaskId)
	}

	if val, ok := responseMap["taskStatus"]; ok {
		json.Unmarshal(*val, &ccaResponse.TaskStatus)
	}

	if val, ok := responseMap["data"]; ok {
		ccaResponse.Data = []byte(*val)
	}

	if val, ok := responseMap["metadata"]; ok {
		metadata := map[string]interface{}{}
		json.Unmarshal(*val, &metadata)
		ccaResponse.MetaData = metadata
	}

	if val, ok := responseMap["errors"]; ok {
		errors := []CcaError{}
		json.Unmarshal(*val, &errors)
		ccaResponse.Errors = errors
	} else if !isInOKRange(response.StatusCode) {
		return nil, fmt.Errorf("Unexpected. Received status " + response.Status + " but no errors in response body")
	}
	return &ccaResponse, nil
}

func isInOKRange(statusCode int) bool {
	return statusCode >= OK && statusCode < MUTIPLE_CHOICES
}
