package main

import (
	"github.com/go-gl/gl/v2.1/gl"
	"fmt"
)

type vertex struct{
	x float32
	y float32
	z float32
}

type floorTile struct{
	vertices [4]vertex
}

type seabed struct{
	floorTiles []floorTile
}

func CreateSeabed(length int, step float32) seabed {
	fmt.Println("Genreating seabed")

	var bottom seabed
	var tile floorTile

	tile.vertices[0].x = 0
	tile.vertices[0].y = 0
	tile.vertices[0].z = 0

	tile.vertices[1].x = 50
	tile.vertices[1].y = 50
	tile.vertices[1].z = 0

	tile.vertices[2].x = 50
	tile.vertices[2].y = 0
	tile.vertices[2].z = 50

	tile.vertices[3].x = 0
	tile.vertices[3].y = 0
	tile.vertices[3].z = 50

	bottom.floorTiles = append(bottom.floorTiles, tile)

	return bottom
}

func printTile(tile floorTile){
	fmt.Println( tile.vertices )
}

func (s seabed) Render(){
	gl.LineWidth(2.5)
	gl.Color3f(1.0, 1.0, 0.0)
	gl.Begin(gl.QUADS)

	for _, value := range s.floorTiles {
		for _, vert := range value.vertices {
			x := (float32(vert.x) * seabedWitdthRatio) - 1 
		    y := (float32(vert.y) * seabedDepthRatio) - 1 
    		z := (float32(vert.z) * seabedWitdthRatio) - 1 
			gl.Vertex3f(x, y, z)
		}
	}

	gl.End()
}
