package repository

import (
	"gorm.io/gorm"
	"time"
	"venecraft-back/cmd/entity"
)

type LogRepository interface {
	CreateLog(log *entity.Log) error
	GetLogByID(id uint64) (*entity.Log, error)
	GetAllLogs() ([]entity.Log, error)
	UpdateLog(log *entity.Log) error
	DeleteLog(id uint64) error
	CountTransactions(action string, fromDate time.Time) (int, error)
}

type logRepository struct {
	db *gorm.DB
}

func NewLogRepository(db *gorm.DB) LogRepository {
	return &logRepository{db}
}

func (r *logRepository) CreateLog(log *entity.Log) error {
	return r.db.Create(log).Error
}

func (r *logRepository) GetLogByID(id uint64) (*entity.Log, error) {
	var log entity.Log
	err := r.db.First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *logRepository) GetAllLogs() ([]entity.Log, error) {
	var logs []entity.Log
	err := r.db.Find(&logs).Error
	return logs, err
}

func (r *logRepository) UpdateLog(log *entity.Log) error {
	return r.db.Save(log).Error
}

func (r *logRepository) DeleteLog(id uint64) error {
	return r.db.Delete(&entity.Log{}, id).Error
}

// CountTransactions counts logs where the action is "transaction" and timestamp is within a specified range
func (r *logRepository) CountTransactions(action string, fromDate time.Time) (int, error) {
	var count int64
	err := r.db.Model(&entity.Log{}).
		Where("action = ? AND timestamp >= ?", action, fromDate).
		Count(&count).Error
	return int(count), err
}
