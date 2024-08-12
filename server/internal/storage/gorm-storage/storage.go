package gormstorage

import (
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) GetAll(insts interface{}) error {
	return s.db.Find(insts).Error
}

func (s *Storage) Get(inst interface{}, id int64) error {
	return s.db.First(inst, id).Error
}

func (s *Storage) Create(inst interface{}) (interface{}, error) {
	return inst, s.db.Create(inst).Error
}

func (s *Storage) Update(inst interface{}) (interface{}, error) {
	return inst, s.db.Updates(inst).Error
}

// delete from table using inst fields
func (s *Storage) Delete(inst interface{}, conds ...interface{}) error {
	return s.db.Delete(inst, conds...).Error
}
