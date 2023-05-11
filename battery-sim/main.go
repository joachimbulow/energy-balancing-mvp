package main

import (
	"time"

	"github.com/joachimbulow/pem-energy-balance/src"
	"github.com/joachimbulow/pem-energy-balance/src/util"
)

var (
	nBatteries = util.GetNumberOfBatteries()
)

func main() {
	initializeBatteries()
}

func initializeBatteries() {
	for i := 0; i < nBatteries; i++ {
		// To spread out the spawn of the batteries so they spam less
		time.Sleep(100 * time.Millisecond)
		go startBattery()
	}

	// To keep routines running' we start sleepin'
	for {
		time.Sleep(1 * time.Second)
	}
}

func startBattery() {
	src.NewBattery()
}
