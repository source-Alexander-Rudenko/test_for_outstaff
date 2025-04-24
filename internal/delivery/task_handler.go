package delivery

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/source-Alexander-Rudenko/test_for_outstaff/internal/domain"
	"io"
	"net/http"
	"strconv"
)

type TaskUsecase interface {
	CreateTask(ctx context.Context) (*domain.Task, error)
	Get(ctx context.Context, task *domain.Task) (*domain.Task, error)
	List(ctx context.Context, offset, limit int) ([]*domain.Task, error)
	Delete(ctx context.Context, task *domain.Task) error
	Update(ctx context.Context, task *domain.Task) error
}

type TaskHandler struct {
	Usecase TaskUsecase
}

func NewTaskHandler(uc TaskUsecase) *TaskHandler {
	return &TaskHandler{uc}
}

func (h *TaskHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/task", h.CreateTask).Methods("POST")
	router.HandleFunc("/task/{id}", h.UpdateTask).Methods("PUT")
	router.HandleFunc("/task/{id}", h.DeleteTask).Methods("DELETE")
	router.HandleFunc("/task", h.ListTask).Methods("GET")
	router.HandleFunc("/task/{id}", h.GetTask).Methods("GET")
	router.HandleFunc("/ping", h.Ping).Methods("GET")
}

func (h *TaskHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]bool{"alive": true})
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	task, err := h.Usecase.CreateTask(r.Context())
	if err != nil {
		WriteError(w, WrapError(ErrFailedToCreateTask, err))
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	task := &domain.Task{ID: id}
	resp, err := h.Usecase.Get(r.Context(), task)
	if err != nil {
		WriteError(w, WrapError(ErrTaskNotFound, err))
		return
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *TaskHandler) ListTask(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	offset, _ := strconv.Atoi(q.Get("offset"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit == 0 {
		limit = 10
	}
	tasks, err := h.Usecase.List(r.Context(), offset, limit)
	if err != nil {
		WriteError(w, WrapError(ErrFailedToListTasks, err))
		return
	}
	_ = json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteError(w, WrapError(ErrFailedToReadBody, err))
		return
	}
	defer r.Body.Close()

	var task domain.Task
	if err := json.Unmarshal(body, &task); err != nil {
		WriteError(w, WrapError(ErrFailedToUnmarshal, err))
		return
	}

	task.ID = id
	if err := h.Usecase.Update(r.Context(), &task); err != nil {
		WriteError(w, WrapError(ErrFailedToUpdateTask, err))
		return
	}
	_ = json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	task := &domain.Task{ID: id}
	if err := h.Usecase.Delete(r.Context(), task); err != nil {
		WriteError(w, WrapError(ErrTaskNotFound, err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
