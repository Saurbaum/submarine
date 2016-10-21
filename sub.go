package main

import (
	"encoding/json"
	"time"
)

type sub struct {
	Speed    float64
	Heading  float64
	Buoyancy float64
	location position
	updated  time.Time
	Alive    bool
}

func CreateSub(startPosition position) *sub {
	return &sub{Speed: 1, Heading: 0, Buoyancy: 0, location: startPosition, updated: time.Now(), Alive: true}
}

func (s sub) GetLocation() position {
	return s.location
}

func (s *sub) updateLocation(updateInterval time.Duration) {
	distance := (s.Speed * float64(updateInterval/time.Millisecond))
	depth := (s.Buoyancy * float64(updateInterval/time.Millisecond))

	s.location.X = s.location.X + int64(distance)
	s.location.Y = s.location.Y + int64(depth)

	if s.location.Y < 0 {
		s.location.Y = 0
	}
}

func (s sub) getSpeed() float64 {
	return s.Speed
}

func (s *sub) setSpeed(newSpeed string) {
	speedValue := s.getSpeed()

	switch newSpeed {
	case "Full Ahead":
		speedValue = 2.0
	case "Half Ahead":
		speedValue = 1.0
	case "Slow Ahead":
		speedValue = 0.5
	case "Dead Slow Ahead":
		speedValue = 0.2
	case "Stop":
		speedValue = 0
	case "Dead Slow Astern":
		speedValue = -0.1
	case "Slow Astern":
		speedValue = -0.25
	case "Half Astern":
		speedValue = -0.5
	case "Full Astern":
		speedValue = -1.0
	}

	s.Speed = speedValue
}

func (s sub) getBuoyancy() float64 {
	return s.Buoyancy
}

func (s *sub) setBuoyancy(newBuoyancy float64) {
	s.Buoyancy -= newBuoyancy
}

func (s sub) getHeading() float64 {
	return s.Heading
}

func (s *sub) setHeading(newHeading float64) {
	s.Heading = newHeading
}

func (s sub) isAlive() bool {
	return s.Alive
}

func (s *sub) GetStatus() []byte {
	data, err := json.Marshal(s)
	if err == nil {
		return data
	}

	return nil
}
