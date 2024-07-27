package gormstorage

import (
	gormstoragev1 "proxyfinder/internal/storage/gorm-storage"

	"gorm.io/gorm"
)

type Storage struct {
	*gormstoragev1.Storage
	db *gorm.DB
}

func New(db *gorm.DB) *Storage {
	return &Storage{
		Storage: gormstoragev1.New(db),
		db:      db,
	}
}

func (s *Storage) GetAllFilter(insts interface{}, page int, perPage int) error {
	if page == 0 {
		page = 1
	}
	if perPage == 0 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	return s.db.Limit(perPage).Offset(offset).Find(insts).Error
}
