package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/hramov/battleship_server/pkg/battlefield"
	connection "github.com/hramov/battleship_server/pkg/connection"
	"github.com/hramov/battleship_server/pkg/ship"
	"github.com/hramov/battleship_server/pkg/shot"
	"github.com/hramov/battleship_server/pkg/utils"
)

var ships []ship.Ship
var clients = make(map[int]connection.Client)
var battlefields = make(map[int]battlefield.BattleField)
var shots = make(map[int]shot.Shot)
var turn bool = false

func main() {

	s := connection.Execute("tcp", "127.0.0.1", "5000")
	handlers := make(map[string]func(client *connection.Client, data string))

	handlers["connect"] = func(client *connection.Client, data string) {
		s.Emit(client, "connect", "")
		s.Emit(client, "whoami", strconv.Itoa((*client).ID))

		if len(clients) < 2 {
			clients[(*client).ID] = *client
		}

		if len(clients)%2 == 0 { // Есть пара соперников
			utils.Log("Start game")
			Roll()                       // Задаем случайное значение переменной turn
			for _, cl := range clients { // Перебираем всех игроков
				cl.Turn = turn // Назначаем полю Turn клиента значение переменной turn
				for i := 1; i <= len(clients); i++ {
					if i != cl.ID { // Если это не текущий клиент
						cl.EnemyID = i // назначаем ID оппонента
					}
				}
				clients[cl.ID] = cl
				turn = !turn // Переключаем переменную turn
				uData, _ := json.Marshal(cl)
				s.Emit(&cl, "update", string(uData))
			}

			s.BroadCast(&clients, "drawField", "")
			s.BroadCast(&clients, "placeShip", "")
		}
	}

	handlers["sendShip"] = func(client *connection.Client, data string) {
		client = FindClientByID((*client).ID)
		sh := ship.Ship{}
		sh.Player = (*client).ID
		sh.ID = GetShipID()
		json.Unmarshal([]byte(data), &sh)
		b := FindFieldByClientID((*client).ID)
		clientShips := FindShipsByClientID(client.ID)
		err := b.CheckShip(sh, &clientShips)
		if err != nil {
			fmt.Println(err)
			s.Emit(client, "wrongShip", err.Error())
			s.Emit(client, "placeShip", err.Error())
		} else {
			b.CreateShip(sh)
			battlefields[(*client).ID] = b
			ships = append(ships, sh)
			s.Emit(client, "rightShip", "Ваш корабль успешно добавлен! Осталось: "+strconv.Itoa(10-len(FindShipsByClientID(client.ID))))
			shipData, _ := json.Marshal(b)
			s.Emit(client, "drawField", string(shipData))

			if len(FindShipsByClientID(client.ID)) < 1 {
				s.Emit(client, "placeShip", "")
			} else {
				s.Emit(client, "makeShot", strconv.FormatBool((*client).Turn))
			}
		}
	}

	handlers["shot"] = func(client *connection.Client, data string) {

		client = FindClientByID(client.ID)
		enemy := FindClientByID(client.EnemyID)

		enemyShips := FindShipsByClientID(enemy.ID)

		newShot := shot.Shot{}
		json.Unmarshal([]byte(data), &newShot)

		b := FindFieldByClientID(enemy.ID)

		if err := b.CheckShot(&newShot); err != nil {
			utils.Log(err.Error())
			s.Emit(client, "makeShot", strconv.FormatBool(client.Turn))
			return
		}

		shots[client.ID] = newShot

		clientB := FindFieldByClientID(client.ID)
		enemyB := FindFieldByClientID(enemy.ID)

		shipID, err := newShot.CheckHit(&enemyShips)

		if err != nil {

			clientB.CreateShot(true, true, newShot)
			enemyB.CreateShot(false, true, newShot)

			battlefields[client.ID] = clientB
			battlefields[enemy.ID] = enemyB

			if playerID, isLose := UpdateShip(shipID); isLose {
				fmt.Printf("%d %s\n", playerID, "програл!")
				// s.BroadCast(&clients, "restart", "")
				// Restart()
				os.Exit(0)
				return
			}

			clientData, _ := json.Marshal(clientB)
			enemyData, _ := json.Marshal(enemyB)

			s.Emit(client, "updateField", string(clientData))
			s.Emit(enemy, "updateField", string(enemyData))

			s.Emit(client, "makeShot", strconv.FormatBool(client.Turn))
			s.Emit(enemy, "makeShot", strconv.FormatBool(enemy.Turn))

		} else {

			clientB.CreateShot(true, false, newShot)
			enemyB.CreateShot(false, false, newShot)

			clientData, _ := json.Marshal(clientB)
			enemyData, _ := json.Marshal(enemyB)

			battlefields[client.ID] = clientB
			battlefields[enemy.ID] = enemyB

			client.Turn = !client.Turn
			enemy.Turn = !enemy.Turn

			clients[client.ID] = *client
			clients[client.EnemyID] = *enemy

			s.Emit(client, "updateField", string(clientData))
			s.Emit(enemy, "updateField", string(enemyData))

			s.Emit(enemy, "makeShot", strconv.FormatBool(enemy.Turn))
			s.Emit(client, "makeShot", strconv.FormatBool(client.Turn))

		}
	}

	handlers["disconnect"] = func(client *connection.Client, data string) {
		client = FindClientByID((*client).ID)
		utils.Log("User: " + strconv.Itoa(client.ID) + " disconnected!")
		client.Socket.Close()
		delete(clients, client.ID)
		delete(battlefields, client.ID)
		utils.Log("Users remain: " + strconv.Itoa(len(clients)))
	}

	maintainConnections(s, &handlers)
}

func maintainConnections(s *connection.Server, handlers *(map[string]func(client *connection.Client, data string))) {
	for {
		conn, _ := s.Ln.Accept()
		ID := len(clients) + 1
		client := connection.Client{ID, 0, conn, make(chan string), make(chan string), false}
		clients[ID] = client

		utils.Log("User " + strconv.Itoa(client.ID) + " connected!")

		b := battlefield.BattleField{}
		b.CreateField()
		battlefields[client.ID] = b

		go s.Speak(&client)
		go s.Listen(&client)
		go s.On(&client, handlers)

		client.From <- "connect|" + strconv.Itoa(client.ID)
	}
}

func FindFieldByClientID(ID int) battlefield.BattleField {
	var b battlefield.BattleField
	for id, battlefield := range battlefields {
		if id == ID {
			b = battlefield
		}
	}
	return b
}

func FindShipsByClientID(ID int) []ship.Ship {
	var clientShips []ship.Ship
	for _, sh := range ships {
		if sh.Player == ID {
			clientShips = append(clientShips, sh)
		}
	}
	return clientShips
}

func FindClientByID(ID int) *connection.Client {
	cl := clients[ID]
	return &cl
}

func Roll() {
	rand.Seed(time.Now().UnixNano())
	player := rand.Intn(2)
	if player == 0 {
		turn = true
	} else {
		turn = false
	}
}

func GetShipID() int {
	return len(ships) + 1
}

func FindShip(shipID int) (int, int, ship.Ship) {
	var target ship.Ship
	var targetID int
	var playerID int

	for id, sh := range ships {
		if sh.ID == shipID {
			targetID = id
			playerID = sh.Player
			target = sh
			break
		}
	}
	return targetID, playerID, target
}

func UpdateShip(shipID int) (int, bool) {

	var updatedShips []ship.Ship

	id, playerID, target := FindShip(shipID)

	target.LivePoints--
	ships[id] = target

	for _, sh := range ships {
		if sh.LivePoints > 0 {
			updatedShips = append(updatedShips, sh)
		}
	}

	ships = updatedShips

	if lose := CheckLose(playerID); lose {
		return playerID, true
	}
	return playerID, false
}

func CheckLose(playerID int) bool {
	clientShips := FindShipsByClientID(playerID)
	if len(clientShips) == 0 {
		return true
	}
	return false
}

func Restart() {
	ships = make([]ship.Ship, 0)
	clients = make(map[int]connection.Client)
	battlefields = make(map[int]battlefield.BattleField)
	shots = make(map[int]shot.Shot)
	turn = false
}
