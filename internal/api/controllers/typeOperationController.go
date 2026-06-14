package controllers

import (
	"fmt"
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TypeOperationController struct {
	TypeOperationUsecase interfaces.TypeOperationUsecase
}

func (mc *TypeOperationController) SearchAllTypeOperations(c *gin.Context) {

	clpList, err := mc.TypeOperationUsecase.SearchAllOperations()
	if err != nil {
		fmt.Println("error.seacrh")
		return
	}

	if len(*clpList) == 0 {
		c.JSON(http.StatusNotFound, models.TypeOperation{})
		return
	}

	c.JSON(http.StatusOK, clpList)
}
