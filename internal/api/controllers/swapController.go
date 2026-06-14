package controllers

import (
	"fmt"
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SwapController struct {
	SwapUsecase interfaces.SwapUsecase
}

func (mc *SwapController) SearchAllSwaps(c *gin.Context) {

	clpList, err := mc.SwapUsecase.SearchAllSwaps()
	if err != nil {
		fmt.Println("error.seacrh")
		return
	}

	if len(*clpList) == 0 {
		c.JSON(http.StatusNotFound, models.Swap{})
		return
	}

	c.JSON(http.StatusOK, clpList)
}
