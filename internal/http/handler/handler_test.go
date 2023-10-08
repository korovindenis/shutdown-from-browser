package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/shutdown-from-browser/v2/internal/config"
	"github.com/korovindenis/shutdown-from-browser/v2/internal/domain/entity"
	"github.com/korovindenis/shutdown-from-browser/v2/internal/http/handler"
	"github.com/korovindenis/shutdown-from-browser/v2/internal/http/handler/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type domainEntity struct {
	computerHandler *handler.ComputerHandler
	mockUsecase     *mocks.Usecase
	mockError       error
}

func setUpDomain(t *testing.T, cfg *config.Config, logger *zap.Logger) domainEntity {
	mockUsecase := mocks.New(t)
	computerHandler := handler.New(mockUsecase, cfg, logger)

	return domainEntity{
		computerHandler: computerHandler,
		mockUsecase:     mockUsecase,
	}
}

func setUpRouter(t *testing.T) *gin.Engine {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %s", err)
	}
	router := gin.Default()
	router.LoadHTMLGlob(filepath.Dir(currentDir) + "/../../web/build/index.html")
	return router
}

func setUpLogger() (*zap.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zap.InfoLevel)
	return zap.New(core), logs
}

func setUpConfig() *config.Config {
	return &config.Config{}
}

func TestHandler_Static(t *testing.T) {
	t.Run("Static Handler", func(t *testing.T) {
		t.Parallel()

		// Arrange
		r := setUpRouter(t)
		currentDir, err := os.Getwd()
		if err != nil {
			log.Print("Failed to get current directory", err)
			return
		}
		folderPath := filepath.Dir(currentDir) + "/../../web/build/static/"
		files, err := filepath.Glob(filepath.Join(folderPath+"/css/", "*.css"))
		if err != nil {
			log.Print("Failed to get files in directory:", err)
			return
		}
		firstFile := ""
		if len(files) > 0 {
			firstFile = filepath.Base(files[0])
		} else {
			log.Print("Failed to get css in directory:", err)
			return
		}
		r.Static("/static/", folderPath)

		// Act
		req, err := http.NewRequest("GET", "/static/css/"+firstFile, nil)
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandler_GetTimePoHandler(t *testing.T) {
	t.Run("Get Time PowerOff Handler", func(t *testing.T) {
		t.Parallel()

		// Arrange
		r := setUpRouter(t)
		logger, _ := setUpLogger()
		cfg := setUpConfig()
		domain := setUpDomain(t, cfg, logger)
		mockData := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
		mockDataJs, _ := json.Marshal(mockData)
		domain.mockUsecase.On("GetTimePowerOff").Return(mockData, domain.mockError).Once()
		r.GET("/api/v1/get-time-autopoweroff/", domain.computerHandler.GetTimePoHandler)

		// Act
		req, err := http.NewRequest("GET", "/api/v1/get-time-autopoweroff/", nil)
		assert.NoError(t, err)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assert
		responseData, _ := io.ReadAll(w.Body)
		assert.Equal(t, string(mockDataJs), string(responseData))
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandler_SetTimePowerOffHandler(t *testing.T) {
	mockData := entity.MyPc{
		ModePowerOff: "reboot",
		TimePowerOff: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
	}

	cases := []struct {
		name             string
		input            string
		expectStatusCode int
		expectErr        error
	}{
		{
			name:             "wrong js",
			expectStatusCode: http.StatusInternalServerError,
			expectErr:        errors.New("unexpect data"),
		},
		{
			name:             "good js",
			input:            mockData.String(),
			expectStatusCode: http.StatusOK,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run("Set Time Power OffHandler : "+tc.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			r := setUpRouter(t)
			logger, _ := setUpLogger()
			cfg := setUpConfig()
			domain := setUpDomain(t, cfg, logger)
			domain.mockUsecase.On("GetTimePowerOff").Return("", tc.expectErr).Once()
			r.POST("/api/v1/server-power/", domain.computerHandler.GetTimePoHandler)
			inputJs, _ := json.Marshal(tc.input)

			// Act
			req, err := http.NewRequest("POST", "/api/v1/server-power/", bytes.NewReader([]byte(string(inputJs))))
			assert.NoError(t, err)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectStatusCode, w.Code)
		})
	}
}
