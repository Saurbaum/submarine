package main

type sub struct {
	speed    int64
	heading  int64
	bouency  int64
	location position
}

func CreateSub(startPosition position) sub {
	return sub{speed: 0, heading: 0, bouency: 0, location: startPosition}
}

func (s sub) GetLocation() position {
	return s.location
}

func (s *sub) SetLocation(newPosition position) {
	s.location = newPosition
}
