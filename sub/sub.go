package sub

import (
	"github.com/saurbaum/submarine/position"
)

type Sub struct {
	speed    int64
	heading  int64
	bouency  int64
	location position.Position
}

func Create(startPosition position.Position) Sub {
	return Sub{speed: 0, heading: 0, bouency: 0, location: startPosition}
}

func (s Sub) GetLocation() position.Position {
	return s.location
}

func (s *Sub) SetLocation(newPosition position.Position) {
	s.location = newPosition
}
