package service

import (
	"time"
	"venecraft-back/cmd/entity"
	"venecraft-back/cmd/repository"
)

type LogService interface {
	CreateLog(log *entity.Log) error
	GetLogByID(id uint64) (*entity.Log, error)
	GetAllLogs() ([]entity.Log, error)
	UpdateLog(log *entity.Log) error
	DeleteLog(id uint64) error
	CountTransactions(fromDate time.Time) (int, error)
}

type logService struct {
	logRepo repository.LogRepository
}

func NewLogService(logRepo repository.LogRepository) LogService {
	return &logService{logRepo: logRepo}
}

// CreateLog adds a new log entry
func (s *logService) CreateLog(log *entity.Log) error {
	return s.logRepo.CreateLog(log)
}

// GetLogByID retrieves a log entry by its ID
func (s *logService) GetLogByID(id uint64) (*entity.Log, error) {
	return s.logRepo.GetLogByID(id)
}

// GetAllLogs retrieves all log entries
func (s *logService) GetAllLogs() ([]entity.Log, error) {
	return s.logRepo.GetAllLogs()
}

// UpdateLog updates an existing log entry
func (s *logService) UpdateLog(log *entity.Log) error {
	return s.logRepo.UpdateLog(log)
}

// DeleteLog performs a logical delete on a log entry
func (s *logService) DeleteLog(id uint64) error {
	return s.logRepo.DeleteLog(id)
}

// CountTransactions counts the number of logs with "transaction" action within a specified date range
func (s *logService) CountTransactions(fromDate time.Time) (int, error) {
	return s.logRepo.CountTransactions("transaction", fromDate)
}
