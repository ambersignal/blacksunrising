package state

import "github.com/ambersignal/blacksunrising/pkg/geom"

func (s *State) MoveCameraTo(pos geom.Vec2) {
	s.Camera = geom.RectangleBySize(pos, s.Camera.Size())

	s.normalizeCamera()
}

func (s *State) MoveCameraOn(shift geom.Vec2) {
	newPos := s.Camera.Center().Add(shift)

	s.MoveCameraTo(newPos)
}

// normalizeCamera ensures that camera stays within world bounds.
func (s *State) normalizeCamera() {
	semiSize := s.Camera.Size().Div(2)
	newPos := s.Camera.Center()

	for i := range 2 {
		if s.Camera.Min[i] < 0 {
			newPos[i] = semiSize[i]
		}

		if s.Camera.Max[i] > s.WorldSize[i] {
			newPos[i] = s.WorldSize[i] - semiSize[i]
		}
	}

	s.Camera = geom.RectangleBySize(newPos, s.Camera.Size())
}
