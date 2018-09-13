// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/olivere/elastic/uritemplates"
)

// TasksListService retrieves the list of currently executing tasks
// on one ore more nodes in the cluster. It is part of the Task Management API
// documented at https://www.elastic.co/guide/en/elasticsearch/reference/6.2/tasks.html.
//
// It is supported as of Elasticsearch 2.3.0.
type TasksListService struct {
	client            *Client
	pretty            bool
	taskId            []string
	actions           []string
	detailed          *bool
	human             *bool
	nodeId            []string
	parentTaskId      string
	waitForCompletion *bool
	groupBy           string
	headers           http.Header
}

// NewTasksListService creates a new TasksListService.
func NewTasksListService(client *Client) *TasksListService {
	return &TasksListService{
		client: client,
	}
}

// TaskId indicates to returns the task(s) with specified id(s).
// Notice that the caller is responsible for using the correct format,
// i.e. node_id:task_number, as specified in the REST API.
func (s *TasksListService) TaskId(taskId ...string) *TasksListService {
	s.taskId = append(s.taskId, taskId...)
	return s
}

// Actions is a list of actions that should be returned. Leave empty to return all.
func (s *TasksListService) Actions(actions ...string) *TasksListService {
	s.actions = append(s.actions, actions...)
	return s
}

// Detailed indicates whether to return detailed task information (default: false).
func (s *TasksListService) Detailed(detailed bool) *TasksListService {
	s.detailed = &detailed
	return s
}

// Human indicates whether to return time and byte values in human-readable format.
func (s *TasksListService) Human(human bool) *TasksListService {
	s.human = &human
	return s
}

// NodeId is a list of node IDs or names to limit the returned information;
// use `_local` to return information from the node you're connecting to,
// leave empty to get information from all nodes.
func (s *TasksListService) NodeId(nodeId ...string) *TasksListService {
	s.nodeId = append(s.nodeId, nodeId...)
	return s
}

// ParentTaskId returns tasks with specified parent task id.
// Notice that the caller is responsible for using the correct format,
// i.e. node_id:task_number, as specified in the REST API.
func (s *TasksListService) ParentTaskId(parentTaskId string) *TasksListService {
	s.parentTaskId = parentTaskId
	return s
}

// WaitForCompletion indicates whether to wait for the matching tasks
// to complete (default: false).
func (s *TasksListService) WaitForCompletion(waitForCompletion bool) *TasksListService {
	s.waitForCompletion = &waitForCompletion
	return s
}

// GroupBy groups tasks by nodes or parent/child relationships.
// As of now, it can either be "nodes" (default) or "parents" or "none".
func (s *TasksListService) GroupBy(groupBy string) *TasksListService {
	s.groupBy = groupBy
	return s
}

// Header sets headers on the request
func (s *TasksListService) Header(name string, value string) *TasksListService {
	if s.headers == nil {
		s.headers = http.Header{}
	}
	s.headers.Add(name, value)
	return s
}

// Pretty indicates that the JSON response be indented and human readable.
func (s *TasksListService) Pretty(pretty bool) *TasksListService {
	s.pretty = pretty
	return s
}

// buildURL builds the URL for the operation.
func (s *TasksListService) buildURL() (string, url.Values, error) {
	// Build URL
	var err error
	var path string
	if len(s.taskId) > 0 {
		path, err = uritemplates.Expand("/_tasks/{task_id}", map[string]string{
			"task_id": strings.Join(s.taskId, ","),
		})
	} else {
		path = "/_tasks"
	}
	if err != nil {
		return "", url.Values{}, err
	}

	// Add query string parameters
	params := url.Values{}
	if s.pretty {
		params.Set("pretty", "true")
	}
	if len(s.actions) > 0 {
		params.Set("actions", strings.Join(s.actions, ","))
	}
	if s.detailed != nil {
		params.Set("detailed", fmt.Sprintf("%v", *s.detailed))
	}
	if s.human != nil {
		params.Set("human", fmt.Sprintf("%v", *s.human))
	}
	if len(s.nodeId) > 0 {
		params.Set("node_id", strings.Join(s.nodeId, ","))
	}
	if s.parentTaskId != "" {
		params.Set("parent_task_id", s.parentTaskId)
	}
	if s.waitForCompletion != nil {
		params.Set("wait_for_completion", fmt.Sprintf("%v", *s.waitForCompletion))
	}
	if s.groupBy != "" {
		params.Set("group_by", s.groupBy)
	}
	return path, params, nil
}

// Validate checks if the operation is valid.
func (s *TasksListService) Validate() error {
	return nil
}

// Do executes the operation.
func (s *TasksListService) Do(ctx context.Context) (*TasksListResponse, error) {
	// Check pre-conditions
	if err := s.Validate(); err != nil {
		return nil, err
	}

	// Get URL for request
	path, params, err := s.buildURL()
	if err != nil {
		return nil, err
	}

	// Get HTTP response
	res, err := s.client.PerformRequest(ctx, PerformRequestOptions{
		Method:  "GET",
		Path:    path,
		Params:  params,
		Headers: s.headers,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(TasksListResponse)
	if err := s.client.decoder.Decode(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// TasksListResponse is the response of TasksListService.Do.
type TasksListResponse struct {
	TaskFailures []*TaskOperationFailure `json:"task_failures"`
	NodeFailures []*FailedNodeException  `json:"node_failures"`
	// Nodes returns the tasks per node. The key is the node id.
	Nodes map[string]*DiscoveryNode `json:"nodes"`
}

type TaskOperationFailure struct {
	TaskId int64         `json:"task_id"` // this is a long in the Java source
	NodeId string        `json:"node_id"`
	Status string        `json:"status"`
	Reason *ErrorDetails `json:"reason"`
}

type FailedNodeException struct {
	*ErrorDetails
	NodeId string `json:"node_id"`
}

type DiscoveryNode struct {
	Name             string                 `json:"name"`
	TransportAddress string                 `json:"transport_address"`
	Host             string                 `json:"host"`
	IP               string                 `json:"ip"`
	Roles            []string               `json:"roles"` // "master", "data", or "ingest"
	Attributes       map[string]interface{} `json:"attributes"`
	// Tasks returns the tasks by its id (as a string).
	Tasks map[string]*TaskInfo `json:"tasks"`
}

// TaskInfo represents information about a currently running task.
type TaskInfo struct {
	Node               string      	     `json:"node"`
	Id                 int64       	     `json:"id"` // the task id (yes, this is a long in the Java source)
	Type               string      	     `json:"type"`
	Action             string      	     `json:"action"`
	Status             interface{} 	     `json:"status"`      // has separate implementations of Task.Status in Java for reindexing, replication, and "RawTaskStatus"
	Description        interface{} 	     `json:"description"` // same as Status
	StartTime          string      	     `json:"start_time"`
	StartTimeInMillis  int64       	     `json:"start_time_in_millis"`
	RunningTime        string      	     `json:"running_time"`
	RunningTimeInNanos int64       	     `json:"running_time_in_nanos"`
	Cancellable        bool        	     `json:"cancellable"`
	ParentTaskId       string      	     `json:"parent_task_id"` // like "YxJnVYjwSBm_AUbzddTajQ:12356"
	Headers            map[string]string `json:"headers"`
}

// StartTaskResult is used in cases where a task gets started asynchronously and
// the operation simply returnes a TaskID to watch for via the Task Management API.
type StartTaskResult struct {
	TaskId string `json:"task"`
}
