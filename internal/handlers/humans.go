package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kiriksik/TestTaskEffectiveMobile/config"
	_ "github.com/kiriksik/TestTaskEffectiveMobile/docs"
	"github.com/kiriksik/TestTaskEffectiveMobile/internal/models"
	service "github.com/kiriksik/TestTaskEffectiveMobile/internal/services"
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
// @Description	Создаёт человека по полям имени, фамилии и отчества(необязательно)
// @Tags	humans
// @Accept	json
// @Produce	json
// @Param	request body models.HumanRequest true "Данные человека"
// @Success	201 {object} models.HumanResponse
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

// @Summary Удаление человека
// @Description	Удаляет человека по его ID
// @Tags	humans
// @Accept	json
// @Produce	json
// @Param	humanID path string true "ID человека"
// @Success	201 {object} models.HumanResponse
// @Router /api/humans [delete]
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

// @Summary Получение человека по ID
// @Description	Возвращает данные человека по его ID
// @Tags	humans
// @Produce	json
// @Param	humanID path string true "ID человека"
// @Success	200 {object} models.HumanResponse
// @Router /api/humans/{humanID} [get]
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

// @Summary Получение списка людей
// @Description	Возвращает всех человека
// @Tags	humans
// @Produce	json
// @Success	200 {array} models.HumanResponse
// @Router /api/humans [get]
func (ah *ApiHandler) getHumans(rw http.ResponseWriter, req *http.Request) {
	humanService := service.UserService{ApiConfig: ah.ApiCfg}

	humans, status, err := humanService.GetHumans(req.Context())
	if err != nil {
		respondWithError(rw, status, err.Error())
		return
	}

	respondWithJson(rw, status, humans)
}

// @Summary Обновление человека
// @Description	Обновляет данные человека по его ID
// @Tags	humans
// @Accept	json
// @Produce	json
// @Param	humanID path string true "ID человека"
// @Param	request body models.HumanRequest true "данные человека"
// @Success	200 {object} models.HumanResponse
// @Router /api/humans/{humanID} [put]
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
