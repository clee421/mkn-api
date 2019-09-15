package main

import (
	"encoding/json"
	"log"
	stdHttp "net/http"

	"github.com/aaclee/mkn-api/pkg/http"
	"github.com/aaclee/mkn-api/pkg/logger"
	"github.com/gorilla/mux"
)

const (
	configPath = "./config/config.development.json"
)

func main() {
	// Getting server configs
	c, err := http.GetServerConfigs(configPath)
	if err != nil {
		log.Fatalf("Error procesing config file at: %s, %s\n", configPath, err)
	}

	// Setting up logger
	logger := logger.CreateLogger()

	// Create server configs
	config := http.Config{
		Log:  logger,
		Port: c.Port,
	}

	// Create server
	server := http.CreateServer(config)

	// Create server mux
	r := mux.NewRouter()

	r.Use(middleware)

	// Catch all route handler
	r.PathPrefix("/").HandlerFunc(catchAllHandler)

	// Listen and Serve
	err = server.ListenAndServe(r)
	if err != nil {
		log.Fatalf("Server could not start on port: %d\n", config.Port)
	}
}

func catchAllHandler(w stdHttp.ResponseWriter, r *stdHttp.Request) {
	w.WriteHeader(stdHttp.StatusNotFound)

	json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{
		Error: "Path not found!",
	})
}

func middleware(next stdHttp.Handler) stdHttp.Handler {
	return stdHttp.HandlerFunc(func(w stdHttp.ResponseWriter, r *stdHttp.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}