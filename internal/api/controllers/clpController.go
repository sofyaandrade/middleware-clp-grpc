package controllers

import (
	"fmt"
	"middleware/internal/domain/conversion"
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CLPController struct {
	CLPUsecase interfaces.CLPUsecase
}

func (mc CLPController) NewCLP(c *gin.Context) {
	var clp *models.CLP

	err := c.ShouldBindJSON(&clp)
	if err != nil {
		fmt.Println("error.json")
		return
	}

	err = mc.CLPUsecase.CreateClp(clp)
	if err != nil {
		fmt.Println("error.create")
		return
	}

	c.JSON(http.StatusOK, "create.success")
}

func (mc CLPController) SearchClpById(c *gin.Context) {
	id := c.Params.ByName("id")

	clp, err := mc.CLPUsecase.SearchClpById(conversion.StringToUint(id))
	if err != nil {
		fmt.Println("error.search")
		return
	}

	if clp.ID == 0 {
		c.JSON(http.StatusNotFound, models.CLP{})
		return
	}

	c.JSON(http.StatusOK, clp)
}

func (mc CLPController) SearchAllClps(c *gin.Context) {

	clpList, err := mc.CLPUsecase.SearchAllClps()
	if err != nil {
		fmt.Println("error.seacrh")
		return
	}

	if len(*clpList) == 0 {
		c.JSON(http.StatusNotFound, models.CLP{})
		return
	}

	c.JSON(http.StatusOK, clpList)
}

func (mc CLPController) SearchClpByType(c *gin.Context) {
	typeID := c.Params.ByName("idType")

	clpList, err := mc.CLPUsecase.SearchClpByType(conversion.StringToUint(typeID))
	if err != nil {
		fmt.Println("")
		return
	}

	if len(*clpList) == 0 {
		c.JSON(http.StatusNotFound, models.CLP{})
		return
	}

	c.JSON(http.StatusOK, clpList)
}

func (mc CLPController) UpdateClp(c *gin.Context) {
	var clp *models.CLP
	id := c.Params.ByName("id")

	err := c.ShouldBindJSON(&clp)
	if err != nil {
		fmt.Println("error.json")
		return
	}

	_, err = mc.CLPUsecase.SearchClpById(conversion.StringToUint(id))
	if err != nil {
		fmt.Println("not.foud")
		return
	}

	err = mc.CLPUsecase.UpdateClp(conversion.StringToUint(id), clp)
	if err != nil {
		fmt.Println("error.update")
		return
	}

	c.JSON(http.StatusOK, "update.sucess")
}

func (mc CLPController) DeleteClp(c *gin.Context) {
	var clp *models.CLP
	id := c.Params.ByName("id")

	err := c.ShouldBindJSON(&clp)
	if err != nil {
		fmt.Println("error.json")
		return
	}

	_, err = mc.CLPUsecase.SearchClpById(conversion.StringToUint(id))
	if err != nil {
		fmt.Println("not.exist")
		return
	}

	err = mc.CLPUsecase.DeleteClp(conversion.StringToUint(id))
	if err != nil {
		fmt.Println("error.delete")
		return
	}

	c.JSON(http.StatusOK, "delete.success")
}
