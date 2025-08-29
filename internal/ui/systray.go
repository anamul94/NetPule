package ui

import (
	"fmt"
	"time"

	"github.com/getlantern/systray"

	"netpulse/internal/domain"
)

type SystrayUI struct {
	networkService domain.NetworkService
	speedCalc      domain.SpeedCalculator
	selectedIface  string
	interfaceItems map[string]*systray.MenuItem
	prevStats      map[string]domain.NetStats
}

func NewSystrayUI(networkService domain.NetworkService, speedCalc domain.SpeedCalculator) *SystrayUI {
	return &SystrayUI{
		networkService: networkService,
		speedCalc:      speedCalc,
		selectedIface:  "all",
		interfaceItems: make(map[string]*systray.MenuItem),
	}
}

func (ui *SystrayUI) Run() {
	systray.Run(ui.onReady, func() {})
}

func (ui *SystrayUI) onReady() {
	ui.prevStats, _ = ui.networkService.GetInterfaceStats()

	systray.SetTitle("NetSpeed")
	systray.SetTooltip("Network speed monitor")

	ui.setupMenu()
	ui.startMonitoring()
}

func (ui *SystrayUI) setupMenu() {
	interfaceMenu := systray.AddMenuItem("Select Interface", "Choose network interface")
	allItem := interfaceMenu.AddSubMenuItem("All Interfaces", "Monitor all interfaces")
	allItem.Check()
	ui.interfaceItems["all"] = allItem

	interfaces := ui.networkService.GetInterfaceNames()
	for _, iface := range interfaces {
		item := interfaceMenu.AddSubMenuItem(iface, fmt.Sprintf("Monitor %s interface", iface))
		ui.interfaceItems[iface] = item
	}

	quitItem := systray.AddMenuItem("Quit", "Exit app")

	go ui.handleMenuClicks(allItem, quitItem)
}

func (ui *SystrayUI) handleMenuClicks(allItem, quitItem *systray.MenuItem) {
	for {
		select {
		case <-allItem.ClickedCh:
			ui.selectInterface("all")
		case <-quitItem.ClickedCh:
			systray.Quit()
			return
		}
		for iface, item := range ui.interfaceItems {
			if iface != "all" {
				select {
				case <-item.ClickedCh:
					ui.selectInterface(iface)
				default:
				}
			}
		}
	}
}

func (ui *SystrayUI) selectInterface(iface string) {
	for name, item := range ui.interfaceItems {
		item.Uncheck()
		if name == iface {
			item.Check()
		}
	}
	ui.selectedIface = iface
}

func (ui *SystrayUI) startMonitoring() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			ui.updateDisplay()
		}
	}()
}

func (ui *SystrayUI) updateDisplay() {
	curStats, _ := ui.networkService.GetInterfaceStats()
	var totalRx, totalTx uint64

	if ui.selectedIface == "all" {
		for iface, cur := range curStats {
			if iface == "lo" {
				continue
			}
			if prev, ok := ui.prevStats[iface]; ok {
				rx, tx := ui.speedCalc.CalculateSpeed(prev, cur)
				totalRx += rx
				totalTx += tx
			}
		}
	} else {
		if cur, ok := curStats[ui.selectedIface]; ok {
			if prev, ok := ui.prevStats[ui.selectedIface]; ok {
				totalRx, totalTx = ui.speedCalc.CalculateSpeed(prev, cur)
			}
		}
	}

	downStr := ui.speedCalc.FormatSpeed(totalRx)
	upStr := ui.speedCalc.FormatSpeed(totalTx)
	prefix := ui.selectedIface
	if ui.selectedIface == "all" {
		prefix = "All"
	}
	systray.SetTitle(fmt.Sprintf("%s: ↓%s ↑%s", prefix, downStr, upStr))

	ui.prevStats = curStats
}