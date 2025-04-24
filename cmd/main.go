package main

import (
	"github.com/gorilla/mux"
	"github.com/source-Alexander-Rudenko/test_for_outstaff/internal/delivery"
	"github.com/source-Alexander-Rudenko/test_for_outstaff/internal/repo"
	"github.com/source-Alexander-Rudenko/test_for_outstaff/internal/usecase"
	"log"
	"net/http"
)

func main() {

	taskRepo := repo.NewTaskRepo()
	taskUsecase := usecase.NewTaskUsecase(taskRepo)
	taskHandler := delivery.NewTaskHandler(taskUsecase)

	r := mux.NewRouter()
	r.Use(delivery.LoggingMiddleware)
	taskHandler.RegisterRoutes(r)

	addr := ":8080"
	log.Printf("Listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("failed to start server %v", err)
	}

}
