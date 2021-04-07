package shot

import (
	"fmt"

	"github.com/hramov/battleship_server/pkg/ship"
)

type Shot struct {
	X int
	Y int
}

func (s *Shot) CheckHit(ships *[]ship.Ship) (int, error) {
	x, y := 0, 0
	for _, sh := range *ships {
		if sh.LivePoints > 0 {
			x = sh.StartX
			y = sh.StartY
			for i := 0; i < sh.Length; i++ {
				if sh.Direction == 0 {
					if s.X == x && s.Y == y+i {
						return sh.ID, fmt.Errorf("Попал!")
					}
				} else if sh.Direction == 1 {
					if s.X == x+i && s.Y == y {
						return sh.ID, fmt.Errorf("Попал!")
					}
				}
			}
		}
	}
	return 0, nil
}
