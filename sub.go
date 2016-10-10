package main

type sub struct {
	speed    int64
	heading  int64
	buoyancy int64
	location position
}

func CreateSub(startPosition position) sub {
	return sub{speed: 0, heading: 0, buoyancy: 0, location: startPosition}
}

func (s sub) GetLocation() position {
	return s.location
}

func (s *sub) SetLocation(newPosition position) {
	s.location = newPosition
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

func (s sub) GetHeading() int64 {
	return s.heading
}

func (s *sub) SetHeading(newHeading int64) {
	s.heading = heading
}
