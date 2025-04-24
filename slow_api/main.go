package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/slow", slowHandler)

	log.Println("Slow API listening on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	delay := time.Duration(3+rand.Intn(3)) * time.Second
	time.Sleep(delay)

	// Случайный результат: ok или fail
	var result string
	if rand.Intn(2) == 0 {
		result = "ok"
	} else {
		result = "fail"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"result": result,
	})
	log.Printf("Handled /slow - result=%s after %v\n", result, time.Since(start))
}
