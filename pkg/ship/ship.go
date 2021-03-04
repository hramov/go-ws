package ship

import (
	"fmt"
)

type Ship struct {
	Player     bool
	Length     int
	StartX     int
	StartY     int
	Direction  int
	Hit        bool
	LivePoints int
	Live       bool
}

func (s Ship) CreateShip(Player bool) Ship {
	s.Player = Player
	fmt.Printf("%s ", "Введите числовую позицию (1-10):")
	fmt.Scanf("%d", &s.StartY)
	fmt.Printf("%s ", "Введите буквенную позицию (1-10):")
	fmt.Scanf("%d", &s.StartX)
	fmt.Printf("%s ", "Введите направление (0: -, 1: |):")
	fmt.Scanf("%d", &s.Direction)
	fmt.Printf("%s ", "Введите длину корабля (1,2,3,4):")
	fmt.Scanf("%d", &s.Length)
	s.LivePoints = s.Length
	fmt.Println(s)
	return s
}

func HitShip() (int, int) {
	var ShotX, ShotY int
	fmt.Println("Введите координаты выстрела! Число:")
	fmt.Scanf("%d", &ShotX)
	fmt.Println("Введите координаты выстрела! Буква:")
	fmt.Scanf("%d", &ShotY)
	return ShotX, ShotY
}
