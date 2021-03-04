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
	myField    Field
	enemyField Field
}

func (b BattleField) CreateField() BattleField {
	for i := 0; i < FIELD_HEIGHT; i++ {
		for j := 0; j < FIELD_WIDTH; j++ {
			if i == 0 || i == FIELD_HEIGHT-1 {
				b.myField[i][j] = "*"
				b.enemyField[i][j] = "*"
				continue
			}
			if j == 0 || j == FIELD_WIDTH-1 {
				b.myField[i][j] = "*"
				b.enemyField[i][j] = "*"
			} else {
				b.myField[i][j] = "_"
				b.enemyField[i][j] = "_"
			}
		}
	}
	return b
}

func (b BattleField) DrawField() {

	fmt.Printf(LETTER_STRING)
	for i := 1; i < FIELD_HEIGHT-1; i++ {

		//My field drawing
		if i != FIELD_HEIGHT-2 {
			fmt.Printf(" %d", i)
		} else {
			fmt.Printf("%d", i)
		}
		for j := 1; j < FIELD_WIDTH-1; j++ {
			if j != FIELD_WIDTH-2 {
				fmt.Printf("|%s", b.myField[i][j])
			} else {
				fmt.Printf("|%s|", b.myField[i][j])
			}
		}
		fmt.Printf("\t\t")

		//Enemy field drawing
		if i != FIELD_WIDTH-2 {
			fmt.Printf(" %d", i)
		} else {
			fmt.Printf("%d", i)
		}
		for j := 1; j < FIELD_HEIGHT-1; j++ {
			if j != FIELD_HEIGHT-2 {
				fmt.Printf("|%s", b.enemyField[i][j])
			} else {
				fmt.Printf("|%s|", b.enemyField[i][j])
			}
		}
		fmt.Println()
	}
}

func (b BattleField) DrawShip(s ship.Ship) BattleField {
	for i := 0; i < s.Length; i++ {
		b.CheckShip(s)
		if s.Player {
			if s.Direction == 0 {
				b.myField[s.StartY][s.StartX+i] = "O"
				b.myField[s.StartY+1][s.StartX+i] = "*"
				b.myField[s.StartY-1][s.StartX+i] = "*"
				b.myField[s.StartY][s.StartX-1] = "*"
				b.myField[s.StartY][s.StartX+s.Length] = "*"
				b.myField[s.StartY+1][s.StartX+s.Length] = "*"
				b.myField[s.StartY-1][s.StartX+s.Length] = "*"
				b.myField[s.StartY+1][s.StartX-1] = "*"
				b.myField[s.StartY-1][s.StartX-1] = "*"
			} else if s.Direction == 1 {
				b.myField[s.StartY+i][s.StartX] = "O"
				b.myField[s.StartY+i][s.StartX+1] = "*"
				b.myField[s.StartY+i][s.StartX-1] = "*"
				b.myField[s.StartY-1][s.StartX] = "*"
				b.myField[s.StartY+s.Length][s.StartX] = "*"
				b.myField[s.StartY+s.Length][s.StartX+1] = "*"
				b.myField[s.StartY+s.Length][s.StartX-1] = "*"
				b.myField[s.StartY-1][s.StartX+1] = "*"
				b.myField[s.StartY-1][s.StartX-1] = "*"
			}
		} else {
			if s.Direction == 0 {
				b.enemyField[s.StartY][s.StartX+i] = "O"
				b.enemyField[s.StartY+1][s.StartX+i] = "*"
				b.enemyField[s.StartY-1][s.StartX+i] = "*"
				b.enemyField[s.StartY][s.StartX-1] = "*"
				b.enemyField[s.StartY][s.StartX+s.Length] = "*"
				b.enemyField[s.StartY+1][s.StartX+s.Length] = "*"
				b.enemyField[s.StartY-1][s.StartX+s.Length] = "*"
				b.enemyField[s.StartY+1][s.StartX-1] = "*"
				b.enemyField[s.StartY-1][s.StartX-1] = "*"
			} else if s.Direction == 1 {
				b.enemyField[s.StartY+i][s.StartX] = "O"
				b.enemyField[s.StartY+i][s.StartX+1] = "*"
				b.enemyField[s.StartY+i][s.StartX-1] = "*"
				b.enemyField[s.StartY-1][s.StartX] = "*"
				b.enemyField[s.StartY+s.Length][s.StartX] = "*"
				b.enemyField[s.StartY+s.Length][s.StartX+1] = "*"
				b.enemyField[s.StartY+s.Length][s.StartX-1] = "*"
				b.enemyField[s.StartY-1][s.StartX+1] = "*"
				b.enemyField[s.StartY-1][s.StartX-1] = "*"
			}
		}
	}
	b.DrawField()
	return b
}

func (b BattleField) DrawShot(Player bool, ShotX, ShotY int, Result int) BattleField {
	if Player {
		if Result == 0 {
			b.enemyField[ShotX][ShotY] = "*"
		} else {
			b.enemyField[ShotX][ShotY] = "X"
		}
	} else {
		if Result == 0 {
			b.myField[ShotX][ShotY] = "*"
		} else {
			b.myField[ShotX][ShotY] = "X"
		}
	}
	b.DrawField()
	return b
}

func (b BattleField) ClearField() {
	b.myField = Field{}
	b.enemyField = Field{}
	b.DrawField()
}

func (b BattleField) CheckShip(s ship.Ship) (bool, error) {

	errorMessage := "Начальное сообщение"
	var turnCheck Field

	if s.Player {
		turnCheck = b.myField
	} else {
		turnCheck = b.enemyField
	}

	if turnCheck[s.StartY][s.StartX] == "_" { //В начальной точке нет другого корабля
		if s.Direction == 0 {
			if s.StartY+s.Length < FIELD_HEIGHT { //Проверка выхода за границы поля
				if turnCheck[s.StartY+s.Length][s.StartX] != "_" { //Проверка доступности клетки в конечной точке
					errorMessage = "Уперся в *"
				} else {
					return true, nil
				}
			} else {
				errorMessage = "Вышел за границу"
			}
		} else {
			if s.StartX+s.Length < FIELD_WIDTH {
				if turnCheck[s.StartY][s.StartX+s.Length] != "_" {
					errorMessage = "Уперся в *"
				} else {
					return true, nil
				}
			} else {
				errorMessage = "Вышел за границу"
			}
		}
	} else {
		errorMessage = "Первое условие"
	}
	return false, fmt.Errorf("%s", errorMessage)
}

func (b BattleField) CheckShot(Player bool, ShotX, ShotY int) error {

	if Player {
		if b.enemyField[ShotX][ShotY] == "O" || b.enemyField[ShotX][ShotY] == "*" {
			return nil
		}
	} else {
		if b.myField[ShotX][ShotY] == "O" || b.myField[ShotX][ShotY] == "*" {
			return nil
		}
	}
	return fmt.Errorf("%s", "Сюда нельзя стрелять")
}

func (b BattleField) CheckHit(Player bool, ShotX, ShotY int, ships *[]ship.Ship) bool {

	var result bool = false

	var newShips, myShips, enemyShips []ship.Ship

	newShips = *ships

	for i := 0; i < len(newShips); i++ {
		if newShips[i].Player == !Player {
			enemyShips = append(enemyShips, newShips[i])
		} else {
			myShips = append(myShips, newShips[i])
		}
	}

	for i := 0; i < len(enemyShips); i++ {

		for j := 0; j < enemyShips[i].Length; j++ {
			if enemyShips[i].Direction == 0 {
				if ShotX == enemyShips[i].StartX && ShotY == enemyShips[i].StartY+j {
					enemyShips[i].LivePoints--
					result = true
				}
			} else {
				if ShotX == enemyShips[i].StartX+j && ShotY == enemyShips[i].StartY {
					enemyShips[i].LivePoints--
					result = true
				}
			}
		}

		if enemyShips[i].LivePoints > 0 {
			myShips = append(myShips, enemyShips[i])
			fmt.Println(myShips)
		}

	}

	*ships = myShips
	return result
}
