package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/kiriksik/TestTaskEffectiveMobile/config"
	"github.com/kiriksik/TestTaskEffectiveMobile/internal/database"
	"github.com/kiriksik/TestTaskEffectiveMobile/internal/models"
)

type UserService struct {
	ApiConfig *config.ApiConfig
}

func (humanService *UserService) CreateHuman(ctx context.Context, req *models.HumanRequest) (database.Human, int, error) {
	if req == nil {
		return database.Human{}, http.StatusBadRequest, fmt.Errorf("bad request")
	}

	patronymicValid := true
	if req.Patronymic == "" {
		patronymicValid = false
	}

	params, status, err := GetParamsFromAPI(req.Name)
	if err != nil {
		return database.Human{}, status, err
	}

	human, err := humanService.ApiConfig.Queries.CreateHuman(ctx,
		database.CreateHumanParams{
			Name:       req.Name,
			Surname:    req.Surname,
			Patronymic: sql.NullString{String: req.Patronymic, Valid: patronymicValid},
			Age:        int32(params.Age),
			Gender:     params.Gender,
			Country:    params.Country,
		})
	if err != nil {
		return database.Human{}, http.StatusInternalServerError, fmt.Errorf("error saving human: %s", err)
	}
	fmt.Println("saved human:", human)
	return human, http.StatusCreated, nil
}

func (humanService *UserService) GetHumanByID(ctx context.Context, id string) (database.Human, int, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return database.Human{}, http.StatusBadRequest, fmt.Errorf("bad uuid: %s", err)
	}

	human, err := humanService.ApiConfig.Queries.GetHumanByID(ctx, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.Human{}, http.StatusNotFound, fmt.Errorf("human does not exists")
		}
		return database.Human{}, http.StatusInternalServerError, fmt.Errorf("failed to get human: %s", err)
	}
	fmt.Println("request for get human:", human)
	return human, http.StatusOK, nil
}

func (humanService *UserService) GetHumans(ctx context.Context) ([]database.Human, int, error) {

	humans, err := humanService.ApiConfig.Queries.GetHumans(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []database.Human{}, http.StatusNotFound, fmt.Errorf("humans does not exists")
		}
		return []database.Human{}, http.StatusInternalServerError, fmt.Errorf("failed to get humans: %s", err)
	}
	fmt.Println("request for get humans")
	return humans, http.StatusOK, nil
}

func (humanService *UserService) DeleteHuman(ctx context.Context, id string) (database.Human, int, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return database.Human{}, http.StatusBadRequest, fmt.Errorf("bad uuid: %s", err)
	}

	human, err := humanService.ApiConfig.Queries.DeleteHuman(ctx, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.Human{}, http.StatusNotFound, fmt.Errorf("human does not exists")
		}
		return database.Human{}, http.StatusInternalServerError, fmt.Errorf("failed to delete human: %s", err)
	}
	fmt.Println("deleted human:", human)
	return human, http.StatusOK, nil
}

func (humanService *UserService) UpdateHuman(ctx context.Context, req *models.HumanRequest, id string) (database.Human, int, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return database.Human{}, http.StatusBadRequest, fmt.Errorf("bad uuid: %s", err)
	}

	if req == nil {
		return database.Human{}, http.StatusBadRequest, fmt.Errorf("bad request")
	}

	patronymicValid := true
	if req.Patronymic == "" {
		patronymicValid = false
	}

	params, status, err := GetParamsFromAPI(req.Name)
	if err != nil {
		return database.Human{}, status, err
	}

	human, err := humanService.ApiConfig.Queries.UpdateHuman(ctx,
		database.UpdateHumanParams{
			ID:         uid,
			Name:       req.Name,
			Surname:    req.Surname,
			Patronymic: sql.NullString{String: req.Patronymic, Valid: patronymicValid},
			Age:        int32(params.Age),
			Gender:     params.Gender,
			Country:    params.Country,
		})
	if err != nil {
		return database.Human{}, http.StatusInternalServerError, fmt.Errorf("error updating human: %s", err)
	}
	fmt.Println("updated human:", human)

	return human, http.StatusOK, nil
}

func GetParamsFromAPI(name string) (models.ExtraParamsResponse, int, error) {
	if name == "" {
		return models.ExtraParamsResponse{}, http.StatusBadRequest, fmt.Errorf("name cant be empty")
	}
	var respBodyAge models.AgeResponse
	var respBodyCountry models.CountryResponse
	var respBodyGender models.GenderResponse

	ageResp, err := http.Get(fmt.Sprintf("https://api.agify.io/?name=%s", name))
	if err != nil {
		return models.ExtraParamsResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to get human age: %s", err)
	}
	err = json.NewDecoder(ageResp.Body).Decode(&respBodyAge)
	defer ageResp.Body.Close()
	if err != nil {
		return models.ExtraParamsResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to decode human age: %s", err)
	}

	genderResp, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
	if err != nil {
		return models.ExtraParamsResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to get human gender: %s", err)
	}
	err = json.NewDecoder(genderResp.Body).Decode(&respBodyGender)
	defer genderResp.Body.Close()
	if err != nil {
		return models.ExtraParamsResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to decode human gender: %s", err)
	}

	countryResp, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		return models.ExtraParamsResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to get human country: %s", err)
	}
	err = json.NewDecoder(countryResp.Body).Decode(&respBodyCountry)
	defer countryResp.Body.Close()
	if err != nil {
		return models.ExtraParamsResponse{}, http.StatusInternalServerError, fmt.Errorf("failed to decode human country: %s", err)
	}
	var mostProbabilityCountry struct {
		Name        string
		Probability float64
	}
	for _, c := range respBodyCountry.Country {
		if c.Probability > mostProbabilityCountry.Probability {
			mostProbabilityCountry.Name = c.CountryID
			mostProbabilityCountry.Probability = c.Probability
		}
	}

	return models.ExtraParamsResponse{Age: respBodyAge.Age, Gender: respBodyGender.Gender, Country: mostProbabilityCountry.Name}, http.StatusOK, nil

}
