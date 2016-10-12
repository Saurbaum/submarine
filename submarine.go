package main

import (
	"crypto/rand"
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"math/big"
	"runtime"
)

const playerIdHeaderKey string = "Playerid"

const maxDepth int64 = 50000

var seabed []position

var players = make(map[string]sub)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	fmt.Println("Starting Submarine")

	generateBottom(10)

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
	window, err := glfw.CreateWindow(1024, 768, "Submarine", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	setupScene()

	for !window.ShouldClose() {
		// Do OpenGL stuff.
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.LineWidth(2.5)
		gl.Color3f(1.0, 1.0, 0.0)
		gl.Begin(gl.LINES)
		gl.Vertex3f(-1, 0, 0)
		gl.Vertex3f(1, 0, 0)
		gl.End()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
