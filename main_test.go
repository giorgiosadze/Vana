package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginRoute(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		payload        map[string]string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Login",
			method:         http.MethodPost,
			payload:        map[string]string{"username": "testuser", "password": "password123"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Login successful"}`,
		},
		{
			name:           "Invalid Credentials",
			method:         http.MethodPost,
			payload:        map[string]string{"username": "wronguser", "password": "wrongpass"},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid credentials\n",
		},
		{
			name:           "Missing Fields",
			method:         http.MethodPost,
			payload:        map[string]string{"username": ""},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid input\n",
		},
		{
			name:           "Wrong HTTP Method",
			method:         http.MethodGet,
			payload:        nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method not allowed\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body *bytes.Buffer
			if tt.payload != nil {
				jsonPayload, _ := json.Marshal(tt.payload)
				body = bytes.NewBuffer(jsonPayload)
			} else {
				body = &bytes.Buffer{}
			}

			req := httptest.NewRequest(tt.method, "/login", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			LoginHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			responseBody := new(bytes.Buffer)
			responseBody.ReadFrom(res.Body)

			if responseBody.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, responseBody.String())
			}
		})
	}
}
