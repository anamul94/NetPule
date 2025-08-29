package infrastructure

import (
	"fmt"

	"netpulse/internal/domain"
)

type SpeedService struct{}

func NewSpeedService() *SpeedService {
	return &SpeedService{}
}

func (s *SpeedService) CalculateSpeed(prev, current domain.NetStats) (rx, tx uint64) {
	return current.Rx - prev.Rx, current.Tx - prev.Tx
}

func (s *SpeedService) FormatSpeed(bytes uint64) string {
	speed := float64(bytes)
	if speed < 1024 {
		return fmt.Sprintf("%.0fB/s", speed)
	}
	speed /= 1024
	if speed < 1024 {
		return fmt.Sprintf("%.1fK/s", speed)
	}
	speed /= 1024
	return fmt.Sprintf("%.1fM/s", speed)
}