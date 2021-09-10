package lib

import (
	"fmt"
	"time"
)

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

func Log(message string) {
	data := time.Now().Format("02-01-2006 15:04:05")
	fmt.Printf("%s | %s\n", data, message)
}
