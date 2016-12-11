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

const playerIdHeaderKey string = "Playerid"
const buoyancyHeaderKey string = "Buoyancy"
const speedHeaderKey string = "Speed"

func generateBottom(length int) {
	fmt.Println("Genreating seabed")
	for i := 0; i < length; i++ {
		yPos, _ := rand.Int(rand.Reader, big.NewInt(maxDepth))
		x := int64(i * seabedStepWidth)
		y := yPos.Int64()
		seabed = append(seabed, position{x, y})

		fmt.Print("x: ", x, " y:", y, " ")
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
	if r.Header[playerIdHeaderKey] == nil || len(r.Header[playerIdHeaderKey]) < 1 {
		return nil, "", errors.New("No playerId")
	}

	var playerId = r.Header[playerIdHeaderKey][0]

	var p, ok = players[playerId]

	if ok {
		return p, playerId, nil
	}

	return nil, playerId, errors.New("Player not found")
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
			if err == nil {
				p.setSpeed(speedString)
				io.WriteString(w, "speed")
			} else {
				io.WriteString(w, err.Error())
			}
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
			io.WriteString(w, string(p.GetStatus()))
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

		io.WriteString(w, string(bed))
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
