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
		fmt.Println("error.json")
		return
	}

	err = tc.TagUsecase.CreateTag(tag)
	if err != nil {
		fmt.Println("")
		return
	}

	c.JSON(http.StatusOK, "create.success")
}

func (tc TagController) SearchTagById(c *gin.Context) {
	id := c.Params.ByName("id")

	tag, err := tc.TagUsecase.SearchTagById(conversion.StringToUint(id))
	if err != nil {
		fmt.Println("error.serach")
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
		fmt.Println("error.search")
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
		fmt.Println("error.json")
		return
	}

	_, err = tc.TagUsecase.SearchTagById(conversion.StringToUint(id))
	if err != nil {
		fmt.Println("not.exist")
		return
	}

	err = tc.TagUsecase.UpdateTag(conversion.StringToUint(id), tag)
	if err != nil {
		fmt.Println("error.update")
		return
	}

	c.JSON(http.StatusOK, "update.sucess")
}

func (tc TagController) DeleteTag(c *gin.Context) {
	id := c.Params.ByName("id")

	_, err := tc.TagUsecase.SearchTagById(conversion.StringToUint(id))
	if err != nil {
		fmt.Println("not.exist")
		return
	}

	err = tc.TagUsecase.DeleteTag(conversion.StringToUint(id))
	if err != nil {
		fmt.Println("error.delete")
		return
	}

	c.JSON(http.StatusOK, "delete.success")
}

func (tc TagController) TagsRealTime(c *gin.Context) {
	Map, err := tc.TagUsecase.TagsRealTime()
	if err != nil {
		fmt.Println("not.exist")
		return
	}

	if Map == nil {
		fmt.Println("not.exist")
		return
	}

	c.JSON(http.StatusOK, Map)
}
