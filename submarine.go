package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/saurbaum/submarine/position"
	"github.com/saurbaum/submarine/sub"
	"io"
	"math/big"
	"net/http"
)

const playerIdHeaderKey string = "Playerid"

var seabed []position.Position

var players = make(map[string]sub.Sub)

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func getPlayer(r *http.Request) (sub.Sub, error) {
	if r.Header[playerIdHeaderKey] == nil || len(r.Header[playerIdHeaderKey]) < 1 {
		return sub.Sub{}, errors.New("No playerId")
	}

	var playerId = r.Header[playerIdHeaderKey][0]

	var p, ok = players[playerId]

	if ok {
		return p, nil
	}

	return sub.Sub{}, errors.New("Player not found")
}

func generateBottom(length int) {
	fmt.Println("Genreating seabed")
	for i := 0; i < length; i++ {
		yPos, _ := rand.Int(rand.Reader, big.NewInt(50))
		x := int64(i * 10)
		y := yPos.Int64()
		seabed = append(seabed, position.Position{x, y})
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("ping")

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
		fmt.Println("get location")

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
		fmt.Println("get createPlayer")

		var uuid, err = newUUID()

		fmt.Println(uuid)

		if err != nil {
			io.WriteString(w, err.Error())
		} else {
			// Create player and retun GUID
			players[uuid] = sub.Create(position.Position{int64(90), int64(10)})
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

	http.ListenAndServe(":80", mux)
}
