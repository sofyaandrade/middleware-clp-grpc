package controllers

import (
	"fmt"
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TypeTagController struct {
	TypeTagUsecase interfaces.TypeTagUsecase
}

func (mc *TypeTagController) SearchAllTypeTags(c *gin.Context) {

	clpList, err := mc.TypeTagUsecase.SearchAllTypeTags()
	if err != nil {
		fmt.Println("error.seacrh")
		return
	}

	if len(*clpList) == 0 {
		c.JSON(http.StatusNotFound, models.TypeTag{})
		return
	}

	c.JSON(http.StatusOK, clpList)
}
