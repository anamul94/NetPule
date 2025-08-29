package main

import (
	"netpulse/internal/infrastructure"
	"netpulse/internal/ui"
)

func main() {
	networkService := infrastructure.NewLinuxNetworkService()
	speedCalc := infrastructure.NewSpeedService()
	systrayUI := ui.NewSystrayUI(networkService, speedCalc)

	systrayUI.Run()
}