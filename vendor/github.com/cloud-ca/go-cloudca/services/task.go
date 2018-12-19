package services

import (
	"encoding/json"
	"github.com/cloud-ca/go-cloudca/api"
	"strings"
	"time"
)

//Task status
const (
	PENDING = "PENDING"
	SUCCESS = "SUCCESS"
	FAILED  = "FAILED"
)

const DEFAULT_POLLING_INTERVAL = 1000

//A Task object. This object can be used to poll asynchronous operations.
type Task struct {
	Id      string
	Status  string
	Created string
	Result  []byte
}

type FailedTask Task

func (ft FailedTask) Error() string {
	return "Task id=" + ft.Id + " failed" //should add reason
}

type TaskService interface {
	Get(id string) (*Task, error)
	Poll(id string, milliseconds time.Duration) ([]byte, error)
	PollResponse(response *api.CcaResponse, milliseconds time.Duration) ([]byte, error)
}

type TaskApi struct {
	apiClient api.ApiClient
}

//Create a new TaskService
func NewTaskService(apiClient api.ApiClient) TaskService {
	return &TaskApi{
		apiClient: apiClient,
	}
}

//Retrieve a Task with sepecified id
func (taskApi *TaskApi) Get(id string) (*Task, error) {
	request := api.CcaRequest{
		Method:   api.GET,
		Endpoint: "tasks/" + id,
	}
	response, err := taskApi.apiClient.Do(request)
	if err != nil {
		return nil, err
	} else if len(response.Errors) > 0 {
		return nil, api.CcaErrorResponse(*response)
	}
	data := response.Data
	taskMap := map[string]*json.RawMessage{}
	json.Unmarshal(data, &taskMap)

	task := Task{}
	json.Unmarshal(*taskMap["id"], &task.Id)
	json.Unmarshal(*taskMap["status"], &task.Status)
	json.Unmarshal(*taskMap["created"], &task.Created)
	if val, ok := taskMap["result"]; ok {
		task.Result = []byte(*val)
	}
	return &task, nil
}

//Poll an the Task API. Blocks until success or failure.
//Returns result on success, an error otherwise
func (taskApi *TaskApi) Poll(id string, milliseconds time.Duration) ([]byte, error) {
	ticker := time.NewTicker(time.Millisecond * milliseconds)
	task, err := taskApi.Get(id)
	if err != nil {
		return nil, err
	}
	done := task.Completed()
	for !done {
		<-ticker.C
		task, err = taskApi.Get(id)
		if err != nil {
			return nil, err
		}
		done = task.Completed()
	}
	if task.Failed() {
		return nil, FailedTask(*task)
	}
	return task.Result, nil
}

//Poll an the Task API. Blocks until success or failure
func (taskApi *TaskApi) PollResponse(response *api.CcaResponse, milliseconds time.Duration) ([]byte, error) {
	if strings.EqualFold(response.TaskStatus, SUCCESS) {
		return response.Data, nil
	} else if strings.EqualFold(response.TaskStatus, FAILED) {
		return nil, api.CcaErrorResponse(*response)
	}
	return taskApi.Poll(response.TaskId, milliseconds)
}

//Returns true if task has failed
func (task Task) Failed() bool {
	return strings.EqualFold(task.Status, FAILED)
}

//Returns true if task was successful
func (task Task) Success() bool {
	return strings.EqualFold(task.Status, SUCCESS)
}

//Returns true if task is still executing
func (task Task) Pending() bool {
	return strings.EqualFold(task.Status, PENDING)
}

//Returns true if task has completed its execution
func (task Task) Completed() bool {
	return !task.Pending()
}
