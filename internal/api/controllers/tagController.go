package controllers

import (
	"fmt"
	"middleware/internal/domain/conversion"
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TagController struct {
	TagUsecase interfaces.TagUsecase
}

func (tc TagController) NewTag(c *gin.Context) {
	var tag *models.Tag

	err := c.ShouldBindJSON(&tag)
	if err != nil {
		fmt.Println("Dados da requisição inválidos.")
		return
	}

	err = tc.TagUsecase.CreateTag(tag)
	if err != nil {
		fmt.Println("")
		return
	}

	c.JSON(http.StatusOK, "Cadastro realizado com sucesso.")
}

func (tc TagController) SearchTagById(c *gin.Context) {
	id := c.Params.ByName("id")

	tag, err := tc.TagUsecase.SearchTagById(conversion.StringToUint(id))
	if err != nil {
		fmt.Println("Não foi possível buscar a tag.")
		return
	}

	if tag.ID == 0 {
		c.JSON(http.StatusNotFound, models.Tag{})
		return
	}

	c.JSON(http.StatusOK, tag)
}

func (tc TagController) SearchAllTags(c *gin.Context) {

	tagList, err := tc.TagUsecase.SearchAllTags()
	if err != nil {
		fmt.Println("Não foi possível buscar a tag.")
		return
	}

	if len(*tagList) == 0 {
		c.JSON(http.StatusNotFound, models.Tag{})
		return
	}

	c.JSON(http.StatusOK, tagList)
}

func (tc TagController) UpdateTag(c *gin.Context) {
	var tag *models.Tag
	id := c.Params.ByName("id")

	err := c.ShouldBindJSON(&tag)
	if err != nil {
		fmt.Println("Dados da requisição inválidos.")
		return
	}

	_, err = tc.TagUsecase.SearchTagById(conversion.StringToUint(id))
	if err != nil {
		fmt.Println("Tag não encontrada.")
		return
	}

	err = tc.TagUsecase.UpdateTag(conversion.StringToUint(id), tag)
	if err != nil {
		fmt.Println("Não foi possível atualizar a tag.")
		return
	}

	c.JSON(http.StatusOK, "Atualização realizada com sucesso.")
}

func (tc TagController) DeleteTag(c *gin.Context) {
	id := c.Params.ByName("id")

	_, err := tc.TagUsecase.SearchTagById(conversion.StringToUint(id))
	if err != nil {
		fmt.Println("Tag não encontrada.")
		return
	}

	err = tc.TagUsecase.DeleteTag(conversion.StringToUint(id))
	if err != nil {
		fmt.Println("Não foi possível excluir a tag.")
		return
	}

	c.JSON(http.StatusOK, "Exclusão realizada com sucesso.")
}

func (tc TagController) TagsRealTime(c *gin.Context) {
	Map, err := tc.TagUsecase.TagsRealTime()
	if err != nil {
		fmt.Println("Tag não encontrada.")
		return
	}

	if Map == nil {
		fmt.Println("Tag não encontrada.")
		return
	}

	c.JSON(http.StatusOK, Map)
}
