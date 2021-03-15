package utils

func Parser(x int) string {
	var a string
	switch x {
	case 1:
		a = "A"
		break
	case 2:
		a = "Б"
		break
	case 3:
		a = "В"
		break
	case 4:
		a = "Г"
		break
	case 5:
		a = "Д"
		break
	case 6:
		a = "Е"
		break
	case 7:
		a = "Ж"
		break
	case 8:
		a = "З"
		break
	case 9:
		a = "И"
		break
	case 10:
		a = "К"
		break
	}
	return a
}

// func FindClientBySocket(s socketio.Conn, sockets *[]socketio.Conn) string {
// 	for _, client := range *sockets {
// 		if
// 	}
// }

func Split(message string, delim string) (string, string) {
	var breakPosition int = 0
	var charArray []string
	var eventString, dataString string
	for i, char := range message {
		charArray = append(charArray, string(char))
		if string(char) == delim {
			breakPosition = i
		}
	}
	event := charArray[:breakPosition]
	data := charArray[breakPosition+1:]
	for i := 0; i < len(event); i++ {
		eventString += string(event[i])
	}
	for i := 0; i < len(data); i++ {
		dataString += string(data[i])
	}
	return eventString, dataString
}
