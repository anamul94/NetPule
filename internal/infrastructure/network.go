package infrastructure

import (
	"bufio"
	"os"
	"sort"
	"strconv"
	"strings"

	"netpulse/internal/domain"
)

type LinuxNetworkService struct{}

func NewLinuxNetworkService() *LinuxNetworkService {
	return &LinuxNetworkService{}
}

func (s *LinuxNetworkService) GetInterfaceStats() (map[string]domain.NetStats, error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats := make(map[string]domain.NetStats)
	scanner := bufio.NewScanner(file)

	for i := 0; i < 2; i++ {
		scanner.Scan()
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		parts := strings.Fields(strings.Replace(line, ":", " ", 1))
		if len(parts) < 10 {
			continue
		}
		iface := parts[0]
		rx, _ := strconv.ParseUint(parts[1], 10, 64)
		tx, _ := strconv.ParseUint(parts[9], 10, 64)
		stats[iface] = domain.NetStats{Rx: rx, Tx: tx}
	}
	return stats, nil
}

func (s *LinuxNetworkService) GetInterfaceNames() []string {
	stats, err := s.GetInterfaceStats()
	if err != nil {
		return nil
	}

	var interfaces []string
	for iface := range stats {
		if iface != "lo" {
			interfaces = append(interfaces, iface)
		}
	}
	sort.Strings(interfaces)
	return interfaces
}