package battlefield

import (
	"fmt"

	"github.com/hramov/battleship_server/pkg/ship"
	"github.com/hramov/battleship_server/pkg/shot"
)

const (
	FIELD_WIDTH   = 12
	FIELD_HEIGHT  = 12
	LETTER_STRING = "   А Б В Г Д Е Ж З И К\t\t   А Б В Г Д Е Ж З И К\n"

	XDirection = 0
	YDirection = 1
)

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

func (b *BattleField) CheckShip(s ship.Ship, ships *[]ship.Ship) error {
	errorMessage := ""

	//Проверка на количество кораблей определенной длины
	if err := CheckQuantity(s, ships); err != nil {
		return fmt.Errorf("%s\n", err)
	}

	//Проверка на корректность расположения
	if b.Field[s.StartY][s.StartX] == "_" { //В начальной точке нет другого корабля
		if s.Direction == YDirection {
			if s.StartY+s.Length-1 < FIELD_HEIGHT { //Проверка выхода за границы поля
				if b.Field[s.StartY+s.Length-1][s.StartX] != "_" { //Проверка доступности клетки в конечной точке
					errorMessage = "Неверное расположение корабля"
				} else {
					return nil
				}
			} else {
				errorMessage = "Вышел за границу поля"
			}
		} else if s.Direction == XDirection {
			if s.StartX+s.Length-1 < FIELD_WIDTH {
				if b.Field[s.StartY][s.StartX+s.Length-1] != "_" {
					errorMessage = "Неверное расположение корабля"
				} else {
					return nil
				}
			} else {
				errorMessage = "Вышел за границу поля"
			}
		}
	} else {
		errorMessage = "Неверно выбрана начальная точка"
	}
	return fmt.Errorf("%s", errorMessage)
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

func CheckQuantity(s ship.Ship, ships *[]ship.Ship) error {
	count := [4]int{4, 3, 2, 1}
	for _, sh := range *ships {
		switch sh.Length {
		case 1:
			count[0]--
			break
		case 2:
			count[1]--
			break
		case 3:
			count[2]--
			break
		case 4:
			count[3]--
			break
		}
	}
	if count[s.Length-1] == 0 {
		return fmt.Errorf("%s %d\n", "Слишком много кораблей с длиной", s.Length)
	}
	return nil
}

func GetShipsByID(id int, ships *[]ship.Ship) []ship.Ship {
	var playerShips []ship.Ship

	for _, sh := range *ships {
		if sh.Player == id {
			playerShips = append(playerShips, sh)
		}
	}
	return playerShips
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

func (b *BattleField) CheckShot(newShot *shot.Shot) error {
	if b.Field[newShot.X][newShot.Y] == "O" || b.Field[newShot.X][newShot.Y] == "_" {
		return nil
	}
	return fmt.Errorf("%s", "Сюда нельзя стрелять")
}

func (b BattleField) CheckHit(newShot *shot.Shot, ships *[]ship.Ship) bool {
	var liveShips []ship.Ship
	var result bool = false

	for i := 0; i < len(*ships); i++ {

		for j := 0; j < (*ships)[i].Length; j++ {
			if (*ships)[i].Direction == 0 {
				if newShot.X == (*ships)[i].StartX && newShot.Y == (*ships)[i].StartY+j {
					(*ships)[i].LivePoints--
					result = true
				}
			} else {
				if newShot.X == (*ships)[i].StartX+j && newShot.Y == (*ships)[i].StartY {
					(*ships)[i].LivePoints--
					result = true
				}
			}
		}

		if (*ships)[i].LivePoints > 0 {
			liveShips = append(liveShips, (*ships)[i])
		}

	}

	*ships = liveShips
	return result
}
