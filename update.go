package main

import (
	"fmt"
	"time"
)

var lastUpdate time.Time = time.Now()

func updatePlayers() {
	for {
		updateInterval := time.Since(lastUpdate)

		if updateInterval.Seconds() < 0.033 {
			sleepDuration := time.Duration(33) * time.Millisecond
			time.Sleep(sleepDuration)
		}

		updateInterval = time.Since(lastUpdate)

		for _, value := range players {
			if value.isAlive() {
				value.updateLocation(updateInterval)

				if testCollision(value.location) {
					value.Alive = false
					fmt.Println("Crashed")
				}

				if value.location.X >= int64(maxDistance) {
					value.Alive = false
					fmt.Println("Finished")
				}
			}
		}

		lastUpdate = time.Now()
	}
}

func testCollision(location position) bool {
	band := 0

	for index, value := range seabed {
		if location.X < value.X {
			band = index - 1
			break
		}
	}

	startDepth := seabed[band].Y
	endDepth := seabed[band+1].Y

	if location.Y < startDepth && location.Y < endDepth {
		return false
	}

	baseDifference := seabed[band+1].X - seabed[band].X

	playerSideDifference := location.X - seabed[band].X

	ratio := float32(playerSideDifference) / float32(baseDifference)

	depthDifference := endDepth - startDepth

	playerMaxDepth := startDepth + int64(float32(depthDifference)*ratio)

	return playerMaxDepth < location.Y
}
