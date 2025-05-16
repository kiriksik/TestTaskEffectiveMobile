package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/kiriksik/TestTaskEffectiveMobile/internal/database"
)

type ApiConfig struct {
	Queries *database.Queries
}

func InitializeApiConfig() *ApiConfig {
	apiCfg := &ApiConfig{
		Queries: initializeDBQueries(),
	}
	return apiCfg
}

func initializeDBQueries() *database.Queries {
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("error in db connection: %s", err)
		return nil
	}
	dbQueries := database.New(db)
	return dbQueries
}
