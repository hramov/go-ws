package shot

import "github.com/hramov/battleship_server/pkg/ship"

type Shot struct {
	X int
	Y int
}

func (s *Shot) CheckHit(ships *map[int]ship.Ship) error {
	for _, sh := range *ships {
		if sh.LivePoints > 0 {
			//
		}
	}
	return nil
}
