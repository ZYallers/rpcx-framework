package mtsc

import (
	"errors"
	"github.com/ZYallers/zgin/libraries/mvcs"
	"gorm.io/gorm"
)

type Model struct {
	mvcs.Model
	Table string
	DB    func() *gorm.DB
}

func (m *Model) Find(dest interface{}, where []interface{}, fields, order string, offset, limit int) {
	db := m.DB().Table(m.Table)
	if fields != "" {
		db = db.Select(fields)
	}
	if where != nil {
		db = db.Where(where[0], where[1:]...)
	}
	if order != "" {
		db = db.Order(order)
	}
	if offset > 0 {
		db = db.Offset(offset)
	}
	if limit > 0 {
		db = db.Limit(limit)
	}
	db.Find(dest)
}

func (m *Model) FindOne(dest interface{}, where []interface{}, fields, order string) {
	m.Find(dest, where, fields, order, 0, 1)
}

func (m *Model) Save(value interface{}, id ...int) (interface{}, error) {
	if len(id) > 0 && id[0] > 0 {
		return value, m.DB().Table(m.Table).Updates(value).Error
	}
	return value, m.DB().Table(m.Table).Create(value).Error
}

func (m *Model) Update(where []interface{}, value interface{}) error {
	return m.DB().Table(m.Table).Where(where[0], where[1:]...).Updates(value).Error
}

func (m *Model) Delete(where []interface{}) error {
	if where == nil {
		return errors.New("query condition cannot be empty")
	}
	return m.DB().Table(m.Table).Where(where[0], where[1:]...).Delete(nil).Error
}

func (m *Model) Count(where []interface{}) int64 {
	var count int64
	db := m.DB().Table(m.Table)
	if where != nil {
		db = db.Where(where[0], where[1:]...)
	}
	db.Count(&count)
	return count
}
