package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/korovindenis/shutdown-from-browser/v1/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestGetTimePOHandler(t *testing.T) {
	// Arrange
	json, _ := json.Marshal(service.Status{})
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
		input        service.Status
		expectedCode int
	}{
		{
			name:         "Check Http Code",
			input:        service.Status{},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Check Bad Mode",
			input:        service.Status{Mode: "bad mode", TimeShutDown: time.Now().Format(time.RFC3339)},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Check Bad Time",
			input:        service.Status{Mode: "shutdown", TimeShutDown: time.Now().AddDate(1, 0, 0).Format(time.RFC3339)},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Check Positive Case",
			input:        service.Status{Mode: "shutdown", TimeShutDown: time.Now().Format(time.RFC3339)},
			expectedCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(test.input)
			if err != nil {
				log.Fatalf("%s", err)
			}
			rw := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/server-power/", buf)

			// Act
			PowerHandler(rw, req)

			// Assert
			assert.Equal(t, test.expectedCode, rw.Code)
		})
	}
}
