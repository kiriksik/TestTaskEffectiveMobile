package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	config "github.com/kiriksik/TestTaskEffectiveMobile/config"
	handler "github.com/kiriksik/TestTaskEffectiveMobile/internal/handlers"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("failed loading enviroment: %s", err)
	}

	cfg := config.InitializeApiConfig()

	serveMux := handler.InitializeMux(cfg)

	httpServer := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	fmt.Println("server started")
	err = httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed: %s", err)
	}

}
