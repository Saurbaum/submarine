package main

import (
	"crypto/rand"
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"math/big"
	"runtime"
	"time"
)

const playerIdHeaderKey string = "Playerid"

const maxDepth int64 = 50000

const screenWidth int = 2000

const screenHeight int = 1000

const seabedSegments int = 10

var seabed []position

var players = make(map[string]sub)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	fmt.Println("Starting Submarine")

	generateBottom(seabedSegments)

	go updatePlayers()

	go startServer()

	render()
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

func setupScene() {
	gl.ClearColor(0, 0, 0, 1)
}

func render() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(screenWidth, screenHeight, "Submarine", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	setupScene()

	seabedStep := 2 / float32(seabedSegments-1)
	seabedRatio := 2 / float32(maxDepth)
	lastUpdate := time.Now()

	for !window.ShouldClose() {
		updateInterval := time.Since(lastUpdate)

		if updateInterval.Seconds() < 0.033 {
			sleepDuration := time.Duration(33) * time.Millisecond
			time.Sleep(sleepDuration)
		}

		// Do OpenGL stuff.
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.LineWidth(2.5)
		gl.Color3f(1.0, 1.0, 0.0)
		gl.Begin(gl.LINES)

		lastX := float32(0)
		lastY := float32(0)

		for index, value := range seabed {
			if index > 1 {
				gl.Vertex3f(lastX, lastY*-1, 0)
			}
			x := (float32(index) * seabedStep) - 1
			y := (float32(value.Y) * seabedRatio) - 1

			gl.Vertex3f(x, y*-1, 0)
			lastX = x
			lastY = y
		}

		gl.End()

		window.SwapBuffers()
		glfw.PollEvents()

		lastUpdate = time.Now()
	}
}
