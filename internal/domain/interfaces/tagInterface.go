package interfaces

import "middleware/internal/domain/models"

type TagUsecase interface {
	CreateTag(tag *models.Tag) error
	SearchAllTags() (*[]models.Tag, error)
	SearchTagById(id uint) (*models.Tag, error)
	UpdateTag(id uint, tag *models.Tag) error
	DeleteTag(id uint) error
	ExistTagWithId(id uint) error
	TagsRealTime() (map[uint]map[uint]models.TagState, error)
}

type TagRepository interface {
	CreateTag(tag *models.Tag) error
	SearchAllTags() (*[]models.Tag, error)
	SearchTagById(id uint) (*models.Tag, error)
	UpdateTag(id uint, tag *models.Tag) error
	DeleteTag(id uint) error
	ExistTagWithId(id uint) error
}
