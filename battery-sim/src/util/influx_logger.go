package util

import (
	"log"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

const (
	influxDB = "influx"
	username = "admin"
	password = "admin"
)

// Map from string to timestamp in millis
var latencyMap = make(map[string]int64)

func LogLatencies(latencyChannel chan string) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     GetInfluxDB(),
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database: influxDB,
	})

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if len(bp.Points()) > 0 {
				err := c.Write(bp)
				if err != nil {
					log.Println("Failed to write batch points:", err)
				}

				// Create a new batch point after the previous one is sent
				bp, _ = client.NewBatchPoints(client.BatchPointsConfig{
					Database: influxDB,
				})
			}
		}
	}()

	for requestId := range latencyChannel {
		departureTime, ok := latencyMap[requestId]
		if !ok {
			// store the departure time
			latencyMap[requestId] = time.Now().UnixNano() / int64(time.Millisecond) // To get the time in millis
			continue
		}
		// if the requestId is already present, calculate the latency and delete the requestId
		arrivalTime := time.Now().UnixNano() / int64(time.Millisecond)
		latency := arrivalTime - departureTime

		tags := map[string]string{"type": "latency"}
		fields := map[string]interface{}{
			"latency": latency,
		}
		pt, err := client.NewPoint("latency", tags, fields, time.Now())
		if err != nil {
			log.Fatal(err)
		}

		// Add the point to the current batch point
		bp.AddPoint(pt)
		delete(latencyMap, requestId) // remove the requestId from the map
	}
}
