package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/source-Alexander-Rudenko/test_for_outstaff/internal/domain"
	"time"
)

// TaskRepository описывает контракт для работы с репозиторием задний
type TaskRepository interface {
	Save(ctx context.Context, task *domain.Task) error
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, task *domain.Task) error
	Get(ctx context.Context, task *domain.Task) (*domain.Task, error)
	List(ctx context.Context, offset, limit int) ([]*domain.Task, error)
	ExecuteTask(ctx context.Context, task *domain.Task) error
}

type Usecase struct {
	repo TaskRepository
}

func NewTaskUsecase(taskRepo TaskRepository) *Usecase {
	return &Usecase{
		repo: taskRepo,
	}
}

func (u *Usecase) CreateTask(ctx context.Context) (*domain.Task, error) {
	id := uuid.New().String()
	taskCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	task := &domain.Task{
		ID:        id,
		Status:    domain.StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Cancel:    cancel,
	}
	err := u.repo.Save(taskCtx, task)
	if err != nil {
		return nil, err
	}
	go func() {
		err = u.repo.ExecuteTask(ctx, task)
		if err != nil {
			return
		}
	}()
	return task, nil
}

func (u *Usecase) Get(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	return u.repo.Get(ctx, task)
}
func (u *Usecase) List(ctx context.Context, offset, limit int) ([]*domain.Task, error) {
	return u.repo.List(ctx, offset, limit)
}

func (u *Usecase) Delete(ctx context.Context, task *domain.Task) error {
	return u.repo.Delete(ctx, task)
}

func (u *Usecase) Update(ctx context.Context, task *domain.Task) error {
	return u.repo.Update(ctx, task)
}
