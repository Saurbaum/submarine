package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
)

const playerIDHeaderKey string = "Playerid"
const buoyancyHeaderKey string = "Buoyancy"
const speedHeaderKey string = "Speed"
const depthRandom int = 3

func getRandomDepth() int64 {
	yPos, _ := rand.Int(rand.Reader, big.NewInt(maxDepth))

	return yPos.Int64()
}

func generateBottom(length int) {
	fmt.Println("Generating seabed")

	startPosition := getRandomDepth()

	for i := 0; i < length; i++ {
		yPos := int64(0)
		if i == 0 {
			yPos = startPosition
		} else {
			for iDepth := 0; iDepth < depthRandom; iDepth++ {
				yTemp := getRandomDepth() / int64(depthRandom*10)

				yPos += yTemp
			}

			postitive, _ := rand.Int(rand.Reader, big.NewInt(2))

			fmt.Println("postitive: ", postitive)

			if postitive.Int64() == int64(1) {
				yPos = yPos * int64(-1)
				fmt.Println("Invert yPos ", yPos)
			}

			proposedY := startPosition + yPos

			if proposedY <= maxDepth/10 || proposedY >= maxDepth {
				yPos = yPos * int64(-1)
			}

			yPos = startPosition + yPos
		}

		x := int64(i * seabedStepWidth)
		y := yPos
		seabed = append(seabed, position{x, y})

		fmt.Print("x: ", x, " y:", y, " ")

		startPosition = y
	}

	fmt.Println(" ")
}

func startServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/seabed", seabedTest)
	mux.HandleFunc("/connect", connect)
	mux.HandleFunc("/start", spawnPlayer)
	mux.HandleFunc("/buoyancy", buoyancy)
	mux.HandleFunc("/speed", speed)
	mux.HandleFunc("/status", status)

	http.ListenAndServe(":8080", mux)
}

func getPlayer(r *http.Request) (*sub, string, error) {
	if r.Header[playerIDHeaderKey] == nil || len(r.Header[playerIDHeaderKey]) < 1 {
		return nil, "", errors.New("no playerid")
	}

	var playerID = r.Header[playerIDHeaderKey][0]

	var p, ok = players[playerID]

	if ok {
		return p, playerID, nil
	}

	return nil, playerID, errors.New("player not found")
}

func ping(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var _, _, err = getPlayer(r)

		if err == nil {
			io.WriteString(w, "ping")
		} else {
			io.WriteString(w, err.Error())
		}
	}
}

func buoyancy(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var p, _, err = getPlayer(r)

		if err == nil {
			var buoyancyString = r.Header[buoyancyHeaderKey][0]
			buoyancy, err := strconv.ParseFloat(buoyancyString, 64)
			if err == nil {
				p.setBuoyancy(buoyancy)
				io.WriteString(w, "buoyancy")
			} else {
				io.WriteString(w, err.Error())
			}
		} else {
			io.WriteString(w, err.Error())
		}
	}
}

func speed(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var p, _, err = getPlayer(r)

		if err == nil {
			var speedString = r.Header[speedHeaderKey][0]
			p.setSpeed(speedString)
			io.WriteString(w, "speed")

		} else {
			io.WriteString(w, err.Error())
		}
	}
}

func status(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var p, _, err = getPlayer(r)

		if err == nil {
			io.Writer.Write(w, p.GetStatus())
		} else {
			io.WriteString(w, err.Error())
		}
	}
}

func seabedTest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("get seabedTest")

		bed, _ := json.Marshal(seabed)

		io.Writer.Write(w, bed)
	}
}

func connect(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("connect")

		var uuid, err = newUUID()

		fmt.Println(uuid)

		if err != nil {
			io.WriteString(w, err.Error())
		} else {
			// Return GUID
			io.WriteString(w, uuid)
		}
	}
}

func spawnPlayer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("spawnPlayer")

		var _, uuid, err = getPlayer(r)

		if err != nil {
			// Pick stating point somewhere above the bottom.
			depthPos, _ := rand.Int(rand.Reader, big.NewInt(seabed[0].Y))

			// Create player and retun GUID
			players[uuid] = CreateSub(position{int64(0), depthPos.Int64()})
			io.WriteString(w, uuid)
		}
	}
}
