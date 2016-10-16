package main

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"runtime"
	"time"
)

const maxDepth int64 = 50000

const screenWidth int = 1500

const screenHeight int = 1000

const seabedStepWidth int = 10000

const seabedSegments int = 10

const seabedStep float32 = 2 / float32(seabedSegments-1)

const seabedDepthRatio float32 = 2 / float32(maxDepth)

const playerVisualHeight float32 = 0.01
const playerVisualWidth float32 = 0.01 * float32(float32(screenHeight)/float32(screenWidth))

const maxDistance int = seabedStepWidth * (seabedSegments - 1)

const seabedWitdthRatio float32 = 2 / float32(maxDistance)

var seabed []position

var players = make(map[string]*sub)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	fmt.Println("Starting Submarine Server")

	generateBottom(seabedSegments)

	go updatePlayers()

	go startServer()

	render()
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

	lastUpdate := time.Now()

	for !window.ShouldClose() {
		updateInterval := time.Since(lastUpdate)

		if updateInterval.Seconds() < 0.033 {
			sleepDuration := time.Duration(33) * time.Millisecond
			time.Sleep(sleepDuration)
		}

		// Do OpenGL stuff.
		gl.Clear(gl.COLOR_BUFFER_BIT)

		drawSeabed()
		drawPlayers()

		window.SwapBuffers()
		glfw.PollEvents()

		lastUpdate = time.Now()
	}
}

func drawSeabed() {
	gl.LineWidth(2.5)
	gl.Color3f(1.0, 1.0, 0.0)
	gl.Begin(gl.LINES)

	lastX := float32(0)
	lastY := float32(0)

	for index, value := range seabed {
		if index > 1 {
			gl.Vertex3f(lastX, lastY*-1, 0)
		}
		x := (float32(value.X) * seabedWitdthRatio) - 1
		y := (float32(value.Y) * seabedDepthRatio) - 1

		gl.Vertex3f(x, y*-1, 0)
		lastX = x
		lastY = y
	}

	gl.End()
}

func drawPlayers() {
	for _, player := range players {
		if player.IsAlive() {
			gl.Color3f(1.0, 0.0, 1.0)
			gl.Begin(gl.QUADS)

			screenLocation := player.GetLocation()

			X := (float32(screenLocation.X) * float32(seabedWitdthRatio)) - 1
			Y := ((float32(screenLocation.Y) * seabedDepthRatio) - 1) * -1

			gl.Vertex3f(float32(X)-playerVisualWidth, float32(Y)-playerVisualHeight, 0)
			gl.Vertex3f(float32(X)+playerVisualWidth, float32(Y)-playerVisualHeight, 0)
			gl.Vertex3f(float32(X)+playerVisualWidth, float32(Y)+playerVisualHeight, 0)
			gl.Vertex3f(float32(X)-playerVisualWidth, float32(Y)+playerVisualHeight, 0)

			gl.End()
		}
	}
}
