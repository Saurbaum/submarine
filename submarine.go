package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/saurbaum/submarine/position"
	"github.com/saurbaum/submarine/sub"
	"io"
	"math/big"
	"net/http"
)

var seabed []position.Position

var playerSub sub.Sub

func setStartPosition() {
	fmt.Println("Setting starting position")

	playerSub = sub.Create(position.Position{int64(90), int64(10)})

	fmt.Println(playerSub.GetLocation())

	playerSub.SetLocation(position.Position{int64(9), int64(1)})

	fmt.Println(playerSub.GetLocation())
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

		io.WriteString(w, "ping")
	}
}

func location(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("get location")

		pos, err := json.Marshal(playerSub.GetLocation())

		if err == nil {
			io.WriteString(w, string(pos))
		} else {
			io.WriteString(w, "Failed to get location")
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

func main() {
	fmt.Println("Starting Submarine")

	generateBottom(10)
	setStartPosition()

	mux := http.NewServeMux()
	mux.HandleFunc("/location", location)
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/seabed", seabedTest)

	http.ListenAndServe(":80", mux)
}
