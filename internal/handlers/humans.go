package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kiriksik/TestTaskEffectiveMobile/config"
	_ "github.com/kiriksik/TestTaskEffectiveMobile/docs"
	"github.com/kiriksik/TestTaskEffectiveMobile/internal/models"
	service "github.com/kiriksik/TestTaskEffectiveMobile/internal/services"
	swaggerFiles "github.com/swaggo/files"
	httpSwagger "github.com/swaggo/http-swagger"
)

type ApiHandler struct {
	ApiCfg *config.ApiConfig
}

type responseError struct {
	Err string `json:"error"`
}

func InitializeMux(ac *config.ApiConfig) *http.ServeMux {

	ah := &ApiHandler{ApiCfg: ac}
	serveMux := http.NewServeMux()

	serveMux.HandleFunc("GET /api/humans/{humanID}", ah.getHumanByID)
	serveMux.HandleFunc("POST /api/humans", ah.createHuman)
	serveMux.HandleFunc("GET /api/humans", ah.getHumans)
	serveMux.HandleFunc("PUT /api/humans/{humanID}", ah.updateHuman)
	serveMux.HandleFunc("DELETE /api/humans/{humanID}", ah.deleteHuman)
	serveMux.HandleFunc("GET /swagger/*any", httpSwagger.WrapHandler(swaggerFiles.Handler))
	return serveMux
}

func respondWithError(rw http.ResponseWriter, code int, errorMessage string) {

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)

	responseError := responseError{Err: errorMessage}
	jsonErr, _ := json.Marshal(responseError)

	rw.Write(jsonErr)
}

func respondWithJson(rw http.ResponseWriter, code int, payload interface{}) {

	rw.Header().Set("Content-Type", "application/json")

	encodedJson, err := json.Marshal(payload)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, fmt.Sprintf("error marshalling json: %s", err))
		return
	}

	rw.WriteHeader(code)
	rw.Write(encodedJson)
}

// @Summary Создание человека
// @Description	Создаёт пользователя по полям имени, фамилии и отчества(необязательно)
// @Tags	humans
// @Accept	json
// @Produce	json
// @Param	request body models.HumanRequest true "query params"
// @Success	201 {object} database.Human
// @Router /api/humans [post]
func (ah *ApiHandler) createHuman(rw http.ResponseWriter, req *http.Request) {

	humanService := service.UserService{ApiConfig: ah.ApiCfg}
	var reqBodyData models.HumanRequest

	err := json.NewDecoder(req.Body).Decode(&reqBodyData)
	defer req.Body.Close()
	if err != nil {
		respondWithError(rw, http.StatusBadRequest, fmt.Sprintf("error marshalling json: %s", err))
		return
	}

	human, status, err := humanService.CreateHuman(req.Context(), &reqBodyData)
	if err != nil {
		respondWithError(rw, status, err.Error())
		return
	}

	respondWithJson(rw, status, human)
}

func (ah *ApiHandler) deleteHuman(rw http.ResponseWriter, req *http.Request) {
	humanService := service.UserService{ApiConfig: ah.ApiCfg}
	humanID := req.PathValue("humanID")
	if humanID == "" {
		respondWithError(rw, http.StatusBadRequest, "missing id")
		return
	}

	human, status, err := humanService.DeleteHuman(req.Context(), humanID)
	if err != nil {
		respondWithError(rw, status, err.Error())
		return
	}

	respondWithJson(rw, status, human)
}

func (ah *ApiHandler) getHumanByID(rw http.ResponseWriter, req *http.Request) {
	humanService := service.UserService{ApiConfig: ah.ApiCfg}
	humanID := req.PathValue("humanID")
	if humanID == "" {
		respondWithError(rw, http.StatusBadRequest, "missing id")
		return
	}
	// fmt.Println(humanID)
	human, status, err := humanService.GetHumanByID(req.Context(), humanID)
	if err != nil {
		respondWithError(rw, status, err.Error())
		return
	}

	respondWithJson(rw, status, human)
}

func (ah *ApiHandler) getHumans(rw http.ResponseWriter, req *http.Request) {
	humanService := service.UserService{ApiConfig: ah.ApiCfg}

	humans, status, err := humanService.GetHumans(req.Context())
	if err != nil {
		respondWithError(rw, status, err.Error())
		return
	}

	respondWithJson(rw, status, humans)
}
func (ah *ApiHandler) updateHuman(rw http.ResponseWriter, req *http.Request) {
	humanService := service.UserService{ApiConfig: ah.ApiCfg}
	var reqBodyData models.HumanRequest
	humanID := req.PathValue("humanID")
	if humanID == "" {
		respondWithError(rw, http.StatusBadRequest, "missing id")
		return
	}

	err := json.NewDecoder(req.Body).Decode(&reqBodyData)
	defer req.Body.Close()
	if err != nil {
		respondWithError(rw, http.StatusBadRequest, fmt.Sprintf("error marshalling json: %s", err))
		return
	}

	human, status, err := humanService.UpdateHuman(req.Context(), &reqBodyData, humanID)
	if err != nil {
		respondWithError(rw, status, err.Error())
		return
	}

	respondWithJson(rw, status, human)
}
