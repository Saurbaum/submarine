package main

import (
	"fmt"
	"time"
)

var lastUpdate time.Time = time.Now()

func updatePlayers() {
	for {
		updatedPlayers := make(map[string]sub)

		updateInterval := time.Since(lastUpdate)

		if updateInterval.Seconds() < 0.033 {
			sleepDuration := time.Duration(33) * time.Millisecond
			time.Sleep(sleepDuration)
		}

		updateInterval = time.Since(lastUpdate)

		for key, value := range players {
			if value.alive {
				value.updateLocation(updateInterval)

				if testCollision(value.location) {
					//value.alive = false
					//fmt.Println("Crashed")
				}

				fmt.Println("player: ", value.location.X, " Max", maxDistance)
				if value.location.X >= int64(maxDistance) {
					value.alive = false
					fmt.Println("Finished")
				}

				updatedPlayers[key] = value
			}
		}

		players = updatedPlayers

		lastUpdate = time.Now()
	}
}

func testCollision(location position) bool {
	band := 0

	for index, value := range seabed {
		if location.X > value.X {
			band = index
			break
		}
	}

	if band == len(seabed) {
		return true
	}

	sideDifference := seabed[band+1].X - seabed[band].X

	playerSideDifference := location.X - seabed[band].X

	ratio := sideDifference / playerSideDifference

	depthDifference := seabed[band].Y - seabed[band+1].Y

	playerMaxDepth := depthDifference * ratio

	if playerMaxDepth < location.Y {
		return true
	}

	return false
}
