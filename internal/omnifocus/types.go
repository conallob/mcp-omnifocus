package omnifocus

// Project represents an OmniFocus project
type Project struct {
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	Status                  string `json:"status"`
	Note                    string `json:"note"`
	Completed               bool   `json:"completed"`
	NumberOfTasks           int    `json:"numberOfTasks"`
	NumberOfCompletedTasks  int    `json:"numberOfCompletedTasks"`
}

// Task represents an OmniFocus task
type Task struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	Note                 string    `json:"note"`
	Completed            bool      `json:"completed"`
	Flagged              bool      `json:"flagged"`
	DueDate              *string   `json:"dueDate"`
	EstimatedMinutes     *int      `json:"estimatedMinutes"`
	Tags                 []string  `json:"tags"`
	ContainingProjectID  *string   `json:"containingProjectId"`
}

// Tag represents an OmniFocus tag
type Tag struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Available bool   `json:"available"`
}

// CreateTaskRequest represents the data needed to create a task
type CreateTaskRequest struct {
	Name             string   `json:"name"`
	Note             string   `json:"note,omitempty"`
	ProjectID        string   `json:"projectId,omitempty"`
	DueDate          string   `json:"dueDate,omitempty"`
	Flagged          bool     `json:"flagged,omitempty"`
	EstimatedMinutes int      `json:"estimatedMinutes,omitempty"`
	Tags             []string `json:"tags,omitempty"`
}

// CreateProjectRequest represents the data needed to create a project
type CreateProjectRequest struct {
	Name   string   `json:"name"`
	Note   string   `json:"note,omitempty"`
	Status string   `json:"status,omitempty"`
	Tags   []string `json:"tags,omitempty"`
}

// UpdateTaskRequest represents the data needed to update a task
type UpdateTaskRequest struct {
	ID               string  `json:"id"`
	Name             *string `json:"name,omitempty"`
	Note             *string `json:"note,omitempty"`
	Completed        *bool   `json:"completed,omitempty"`
	Flagged          *bool   `json:"flagged,omitempty"`
	DueDate          *string `json:"dueDate,omitempty"`
	EstimatedMinutes *int    `json:"estimatedMinutes,omitempty"`
}

// OperationResult represents the result of a create/update operation
type OperationResult struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}
