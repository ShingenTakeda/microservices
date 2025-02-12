package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) LogRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	json_data, _ := json.MarshalIndent(entry, "", "\t")
	log_service_url := "http://logger-service/log"

	request, err := http.NewRequest("POST", log_service_url, bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) authenticate(w http.ResponseWriter, r *http.Request) {
	var request_payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.read_json(w, r, &request_payload)
	if err != nil {
		app.error_json(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(request_payload.Email)
	if err != nil {
		app.error_json(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(request_payload.Password)
	if err != nil || !valid {
		app.error_json(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// Log authentication
	err = app.LogRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.error_json(w, err)
	}

	payload := JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.write_json(w, http.StatusAccepted, payload)
}
