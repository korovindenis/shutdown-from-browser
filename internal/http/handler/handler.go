package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/shutdown-from-browser/v2/api/response"
	"github.com/korovindenis/shutdown-from-browser/v2/internal/config"
	"github.com/korovindenis/shutdown-from-browser/v2/internal/domain/entity"
	"github.com/korovindenis/shutdown-from-browser/v2/internal/domain/usecase"
	"go.uber.org/zap"
)

type ComputerHandler struct {
	computerUsecase usecase.IComputerUsecase
	logger          *zap.Logger
	config          *config.Config
}

func NewComputerHandler(usecase usecase.IComputerUsecase, cfg *config.Config, logger *zap.Logger) *ComputerHandler {
	return &ComputerHandler{
		computerUsecase: usecase,
		logger:          logger,
		config:          cfg,
	}
}
func (h *ComputerHandler) SetPowerOffHandler(c *gin.Context) {
	var tmpMyPc entity.MyPc

	if err := c.ShouldBindJSON(&tmpMyPc); err != nil {
		h.logger.Error("Js ShouldBindJSON SetPowerOffHandler", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.Error("internal error"))

		return
	}

	if err := h.computerUsecase.SetPowerOff(tmpMyPc); err != nil {
		h.logger.Error("SetPowerOff SetPowerOffHandler", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.Error("internal error"))

		return
	}

	c.JSON(http.StatusOK, response.PowerOffResponse{Message: "Server is " + tmpMyPc.ModePowerOff + " on the " + tmpMyPc.TimePowerOff})
}

func (h *ComputerHandler) GetTimePoHandler(c *gin.Context) {
	res, err := h.computerUsecase.GetTimePowerOff()
	if err != nil {
		h.logger.Error("handler GetTimePoHandler", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.Error("internal error"))

		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *ComputerHandler) MainPageHandler(c *gin.Context) {
	if _, err := os.Stat(h.config.HTTPServer.TemplatesPath); os.IsNotExist(err) {
		h.logger.Error("index.html not found MainPageHandler", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.Error("internal error"))

		return
	}

	c.HTML(http.StatusOK, "index.html", nil)
}
