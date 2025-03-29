package main

import (
	"fmt"
	"log"

	"example/greetings"
	// "rsc.io/quote"
)

func main() {
	msg, err := greetings.FuckYou("Gleb")
	if err != nil {
		fmt.Println("Put not empty string")
		log.Fatal("Ты еблан")
	}
	fmt.Println(msg)
}
