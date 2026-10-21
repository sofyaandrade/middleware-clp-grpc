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
		fmt.Println("Dados da requisição inválidos.")
		return
	}

	err = mc.CLPUsecase.CreateClp(clp)
	if err != nil {
		fmt.Println("Não foi possível criar o CLP.")
		return
	}

	c.JSON(http.StatusOK, "Cadastro realizado com sucesso.")
}

func (mc CLPController) SearchClpById(c *gin.Context) {
	id := c.Params.ByName("id")

	clp, err := mc.CLPUsecase.SearchClpById(conversion.StringToUint(id))
	if err != nil {
		fmt.Println("Não foi possível buscar o CLP.")
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
		fmt.Println("Não foi possível buscar o CLP.")
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

func (mc CLPController) ClpsStatus(c *gin.Context) {
	c.JSON(http.StatusOK, mc.CLPUsecase.ClpsStatus())
}

func (mc CLPController) UpdateClp(c *gin.Context) {
	var clp *models.CLP
	id := c.Params.ByName("id")

	err := c.ShouldBindJSON(&clp)
	if err != nil {
		fmt.Println("Dados da requisição inválidos.")
		return
	}

	_, err = mc.CLPUsecase.SearchClpById(conversion.StringToUint(id))
	if err != nil {
		fmt.Println("CLP não encontrado.")
		return
	}

	err = mc.CLPUsecase.UpdateClp(conversion.StringToUint(id), clp)
	if err != nil {
		fmt.Println("Não foi possível atualizar o CLP.")
		return
	}

	c.JSON(http.StatusOK, "Atualização realizada com sucesso.")
}

func (mc CLPController) DeleteClp(c *gin.Context) {
	id := c.Params.ByName("id")

	_, err := mc.CLPUsecase.SearchClpById(conversion.StringToUint(id))
	if err != nil {
		fmt.Println("CLP não encontrado.")
		return
	}

	err = mc.CLPUsecase.DeleteClp(conversion.StringToUint(id))
	if err != nil {
		fmt.Println("Não foi possível excluir o CLP.")
		return
	}

	c.JSON(http.StatusOK, "Exclusão realizada com sucesso.")
}
