package main

import (
	"fmt"
	"os/user"

	"github.com/k20ku/monkey/repl"
)

func main() {
	// greeting
	if user, err := user.Current(); err == nil {
		fmt.Printf("Hello %s! ", user.Username)
	} else {
		fmt.Print("Hello! ")
	}
	fmt.Println("This is the Monkey programming language!")
	fmt.Println("Feel free to type in commands")
	fmt.Println("Ctrl+D to exit!")

	for i, mode := range repl.AllReplModes {
		if i == 0 {
			fmt.Printf("'%s'", mode.String())
			continue
		}
		fmt.Printf(", '%s'", mode.String())
	}
	fmt.Println()
	// start repl
	repl.Start()

	// greeting
	println("Thank you. See you!")
}
