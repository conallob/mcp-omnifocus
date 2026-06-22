package main

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/conall/mcp-omnifocus/internal/omnifocus"
	"github.com/mark3labs/mcp-go/mcp"
)

// ---------- mockClient ----------

type mockClient struct {
	projects []omnifocus.Project
	tasks    []omnifocus.Task
	tags     []omnifocus.Tag
	result   *omnifocus.OperationResult
	err      error

	lastCreateTaskReq    omnifocus.CreateTaskRequest
	lastCreateProjectReq omnifocus.CreateProjectRequest
	lastUpdateTaskReq    omnifocus.UpdateTaskRequest
	lastCompleteTaskID   string
}

func (m *mockClient) ListProjects() ([]omnifocus.Project, error) { return m.projects, m.err }
func (m *mockClient) ListTasks(projectID string) ([]omnifocus.Task, error) {
	return m.tasks, m.err
}
func (m *mockClient) ListTags() ([]omnifocus.Tag, error) { return m.tags, m.err }
func (m *mockClient) CreateTask(req omnifocus.CreateTaskRequest) (*omnifocus.OperationResult, error) {
	m.lastCreateTaskReq = req
	return m.result, m.err
}
func (m *mockClient) CreateProject(req omnifocus.CreateProjectRequest) (*omnifocus.OperationResult, error) {
	m.lastCreateProjectReq = req
	return m.result, m.err
}
func (m *mockClient) UpdateTask(req omnifocus.UpdateTaskRequest) (*omnifocus.OperationResult, error) {
	m.lastUpdateTaskReq = req
	return m.result, m.err
}
func (m *mockClient) CompleteTask(taskID string) (*omnifocus.OperationResult, error) {
	m.lastCompleteTaskID = taskID
	return m.result, m.err
}

// ---------- splitTags ----------

func TestSplitTags_Empty(t *testing.T) {
	if tags := splitTags(""); tags != nil {
		t.Errorf("expected nil, got %v", tags)
	}
}

func TestSplitTags_Single(t *testing.T) {
	tags := splitTags("home")
	if len(tags) != 1 || tags[0] != "home" {
		t.Errorf("unexpected: %v", tags)
	}
}

func TestSplitTags_Multiple(t *testing.T) {
	tags := splitTags("home,work,errands")
	if len(tags) != 3 {
		t.Fatalf("expected 3, got %d: %v", len(tags), tags)
	}
	if tags[0] != "home" || tags[1] != "work" || tags[2] != "errands" {
		t.Errorf("unexpected values: %v", tags)
	}
}

func TestSplitTags_Spaces(t *testing.T) {
	tags := splitTags("home, work , errands")
	if len(tags) != 3 || tags[1] != "work" {
		t.Errorf("unexpected: %v", tags)
	}
}

func TestSplitTags_TrailingComma(t *testing.T) {
	tags := splitTags("home,")
	if len(tags) != 1 || tags[0] != "home" {
		t.Errorf("unexpected: %v", tags)
	}
}

// ---------- handleListProjects ----------

func TestHandleListProjects_ReturnsProjects(t *testing.T) {
	m := &mockClient{
		projects: []omnifocus.Project{
			{ID: "p1", Name: "Alpha", Status: "active"},
			{ID: "p2", Name: "Beta", Status: "on-hold"},
		},
	}
	res, err := handleListProjects(m, map[string]interface{}{})
	if err != nil || res.IsError {
		t.Fatalf("err=%v isError=%v", err, res.IsError)
	}
}

func TestHandleListProjects_Filter(t *testing.T) {
	m := &mockClient{
		projects: []omnifocus.Project{
			{ID: "p1", Name: "Active", Status: "active"},
			{ID: "p2", Name: "OnHold", Status: "on-hold"},
		},
	}
	res, err := handleListProjects(m, map[string]interface{}{"filter": "active"})
	if err != nil || res.IsError {
		t.Fatalf("err=%v isError=%v", err, res.IsError)
	}
	// Extract text from first content element
	text := extractText(t, res)
	var projects []omnifocus.Project
	json.Unmarshal([]byte(text), &projects)
	if len(projects) != 1 || projects[0].Name != "Active" {
		t.Errorf("filter did not work: %v", projects)
	}
}

func TestHandleListProjects_Error(t *testing.T) {
	m := &mockClient{err: errors.New("OmniFocus unavailable")}
	res, err := handleListProjects(m, map[string]interface{}{})
	if err != nil {
		t.Fatalf("handler should not return Go error: %v", err)
	}
	if !res.IsError {
		t.Error("expected IsError=true on client failure")
	}
}

// ---------- handleListTasks ----------

func TestHandleListTasks_All(t *testing.T) {
	m := &mockClient{
		tasks: []omnifocus.Task{{ID: "t1", Name: "Task one"}},
	}
	res, err := handleListTasks(m, map[string]interface{}{})
	if err != nil || res.IsError {
		t.Fatalf("err=%v isError=%v", err, res.IsError)
	}
}

func TestHandleListTasks_WithProjectID(t *testing.T) {
	m := &mockClient{tasks: []omnifocus.Task{}}
	res, err := handleListTasks(m, map[string]interface{}{"project_id": "proj-1"})
	if err != nil || res.IsError {
		t.Fatalf("err=%v isError=%v", err, res.IsError)
	}
}

func TestHandleListTasks_Error(t *testing.T) {
	m := &mockClient{err: errors.New("fail")}
	res, err := handleListTasks(m, map[string]interface{}{})
	if err != nil || !res.IsError {
		t.Errorf("expected IsError=true, got err=%v isError=%v", err, res.IsError)
	}
}

// ---------- handleListTags ----------

func TestHandleListTags_ReturnsTags(t *testing.T) {
	m := &mockClient{
		tags: []omnifocus.Tag{{ID: "tag1", Name: "home"}, {ID: "tag2", Name: "work"}},
	}
	res, err := handleListTags(m, map[string]interface{}{})
	if err != nil || res.IsError {
		t.Fatalf("err=%v isError=%v", err, res.IsError)
	}
	text := extractText(t, res)
	if !strings.Contains(text, "home") {
		t.Errorf("expected 'home' in output: %s", text)
	}
}

func TestHandleListTags_Error(t *testing.T) {
	m := &mockClient{err: errors.New("fail")}
	res, err := handleListTags(m, map[string]interface{}{})
	if err != nil || !res.IsError {
		t.Errorf("expected IsError=true")
	}
}

// ---------- handleCreateTask ----------

func TestHandleCreateTask_BasicInbox(t *testing.T) {
	m := &mockClient{result: &omnifocus.OperationResult{ID: "t1", Name: "Buy milk", Success: true}}
	res, err := handleCreateTask(m, map[string]interface{}{"name": "Buy milk"})
	if err != nil || res.IsError {
		t.Fatalf("err=%v isError=%v", err, res.IsError)
	}
	if m.lastCreateTaskReq.Name != "Buy milk" {
		t.Errorf("unexpected name: %q", m.lastCreateTaskReq.Name)
	}
}

func TestHandleCreateTask_AllFields(t *testing.T) {
	m := &mockClient{result: &omnifocus.OperationResult{ID: "t1", Name: "T", Success: true}}
	args := map[string]interface{}{
		"name":               "T",
		"note":               "A note",
		"project_id":         "p1",
		"due_date":           "2025-12-31T23:59:59Z",
		"flagged":            true,
		"estimated_minutes":  float64(30),
		"tags":               "home, work",
	}
	res, err := handleCreateTask(m, args)
	if err != nil || res.IsError {
		t.Fatalf("err=%v isError=%v", err, res.IsError)
	}
	req := m.lastCreateTaskReq
	if req.Note != "A note" || req.ProjectID != "p1" || !req.Flagged ||
		req.EstimatedMinutes != 30 || len(req.Tags) != 2 {
		t.Errorf("field mapping wrong: %+v", req)
	}
}

func TestHandleCreateTask_Error(t *testing.T) {
	m := &mockClient{err: errors.New("create failed")}
	res, err := handleCreateTask(m, map[string]interface{}{"name": "T"})
	if err != nil || !res.IsError {
		t.Errorf("expected IsError=true")
	}
}

// ---------- handleCreateProject ----------

func TestHandleCreateProject_Basic(t *testing.T) {
	m := &mockClient{result: &omnifocus.OperationResult{ID: "p1", Name: "Work", Success: true}}
	res, err := handleCreateProject(m, map[string]interface{}{"name": "Work"})
	if err != nil || res.IsError {
		t.Fatalf("err=%v isError=%v", err, res.IsError)
	}
	if m.lastCreateProjectReq.Name != "Work" {
		t.Errorf("unexpected name: %q", m.lastCreateProjectReq.Name)
	}
}

func TestHandleCreateProject_WithStatusAndTags(t *testing.T) {
	m := &mockClient{result: &omnifocus.OperationResult{ID: "p1", Name: "P", Success: true}}
	args := map[string]interface{}{
		"name":   "P",
		"note":   "N",
		"status": "active",
		"tags":   "work,home",
	}
	res, err := handleCreateProject(m, args)
	if err != nil || res.IsError {
		t.Fatalf("err=%v isError=%v", err, res.IsError)
	}
	req := m.lastCreateProjectReq
	if req.Status != "active" || req.Note != "N" || len(req.Tags) != 2 {
		t.Errorf("field mapping wrong: %+v", req)
	}
}

func TestHandleCreateProject_Error(t *testing.T) {
	m := &mockClient{err: errors.New("fail")}
	res, err := handleCreateProject(m, map[string]interface{}{"name": "P"})
	if err != nil || !res.IsError {
		t.Errorf("expected IsError=true")
	}
}

// ---------- handleUpdateTask ----------

func TestHandleUpdateTask_NameAndFlag(t *testing.T) {
	m := &mockClient{result: &omnifocus.OperationResult{ID: "t1", Name: "New", Success: true}}
	args := map[string]interface{}{
		"id":      "t1",
		"name":    "New",
		"flagged": true,
	}
	res, err := handleUpdateTask(m, args)
	if err != nil || res.IsError {
		t.Fatalf("err=%v isError=%v", err, res.IsError)
	}
	req := m.lastUpdateTaskReq
	if req.ID != "t1" || req.Name == nil || *req.Name != "New" ||
		req.Flagged == nil || !*req.Flagged {
		t.Errorf("field mapping wrong: %+v", req)
	}
}

func TestHandleUpdateTask_AllOptionals(t *testing.T) {
	m := &mockClient{result: &omnifocus.OperationResult{ID: "t1", Name: "T", Success: true}}
	args := map[string]interface{}{
		"id":                "t1",
		"note":              "note",
		"completed":         true,
		"due_date":          "2025-01-01T00:00:00Z",
		"estimated_minutes": float64(60),
	}
	res, err := handleUpdateTask(m, args)
	if err != nil || res.IsError {
		t.Fatalf("err=%v isError=%v", err, res.IsError)
	}
	req := m.lastUpdateTaskReq
	if req.Note == nil || *req.Note != "note" || req.Completed == nil || !*req.Completed ||
		req.EstimatedMinutes == nil || *req.EstimatedMinutes != 60 {
		t.Errorf("field mapping wrong: %+v", req)
	}
}

func TestHandleUpdateTask_Error(t *testing.T) {
	m := &mockClient{err: errors.New("fail")}
	res, err := handleUpdateTask(m, map[string]interface{}{"id": "t1"})
	if err != nil || !res.IsError {
		t.Errorf("expected IsError=true")
	}
}

// ---------- handleCompleteTask ----------

func TestHandleCompleteTask_Success(t *testing.T) {
	m := &mockClient{result: &omnifocus.OperationResult{ID: "t1", Name: "Done", Success: true}}
	res, err := handleCompleteTask(m, map[string]interface{}{"id": "t1"})
	if err != nil || res.IsError {
		t.Fatalf("err=%v isError=%v", err, res.IsError)
	}
	if m.lastCompleteTaskID != "t1" {
		t.Errorf("expected task ID 't1', got %q", m.lastCompleteTaskID)
	}
}

func TestHandleCompleteTask_Error(t *testing.T) {
	m := &mockClient{err: errors.New("fail")}
	res, err := handleCompleteTask(m, map[string]interface{}{"id": "t1"})
	if err != nil || !res.IsError {
		t.Errorf("expected IsError=true")
	}
}

// ---------- helper ----------

// extractText serialises a CallToolResult and pulls the text from the first
// content element. Uses JSON round-trip to avoid depending on internal types.
func extractText(t *testing.T, res *mcp.CallToolResult) string {
	t.Helper()
	b, err := json.Marshal(res)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var wrapper struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(b, &wrapper); err != nil || len(wrapper.Content) == 0 {
		t.Fatalf("unmarshal content: %v (raw=%s)", err, b)
	}
	return wrapper.Content[0].Text
}
