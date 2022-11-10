package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/korovindenis/shutdown-from-browser/models"
	"github.com/stretchr/testify/assert"
)

func TestGetTimePOHandler(t *testing.T) {
	// Arrange
	json, _ := json.Marshal(models.ServerStatus{})
	expected := string(json)

	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/get-time-autopoweroff/", nil)

	// Act
	GetTimePOHandler(rw, req)

	// Assert
	assert.Equal(t, http.StatusOK, rw.Code, fmt.Sprintf("Incorrect return code. Expected %d got %d", http.StatusOK, rw.Code))
	assert.Equal(t, expected, rw.Body.String(), fmt.Sprintf("Incorrect return json. Expected %+v got %s", expected, rw.Body.String()))
}

func TestPowerHandler(t *testing.T) {
	// Arrange
	tests := []struct {
		name         string
		input        models.ServerStatus
		expectedCode int
	}{
		{
			name:         "Check Http Code",
			input:        models.ServerStatus{},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Check Bad Mode",
			input:        models.ServerStatus{Mode: "bad mode", TimeShutDown: time.Now().Format(time.RFC3339)},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Check Bad Time",
			input:        models.ServerStatus{Mode: "shutdown", TimeShutDown: time.Now().AddDate(1, 0, 0).Format(time.RFC3339)},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Check Positive Case",
			input:        models.ServerStatus{Mode: "shutdown", TimeShutDown: time.Now().Format(time.RFC3339)},
			expectedCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			json.NewEncoder(buf).Encode(test.input)

			rw := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/server-power/", buf)

			// Act
			PowerHandler(rw, req)

			// Assert
			assert.Equal(t, test.expectedCode, rw.Code)
		})
	}
}
