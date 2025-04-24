package delivery

import (
	"encoding/json"
	"net/http"
)

// AppError представляет структуру для описания ошибок API
type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

var (
	ErrFailedToCreateTask = &AppError{Code: http.StatusInternalServerError, Message: "failed to create task"}
	ErrTaskNotFound       = &AppError{Code: http.StatusNotFound, Message: "task not found"}
	ErrFailedToListTasks  = &AppError{Code: http.StatusInternalServerError, Message: "failed to list tasks"}
	ErrFailedToReadBody   = &AppError{Code: http.StatusBadRequest, Message: "failed to read body"}
	ErrFailedToUnmarshal  = &AppError{Code: http.StatusBadRequest, Message: "failed to parse JSON body"}
	ErrFailedToUpdateTask = &AppError{Code: http.StatusInternalServerError, Message: "failed to update task"}
)

func WriteError(w http.ResponseWriter, appErr *AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": appErr.Message})
}

func WrapError(base *AppError, err error) *AppError {
	return &AppError{
		Code:    base.Code,
		Message: base.Message,
		Err:     err,
	}
}
