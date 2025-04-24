package repo

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/source-Alexander-Rudenko/test_for_outstaff/internal/domain"
	"log"
	"net/http"
	"sync"
	"time"
)

type TaskRepo struct {
	mu    sync.RWMutex
	tasks map[string]*domain.Task
}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{
		tasks: make(map[string]*domain.Task),
	}
}

func (r *TaskRepo) Save(ctx context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[task.ID] = task
	return nil
}

func (r *TaskRepo) Get(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, ok := r.tasks[task.ID]
	if !ok {
		return nil, errors.New("task not found")
	}
	log.Printf("Getting task %s", task.ID)
	return task, nil
}

func (r *TaskRepo) Delete(ctx context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	task, ok := r.tasks[task.ID]
	if !ok {
		return errors.New("task not found")
	}
	if task.Status == domain.StatusRunning {
		task.Cancel()
	}
	delete(r.tasks, task.ID)
	return nil
}

func (r *TaskRepo) Update(ctx context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[task.ID]; !ok {
		return errors.New("task not found")
	}
	r.tasks[task.ID] = task
	return nil
}

func (r *TaskRepo) List(ctx context.Context, offset, limit int) ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tasks := make([]*domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}
	if offset > len(tasks) {
		return []*domain.Task{}, nil
	}
	end := offset + limit
	if end > len(tasks) {
		end = len(tasks)
	}
	return tasks[offset:end], nil
}

func (r *TaskRepo) ExecuteTask(ctx context.Context, task *domain.Task) error {
	task.Status = domain.StatusRunning
	task.UpdatedAt = time.Now()
	err := r.Update(ctx, task)
	if err != nil {
		return err
	}
	resp, err := http.Get("http://slow-api:8081/slow")
	if err != nil {
		log.Printf("Error calling external api: %v", err)
		task.Status = domain.StatusCancelled
		task.Result = "error"
		return err
	}
	defer resp.Body.Close()

	var result domain.SlowResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error decoding slow response: %v", err)
	}

	if resp.StatusCode == http.StatusOK && result.Result == "ok" {
		task.Status = domain.StatusCompleted
		task.Result = "success"
	} else {
		task.Status = domain.StatusCancelled
		task.Result = "Failed"
	}

	task.UpdatedAt = time.Now()
	err = r.Update(ctx, task)
	if err != nil {
		return err
	}
	return nil

}
