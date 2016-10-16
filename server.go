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
	mux.HandleFunc("/start", createPlayer)
	mux.HandleFunc("/buoyancy", buoyancy)

	http.ListenAndServe(":8080", mux)
}

func getPlayer(r *http.Request) (*sub, error) {
	if r.Header[playerIdHeaderKey] == nil || len(r.Header[playerIdHeaderKey]) < 1 {
		return nil, errors.New("No playerId")
	}

	var playerId = r.Header[playerIdHeaderKey][0]

	var p, ok = players[playerId]

	if ok {
		return p, nil
	}

	return nil, errors.New("Player not found")
}

func ping(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var _, err = getPlayer(r)

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
		var p, err = getPlayer(r)
		fmt.Println("buoyancy for player: ", p)

		if err == nil {
			fmt.Println("r.Header: ", r.Header)
			var buoyancyString = r.Header[buoyancyHeaderKey][0]
			fmt.Println("buoyancyString: ", buoyancyString)
			buoyancy, err := strconv.ParseFloat(buoyancyString, 64)
			if err == nil {
				p.SetBuoyancy(buoyancy)
				io.WriteString(w, "buoyancy")
			} else {
				io.WriteString(w, err.Error())
			}
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

func createPlayer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("createPlayer")

		var uuid, err = newUUID()

		fmt.Println(uuid)

		if err != nil {
			io.WriteString(w, err.Error())
		} else {
			// Pick stating point somewhere above the bottom.
			depthPos, _ := rand.Int(rand.Reader, big.NewInt(seabed[0].Y))

			// Create player and retun GUID
			players[uuid] = CreateSub(position{int64(0), depthPos.Int64()})
			io.WriteString(w, uuid)
		}
	}
}
