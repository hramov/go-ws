package battlefield

import (
	"fmt"

	"github.com/hramov/battleship_server/pkg/ship"
)

const FIELD_WIDTH = 12
const FIELD_HEIGHT = 12
const LETTER_STRING = "   А Б В Г Д Е Ж З И К\t\t   А Б В Г Д Е Ж З И К\n"

type Field [FIELD_WIDTH][FIELD_HEIGHT]string

type BattleField struct {
	Field     Field
	ShotField Field
}

func (b *BattleField) CreateField() {
	for i := 0; i < FIELD_HEIGHT; i++ {
		for j := 0; j < FIELD_WIDTH; j++ {
			if i == 0 || i == FIELD_HEIGHT-1 {
				b.Field[i][j] = "*"
				b.ShotField[i][j] = "*"
				continue
			}
			if j == 0 || j == FIELD_WIDTH-1 {
				b.Field[i][j] = "*"
				b.ShotField[i][j] = "*"
			} else {
				b.Field[i][j] = "_"
				b.ShotField[i][j] = "_"
			}
		}
	}
}

func (b *BattleField) CheckShip(s ship.Ship) error {
	errorMessage := ""

	if b.Field[s.StartY][s.StartX] == "_" { //В начальной точке нет другого корабля
		if s.Direction == 0 {
			if s.StartY+s.Length < FIELD_HEIGHT { //Проверка выхода за границы поля
				if b.Field[s.StartY+s.Length][s.StartX] != "_" { //Проверка доступности клетки в конечной точке
					errorMessage = "Уперся в *"
				} else {
					return nil
				}
			} else {
				errorMessage = "Вышел за границу"
			}
		} else {
			if s.StartX+s.Length < FIELD_WIDTH {
				if b.Field[s.StartY][s.StartX+s.Length] != "_" {
					errorMessage = "Уперся в *"
				} else {
					return nil
				}
			} else {
				errorMessage = "Вышел за границу"
			}
		}
	} else {
		errorMessage = "Первое условие"
	}
	return fmt.Errorf("%s\n", errorMessage)
}

func (b *BattleField) CreateShip(s ship.Ship) error {
	for i := 0; i < s.Length; i++ {
		if s.Direction == 0 {
			b.Field[s.StartY][s.StartX+i] = "O"
			b.Field[s.StartY+1][s.StartX+i] = "*"
			b.Field[s.StartY-1][s.StartX+i] = "*"
			b.Field[s.StartY][s.StartX-1] = "*"
			b.Field[s.StartY][s.StartX+s.Length] = "*"
			b.Field[s.StartY+1][s.StartX+s.Length] = "*"
			b.Field[s.StartY-1][s.StartX+s.Length] = "*"
			b.Field[s.StartY+1][s.StartX-1] = "*"
			b.Field[s.StartY-1][s.StartX-1] = "*"
		} else if s.Direction == 1 {
			b.Field[s.StartY+i][s.StartX] = "O"
			b.Field[s.StartY+i][s.StartX+1] = "*"
			b.Field[s.StartY+i][s.StartX-1] = "*"
			b.Field[s.StartY-1][s.StartX] = "*"
			b.Field[s.StartY+s.Length][s.StartX] = "*"
			b.Field[s.StartY+s.Length][s.StartX+1] = "*"
			b.Field[s.StartY+s.Length][s.StartX-1] = "*"
			b.Field[s.StartY-1][s.StartX+1] = "*"
			b.Field[s.StartY-1][s.StartX-1] = "*"
		}
	}
	return nil
}

// func (b BattleField) DrawShot(Player bool, ShotX, ShotY int, Result int) BattleField {
// 	if Player {
// 		if Result == 0 {
// 			b.EnemyField[ShotX][ShotY] = "*"
// 		} else {
// 			b.EnemyField[ShotX][ShotY] = "X"
// 		}
// 	} else {
// 		if Result == 0 {
// 			b.MyField[ShotX][ShotY] = "*"
// 		} else {
// 			b.MyField[ShotX][ShotY] = "X"
// 		}
// 	}
// 	b.DrawField()
// 	return b
// }

// func (b BattleField) ClearField() {
// 	b.MyField = Field{}
// 	b.EnemyField = Field{}
// 	b.DrawField()
// }

// func (b BattleField) CheckShot(Player bool, shot ship.Shot) error {

// 	if Player {
// 		if b.EnemyField[shot.X][shot.Y] == "O" || b.EnemyField[shot.X][shot.Y] == "*" {
// 			return nil
// 		}
// 	} else {
// 		if b.MyField[shot.X][shot.Y] == "O" || b.MyField[shot.X][shot.Y] == "*" {
// 			return nil
// 		}
// 	}
// 	return fmt.Errorf("%s", "Сюда нельзя стрелять")
// }

// func (b BattleField) CheckHit(Player bool, shot ship.Shot, ships *[]ship.Ship) bool {

// 	var result bool = false

// 	var newShips, myShips, enemyShips []ship.Ship

// 	newShips = *ships

// 	for i := 0; i < len(newShips); i++ {
// 		if newShips[i].Player == !Player {
// 			enemyShips = append(enemyShips, newShips[i])
// 		} else {
// 			myShips = append(myShips, newShips[i])
// 		}
// 	}

// 	for i := 0; i < len(enemyShips); i++ {

// 		for j := 0; j < enemyShips[i].Length; j++ {
// 			if enemyShips[i].Direction == 0 {
// 				if shot.X == enemyShips[i].StartX && shot.Y == enemyShips[i].StartY+j {
// 					enemyShips[i].LivePoints--
// 					result = true
// 				}
// 			} else {
// 				if shot.X == enemyShips[i].StartX+j && shot.Y == enemyShips[i].StartY {
// 					enemyShips[i].LivePoints--
// 					result = true
// 				}
// 			}
// 		}

// 		if enemyShips[i].LivePoints > 0 {
// 			myShips = append(myShips, enemyShips[i])
// 			fmt.Println(myShips)
// 		}

// 	}

// 	*ships = myShips
// 	return result
// }

// func whosTurn(s ship.Ship) bool {

// }
