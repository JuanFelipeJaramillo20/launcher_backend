package service

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"time"
	"venecraft-back/cmd/repository"
)

type ServerMetrics struct {
	CPUUsage     float64 `json:"cpu_usage"`
	MemoryUsage  float64 `json:"memory_usage"`
	DiskUsage    float64 `json:"disk_usage"`
	ActiveUsers  int     `json:"active_users"`
	Transactions int     `json:"transactions"`
}

type serverStatsService struct {
	userRepo repository.UserRepository
	logRepo  repository.LogRepository
}

type ServerStatsService interface {
	GetMetrics() (*ServerMetrics, error)
}

func NewServerStatsService(userRepo repository.UserRepository, logRepo repository.LogRepository) ServerStatsService {
	return &serverStatsService{userRepo: userRepo, logRepo: logRepo}
}

func (s *serverStatsService) GetMetrics() (*ServerMetrics, error) {
	// Get CPU usage
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	// Get Memory usage
	vMem, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	// Get Disk usage
	diskStat, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	// Get active users count from the database
	activeUsers, err := s.userRepo.CountActiveUsers() // Assuming CountActiveUsers method in UserRepository
	if err != nil {
		return nil, err
	}

	// Get transaction count from the database
	transactionCount, err := s.logRepo.CountTransactions("transaction", time.Now().Add(-24*time.Hour))
	if err != nil {
		return nil, err
	}

	// Assemble metrics
	return &ServerMetrics{
		CPUUsage:     cpuPercent[0],
		MemoryUsage:  vMem.UsedPercent,
		DiskUsage:    diskStat.UsedPercent,
		ActiveUsers:  activeUsers,
		Transactions: transactionCount,
	}, nil
}
