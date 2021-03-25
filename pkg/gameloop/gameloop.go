package gameloop

func Init() {
	Start()
	// Creating and drawing ships
	// var ships []ship.Ship
	// Player := true

	// PlaceShips(Player, &ships, &b)
	// PlaceShips(!Player, &ships, &b)

	// Game(Player, &b, &ships)
}

func Start() {
	// Creating and drawing battlefield
	// c := battlefield.Client{}
	// c = c.CreateField()
	// b.DrawField()

	// Creating and drawing ships
	// var ships []ship.Ship
	// Player := true

	// PlaceShips(Player, &ships, &b)
	// PlaceShips(!Player, &ships, &b)

	// Game(Player, &b, &ships)
}

// func PlaceShips(Player bool, ships *[]ship.Ship, b *battlefield.BattleField) {
// 	i := 0
// 	fmt.Printf("Расстановка кораблей для игрока %t\n", Player)
// 	for i < 1 {
// 		fmt.Printf("Корабль %d:\n", i+1)
// 		s := ship.Ship{}
// 		s = s.CreateShip(Player)
// 		_, err := b.CheckShip(s)
// 		if err != nil {
// 			fmt.Println(err)
// 		} else {
// 			*ships = append(*ships, s)
// 			*b = b.DrawShip(s)
// 			i++
// 		}
// 	}
// }

// func Game(Player bool, b *battlefield.BattleField, ships *[]ship.Ship) {
// 	for {
// 		ShotX, ShotY := ship.HitShip()
// 		err := b.CheckShot(Player, ShotX, ShotY)
// 		if err != nil {
// 			fmt.Println(err)
// 			continue
// 		}
// 		err1 := b.CheckHit(Player, ShotX, ShotY, ships)
// 		if err1 == false {
// 			fmt.Println(err1)
// 			*b = b.DrawShot(Player, ShotX, ShotY, 0)
// 			Player = !Player
// 		} else {
// 			fmt.Println("Попал!")
// 			*b = b.DrawShot(Player, ShotX, ShotY, 1)

// 			result := IsWon(Player, ships)
// 			fmt.Println(result)
// 			if result {
// 				fmt.Printf("Игрок %t победил!\n", Player)
// 				break
// 			}
// 		}
// 	}
// }

// func IsWon(Player bool, ships *[]ship.Ship) bool {
// 	var newShips, myShips, enemyShips []ship.Ship
// 	newShips = *ships
// 	for i := 0; i < len(newShips); i++ {
// 		if newShips[i].Player == !Player {
// 			enemyShips = append(enemyShips, newShips[i])
// 		} else {
// 			myShips = append(myShips, newShips[i])
// 		}
// 	}
// 	if len(enemyShips) == 0 {
// 		return true
// 	}
// 	return false
// }
