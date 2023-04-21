package main

import (
	"fmt"

	"github.com/joachimbulow/pem-energy-balance/src"
)

var (
	nBatteries = 1
)

func main() {
	initialize()
}

func initialize() {
	for i := 0; i < nBatteries; i++ {
		go startBattery()
	}
}

func startBattery() {
	battery := src.NewBattery()
	fmt.Print("New battery created: ", battery.ID)
}
