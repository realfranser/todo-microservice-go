package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/realfranser/todo-microservice-go/internal"
)

const uuidRegEx string = `[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`

// TaskService ...
type TaskService interface {
	Create(c context.Context, description string, priority internal.Priority, dates internal.Dates) (internal.Task, error)
	Task(c context.Context, id string) (internal.Task, error)
	Update(c context.Context, id string, description string, priority internal.Priority, dates internal.Dates, isDone bool) error
}

// TaskHandler ...
type TaskHandler struct {
	svc TaskService
}

// NewTaskHandler ...
func NewTaskHandler(svc TaskService) *TaskHandler {
	return &TaskHandler{
		svc: svc,
	}
}

// Register connects the handlers to the router.
func (t *TaskHandler) Register(r fiber.Router) {
	r.Get(fmt.Sprintf("/tasks/:id<regex(%s)>", uuidRegEx), t.task)
	r.Post("/tasks", t.create)
	r.Put(fmt.Sprintf("/tasks/:id<regex(%s)>", uuidRegEx), t.update)
}

// Task is an activity that needs to be completed within a period of time.
type Task struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// CreateTasksRequest defines the request used for creating tasks.
type CreateTasksRequest struct {
	Description string `json:"description"`
	// XXX `Priority` and `Dates` are intentionally missing, to be covered in future videos
}

// CreateTasksResponse defines the response returned back after creating tasks.
type CreateTasksResponse struct {
	Task Task `json:"task"`
}

func (t *TaskHandler) create(c *fiber.Ctx) error {
	req := new(CreateTasksRequest)
	if err := c.BodyParser(req); err != nil {
		renderErrorResponse(c, "invalid request", http.StatusBadRequest)
		return err
	}

	task, err := t.svc.Create(c.Context(), req.Description, internal.PriorityNone, internal.Dates{})
	if err != nil {
		renderErrorResponse(c, "create failed", http.StatusInternalServerError)
		return err
	}

	renderResponse(c,
		&CreateTasksResponse{
			Task: Task{
				ID:          task.ID,
				Description: task.Description,
			},
		},
		http.StatusCreated)
	return nil
}

// GetTasksResponse defines the response returned back after searching one task.
type GetTasksResponse struct {
	Task Task `json:"task"`
}

func (t *TaskHandler) task(c *fiber.Ctx) error {
	id := c.Params("id")

	task, err := t.svc.Task(c.Context(), id)
	if err != nil {
		// XXX: Differentiating between NotFound and Internal errors will be covered in future episodes.
		renderErrorResponse(c, "find failed", http.StatusInternalServerError)
		return err
	}

	renderResponse(c,
		&GetTasksResponse{
			Task: Task{
				ID:          task.ID,
				Description: task.Description,
			},
		},
		http.StatusOK)
	return nil
}

// UpdateTasksRequest defines the request used for updating a task.
type UpdateTasksRequest struct {
	Description string `json:"description"`
	IsDone      bool   `json:"is_done"`
	// XXX `Priority` and `Dates` are intentionally missing, to be covered in future videos
}

func (t *TaskHandler) update(c *fiber.Ctx) error {
	var req UpdateTasksRequest
	if err := c.BodyParser(&req); err != nil {
		renderErrorResponse(c, "invalid request", http.StatusBadRequest)
	}

	id := c.Params("id")

	err := t.svc.Update(c.Context(), id, req.Description, internal.PriorityNone, internal.Dates{}, req.IsDone)
	if err != nil {
		// XXX: Differentiating between NotFound and Internal errors will be covered in future episodes.
		renderErrorResponse(c, "update failed", http.StatusInternalServerError)
		return err
	}

	renderResponse(c, &struct{}{}, http.StatusOK)
	return nil
}
