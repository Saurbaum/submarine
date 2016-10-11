package main

import (
	"time"
)

type sub struct {
	speed    int64
	heading  float64
	buoyancy int64
	location position
	updated  time.Time
	alive    bool
}

func CreateSub(startPosition position) sub {
	return sub{speed: 1, heading: 0, buoyancy: 0, location: startPosition, updated: time.Now(), alive: true}
}

func (s sub) GetLocation() position {
	return s.location
}

func (s *sub) updateLocation(updateInterval time.Duration) {
	distance := (s.speed * int64(updateInterval/time.Millisecond))

	s.location.X = s.location.X + distance
}

func (s sub) GetSpeed() int64 {
	return s.speed
}

func (s *sub) SetSpeed(newSpeed int64) {
	s.speed = newSpeed
}

func (s sub) GetBuoyancy() int64 {
	return s.buoyancy
}

func (s *sub) SetBuoyancy(newBuoyancy int64) {
	s.buoyancy = newBuoyancy
}

func (s sub) GetHeading() float64 {
	return s.heading
}

func (s *sub) SetHeading(newHeading float64) {
	s.heading = newHeading
}
