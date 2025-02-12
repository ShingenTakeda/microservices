package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

func (app *Config) broker(w http.ResponseWriter, r *http.Request) {
	payload := JsonResponse{
		Error:   false,
		Message: "Hit the broker!",
	}

	_ = app.write_json(w, http.StatusOK, payload)
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// Create json to send to auth

	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// Call auth service

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.error_json(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.error_json(w, err)
		return
	}
	defer response.Body.Close()

	// Make sure that we get the correct code
	if response.StatusCode == http.StatusUnauthorized {
		app.error_json(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.error_json(w, errors.New("error calling auth service"))
		return
	}

	// Variable to store reponse Body
	var jsonFromService JsonResponse

	// Decode json from auth
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.error_json(w, errors.New("error calling auth service"))
		return
	}

	if jsonFromService.Error {
		app.error_json(w, err, http.StatusUnauthorized)
		return
	}

	var payload JsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	app.write_json(w, http.StatusAccepted, payload)
}

func (app *Config) log_item(w http.ResponseWriter, entry LogPayload) {
	json_data, _ := json.MarshalIndent(entry, "", "\t")

	log_service_url := "http://logger-service/log"

	request, err := http.NewRequest("POST", log_service_url, bytes.NewBuffer(json_data))
	if err != nil {
		app.error_json(w, err)
		return
	}

	request.Header.Set("Conten-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.error_json(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.error_json(w, err)
		return
	}

	var payload JsonResponse
	payload.Error = false
	payload.Message = "logged"

	app.write_json(w, http.StatusAccepted, payload)
}

func (app *Config) handle_submission(w http.ResponseWriter, r *http.Request) {
	var request_payload RequestPayload

	err := app.read_json(w, r, &request_payload)
	if err != nil {
		app.error_json(w, err)
		return
	}

	switch request_payload.Action {
	case "auth":
		app.authenticate(w, request_payload.Auth)
	case "log":
		app.log_item(w, request_payload.Log)
	default:
		app.error_json(w, errors.New("unkown action"))
	}
}
