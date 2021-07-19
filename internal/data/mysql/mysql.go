package mysql

import (
	data2 "github.com/rgaquino/rocinante-books/internal/data"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type strategy struct {
	db *gorm.DB
}

func New(dsn string) (data2.Strategy, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	return &strategy{
		db: db,
	}, nil
}

func (s *strategy) Create(entity data2.Entity) error {
	return s.db.Create(entity).Error
}

func (s *strategy) CreateAll(entities []data2.Entity) error {
	db := s.db

	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, v := range entities {
		err := tx.Create(v).Error
		if err == nil {
			continue
		}
		if rollback := tx.Rollback(); rollback.Error != nil {
			// TODO: log
		}
		return err
	}

	if commit := tx.Commit(); commit.Error != nil {
		return commit.Error
	}
	return nil
}
