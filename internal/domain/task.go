package domain

import (
	"context"
	"time"
)

type TaskStatus string

const (
	StatusPending   TaskStatus = "PENDING"
	StatusRunning   TaskStatus = "RUNNING"
	StatusCompleted TaskStatus = "COMPLETED"
	StatusCancelled TaskStatus = "CANCELLED"
)

type Task struct {
	ID        string             `json:"id"`
	Status    TaskStatus         `json:"status"`
	Result    string             `json:"result"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	Cancel    context.CancelFunc `json:"-"`
}

type SlowResponse struct {
	Result string `json:"result"`
}
