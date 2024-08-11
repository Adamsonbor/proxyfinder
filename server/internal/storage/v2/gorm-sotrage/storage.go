package gormstorage

import (
	"proxyfinder/internal/storage"
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

func (s *Storage) GetAllFilter(insts interface{}, options *storage.Options) error {
	if options.Page == 0 {
		options.Page = 1
	}
	if options.PerPage == 0 {
		options.PerPage = 10
	}

	offset := (options.Page - 1) * options.PerPage

	return s.db.Limit(options.PerPage).Offset(offset).Find(insts).Error
}

// Get all rows by fild
func (self *Storage) GetAllBy(inst interface{}, field string, value interface{}) error {
	return self.db.Where(field+" = ?", value).Find(inst).Error
}

// Get all rows by multiple fields
func (self *Storage) GetAllByFields(inst interface{}, fields interface{}) error {
	return self.db.Where(fields).Find(inst).Error
}

// Get by field like email or username
func (s *Storage) GetBy(inst interface{}, field string, value interface{}) error {
	return s.db.Where(field+" = ?", value).First(inst).Error
}
