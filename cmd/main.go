package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	config "github.com/kiriksik/TestTaskEffectiveMobile/config"
	_ "github.com/kiriksik/TestTaskEffectiveMobile/docs"
	handler "github.com/kiriksik/TestTaskEffectiveMobile/internal/handlers"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Humans API
// @version         1.0
// @description     API для управления людьми

// @host      localhost:8080
// @BasePath  /

// @schemes http
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("failed loading enviroment: %s", err)
	}

	cfg := config.InitializeApiConfig()

	serveMux := handler.InitializeMux(cfg)
	serveMux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

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
