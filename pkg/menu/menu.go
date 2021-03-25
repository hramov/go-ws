package menu

import (
	"fmt"
	"os"

	"github.com/hramov/battleship_server/pkg/gameloop"
)

func create() []string {
	menuItems := []string{
		"Начать игру",
		"О создателях",
		"Выйти",
	}
	return menuItems
}

func draw(menuItems []string) {
	fmt.Println("Hello! Welcome to BattleShips! Here you have a menu.")
	for i := 0; i < len(menuItems); i++ {
		fmt.Printf("%d) %s\n", i+1, menuItems[i])
	}
}

func MainMenu() {
	menuItems := create()
	draw(menuItems)

	var menuChose int = 0
	fmt.Scanf("%d", &menuChose)
	switch menuChose {
	case 1:
		gameloop.Start()
		break
	case 2:
		fmt.Println("Khramov Sergey")
		break
	case 3:
		os.Exit(0)
	default:
		fmt.Println("You didn't choose anything")
	}
}
