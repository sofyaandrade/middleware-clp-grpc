package usecase

import (
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
	"middleware/internal/infrastructure/jobs"
)

type TagUsecase struct {
	TagReposiotry interfaces.TagRepository
}

func NewTagUsecase(CLPRepository interfaces.TagRepository) interfaces.TagUsecase {
	return &TagUsecase{
		TagReposiotry: CLPRepository,
	}
}

func (tu *TagUsecase) CreateTag(tag *models.Tag) error {
	return tu.TagReposiotry.CreateTag(tag)
}

func (tu *TagUsecase) SearchAllTags() (*[]models.Tag, error) {
	return tu.TagReposiotry.SearchAllTags()
}

func (tu *TagUsecase) SearchTagById(id uint) (*models.Tag, error) {
	return tu.TagReposiotry.SearchTagById(id)
}

func (tu *TagUsecase) UpdateTag(id uint, tag *models.Tag) error {
	return tu.TagReposiotry.UpdateTag(id, tag)
}

func (tu *TagUsecase) DeleteTag(id uint) error {
	return tu.TagReposiotry.DeleteTag(id)
}

func (tu *TagUsecase) ExistTagWithId(id uint) error {
	return tu.TagReposiotry.ExistTagWithId(id)
}

func (tu *TagUsecase) TagsRealTime() (map[uint]map[uint]interface{}, error) {
	return jobs.ReadAllTagsRealTime()
}
