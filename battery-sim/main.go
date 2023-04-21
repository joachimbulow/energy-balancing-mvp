package main

import (
	"time"

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

	for {
		time.Sleep(1 * time.Second)
	}
}

func startBattery() {
	src.NewBattery()
}
