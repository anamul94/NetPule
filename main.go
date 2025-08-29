package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/getlantern/systray"
)

type NetStats struct {
	Rx uint64
	Tx uint64
}

var (
	selectedInterface = "all"
	interfaceItems    = make(map[string]*systray.MenuItem)
)

func parseProcNetDev() (map[string]NetStats, error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats := make(map[string]NetStats)
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
		stats[iface] = NetStats{Rx: rx, Tx: tx}
	}
	return stats, nil
}

func formatSpeed(bytes uint64) string {
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

func updateDisplay(prevStats, curStats map[string]NetStats) {
	var totalRx, totalTx uint64

	if selectedInterface == "all" {
		for iface, cur := range curStats {
			if iface == "lo" {
				continue
			}
			if prev, ok := prevStats[iface]; ok {
				totalRx += cur.Rx - prev.Rx
				totalTx += cur.Tx - prev.Tx
			}
		}
	} else {
		if cur, ok := curStats[selectedInterface]; ok {
			if prev, ok := prevStats[selectedInterface]; ok {
				totalRx = cur.Rx - prev.Rx
				totalTx = cur.Tx - prev.Tx
			}
		}
	}

	downStr := formatSpeed(totalRx)
	upStr := formatSpeed(totalTx)
	prefix := selectedInterface
	if selectedInterface == "all" {
		prefix = "All"
	}
	systray.SetTitle(fmt.Sprintf("%s: ↓%s ↑%s", prefix, downStr, upStr))
}

func onReady() {
	prevStats, _ := parseProcNetDev()

	systray.SetTitle("NetSpeed")
	systray.SetTooltip("Network speed monitor")

	// Create interface selection menu
	interfaceMenu := systray.AddMenuItem("Select Interface", "Choose network interface")
	allItem := interfaceMenu.AddSubMenuItem("All Interfaces", "Monitor all interfaces")
	allItem.Check()
	interfaceItems["all"] = allItem



	// Add individual interfaces
	var interfaces []string
	for iface := range prevStats {
		if iface != "lo" {
			interfaces = append(interfaces, iface)
		}
	}
	sort.Strings(interfaces)

	for _, iface := range interfaces {
		item := interfaceMenu.AddSubMenuItem(iface, fmt.Sprintf("Monitor %s interface", iface))
		interfaceItems[iface] = item
	}

	quitItem := systray.AddMenuItem("Quit", "Exit app")

	// Handle interface selection
	go func() {
		for {
			select {
			case <-allItem.ClickedCh:
				for name, item := range interfaceItems {
					item.Uncheck()
					if name == "all" {
						item.Check()
					}
				}
				selectedInterface = "all"
			}
			for iface, item := range interfaceItems {
				if iface != "all" {
					select {
					case <-item.ClickedCh:
						for name, it := range interfaceItems {
							it.Uncheck()
							if name == iface {
								it.Check()
							}
						}
						selectedInterface = iface
					default:
					}
				}
			}
		}
	}()

	// Update display every second
	go func() {
		for {
			time.Sleep(1 * time.Second)
			curStats, _ := parseProcNetDev()
			updateDisplay(prevStats, curStats)
			prevStats = curStats
		}
	}()

	go func() {
		<-quitItem.ClickedCh
		systray.Quit()
	}()
}

func main() {
	systray.Run(onReady, func() {})
}
