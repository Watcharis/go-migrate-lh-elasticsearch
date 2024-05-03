package db

import (
	"context"
	"watcharis/go-migrate-lotto-history-els/models"

	"gorm.io/gorm"
)

type timeTableRepository struct {
	db *gorm.DB
}

type TimeTableRepository interface {
	GetTimetable(ctx context.Context, fistPeriod string, currentPeriod string) ([]models.TimeTable, error)
}

func NewTimetableRepository(db *gorm.DB) TimeTableRepository {
	return &timeTableRepository{
		db: db,
	}
}

func (r *timeTableRepository) GetTimetable(ctx context.Context, fistPeriod string, currentPeriod string) ([]models.TimeTable, error) {

	var timetable []models.TimeTable

	err := r.db.Debug().
		Table("timetable t").
		Select("round_date").
		Where("t.round_date >= ? AND t.round_date  <= ?", fistPeriod, currentPeriod).
		Order("t.id DESC").
		Find(&timetable).Error
	if err != nil {
		return nil, err
	}
	return timetable, nil
}
