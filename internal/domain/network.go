package domain

type NetStats struct {
	Rx uint64
	Tx uint64
}

type NetworkService interface {
	GetInterfaceStats() (map[string]NetStats, error)
	GetInterfaceNames() []string
}

type SpeedCalculator interface {
	CalculateSpeed(prev, current NetStats) (rx, tx uint64)
	FormatSpeed(bytes uint64) string
}