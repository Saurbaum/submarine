package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
)

type position struct {
	X int64
	Y int64
}

var seabed []position

var submarinePosition position

func setStartPosition() {
	fmt.Println("Setting starting position")
	submarinePosition.X = 0
	submarinePosition.Y = 0
	fmt.Println(submarinePosition)
}

func generateBottom(length int) {
	fmt.Println("Genreating seabed")
	for i := 0; i < length; i++ {
		yPos, _ := rand.Int(rand.Reader, big.NewInt(50))
		x := int64(i * 10)
		y := yPos.Int64()
		seabed = append(seabed, position{x, y})
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("ping")

		io.WriteString(w, "get location")
	}
}

func location(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Println("get location")

		pos, err := json.Marshal(submarinePosition)

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
