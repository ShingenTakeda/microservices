package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) read_json(w http.ResponseWriter, r *http.Request, data any) error {
	max_bytes := 1048576 // 1MB

	r.Body = http.MaxBytesReader(w, r.Body, int64(max_bytes))

	json_decoder := json.NewDecoder(r.Body)
	err := json_decoder.Decode(data)
	if err != nil {
		return err
	}

	err = json_decoder.Decode(&struct{}{})

	if err != io.EOF {
		return errors.New("body must have only a single json value")
	}

	return nil
}

func (app *Config) write_json(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) error_json(w http.ResponseWriter, err error, status ...int) error {
	status_code := http.StatusBadRequest

	if len(status) > 0 {
		status_code = status[0]
	}

	var payload JsonResponse

	payload.Error = true
	payload.Message = err.Error()

	return app.write_json(w, status_code, payload)
}
