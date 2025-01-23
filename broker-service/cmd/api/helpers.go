package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("invalid json format")
	}

	return nil
}

// This function `writJSON` is responsible for writing JSON responses to the HTTP response writer
// (`http.ResponseWriter`). Here is a breakdown of what the function does:
func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	// The code snippet `out, err := json.Marshal(data)` is responsible for marshaling the `data`
	// interface into a JSON format. The `json.Marshal` function in Go converts the provided data
	// structure into a JSON representation.
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// The code snippet `if len(headers) > 0 { ... }` is checking if there are any headers passed as
	// arguments to the `writJSON` function. If the length of the `headers` slice is greater than 0, it
	// means that headers have been provided.
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	// The code snippet `w.Header().Set("Content-Type", "application/json")` is setting the "Content-Type"
	// header of the HTTP response to "application/json". This header informs the client (browser, API
	// consumer, etc.) that the content being returned in the response body is in JSON format.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}
	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}
