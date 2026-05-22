package controllers

import (
	"fmt"
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TypeClpController struct {
	TypeClpUsecase interfaces.TypeClpUsecase
}

func (mc *TypeClpController) SearchAllTypeClps(c *gin.Context) {

	clpList, err := mc.TypeClpUsecase.SearchAllTypeClps()
	if err != nil {
		fmt.Println("error.seacrh")
		return
	}

	if len(*clpList) == 0 {
		c.JSON(http.StatusNotFound, models.TypeClp{})
		return
	}

	c.JSON(http.StatusOK, clpList)
}
