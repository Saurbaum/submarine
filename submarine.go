package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
)

const playerIdHeaderKey string = "Playerid"

const maxDepth int64 = 50

var seabed []position

var players = make(map[string]sub)

func getPlayer(r *http.Request) (sub, error) {
	if r.Header[playerIdHeaderKey] == nil || len(r.Header[playerIdHeaderKey]) < 1 {
		return sub{}, errors.New("No playerId")
	}

	var playerId = r.Header[playerIdHeaderKey][0]

	var p, ok = players[playerId]

	if ok {
		return p, nil
	}

	return sub{}, errors.New("Player not found")
}

func generateBottom(length int) {
	fmt.Println("Genreating seabed")
	for i := 0; i < length; i++ {
		yPos, _ := rand.Int(rand.Reader, big.NewInt(maxDepth))
		x := int64(i * 10000)
		y := yPos.Int64()
		seabed = append(seabed, position{x, y})

		fmt.Print(y, " ")
	}

	fmt.Println(" ")
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

func location(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var player, err = getPlayer(r)

		if err == nil {
			pos, err := json.Marshal(player.GetLocation())

			if err == nil {
				io.WriteString(w, string(pos))
				return
			} else {
				io.WriteString(w, "Failed to get location")
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
			startArea := maxDepth - seabed[0].Y
			depthPos, _ := rand.Int(rand.Reader, big.NewInt(startArea))

			// Create player and retun GUID
			players[uuid] = CreateSub(position{int64(0), depthPos.Int64()})
			io.WriteString(w, uuid)
		}
	}
}

func main() {
	fmt.Println("Starting Submarine")

	generateBottom(10)

	mux := http.NewServeMux()
	mux.HandleFunc("/location", location)
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/seabed", seabedTest)
	mux.HandleFunc("/start", createPlayer)

	go updatePlayers()

	http.ListenAndServe(":80", mux)
}
