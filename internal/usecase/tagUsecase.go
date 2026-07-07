package usecase

import (
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
	"middleware/internal/infrastructure/jobs"
)

type TagUsecase struct {
	TagReposiotry  interfaces.TagRepository
	ReloadNotifier interfaces.CLPReloadNotifier
}

func NewTagUsecase(CLPRepository interfaces.TagRepository, reloadNotifier interfaces.CLPReloadNotifier) interfaces.TagUsecase {
	return &TagUsecase{
		TagReposiotry:  CLPRepository,
		ReloadNotifier: reloadNotifier,
	}
}

func (tu *TagUsecase) CreateTag(tag *models.Tag) error {
	if err := tu.TagReposiotry.CreateTag(tag); err != nil {
		return err
	}

	tu.requestCLPReload(tag.IdClp)
	return nil
}

func (tu *TagUsecase) SearchAllTags() (*[]models.Tag, error) {
	return tu.TagReposiotry.SearchAllTags()
}

func (tu *TagUsecase) SearchTagById(id uint) (*models.Tag, error) {
	return tu.TagReposiotry.SearchTagById(id)
}

func (tu *TagUsecase) UpdateTag(id uint, tag *models.Tag) error {
	previousTag, err := tu.TagReposiotry.SearchTagById(id)
	if err != nil {
		return err
	}

	if err := tu.TagReposiotry.UpdateTag(id, tag); err != nil {
		return err
	}

	tu.requestCLPReload(previousTag.IdClp)
	tu.requestCLPReload(tag.IdClp)
	return nil
}

func (tu *TagUsecase) DeleteTag(id uint) error {
	tag, err := tu.TagReposiotry.SearchTagById(id)
	if err != nil {
		return err
	}

	if err := tu.TagReposiotry.DeleteTag(id); err != nil {
		return err
	}

	tu.requestCLPReload(tag.IdClp)
	return nil
}

func (tu *TagUsecase) ExistTagWithId(id uint) error {
	return tu.TagReposiotry.ExistTagWithId(id)
}

func (tu *TagUsecase) TagsRealTime() (map[uint]map[uint]models.TagState, error) {
	return jobs.ReadAllTagsRealTime()
}

func (tu *TagUsecase) requestCLPReload(clpID uint) {
	if tu.ReloadNotifier != nil {
		tu.ReloadNotifier.RequestCLPReload(clpID)
	}
}
