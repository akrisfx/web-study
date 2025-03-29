package greetings

import (
	"errors"
	"fmt"
	"math/rand"
)

// Hello returns a greeting for the named person.
func Hello(name string) string {
	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
	return message
}

func FuckYou(name string) (string, error) {
	if name == "" {
		return "", errors.New("empty string")
	}
	return fmt.Sprintf("%v, %v", RandomFuckYou(), name), nil
}

func RandomFuckYou() string {
	formats := []string{
		"Sosi huy eblan",
		"Idi nahuy dolboeb",
		"fuck off and suck balls",
		"hi, luckyboy",
	}
	return formats[rand.Intn(len(formats))]
}
