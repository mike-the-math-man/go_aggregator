package main

import (
	"fmt"
	"os"

	"github.com/mike-the-math-man/internal/config"
)

func main() {
	updatedCfg := config.Read()
	new_state := state{
		&updatedCfg,
	}
	new_commands := commands{
		map[string]func(*state, command) error{},
	}
	new_commands.register("login", handlerLogin)
	args := os.Args
	if len(args) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	err := new_commands.run(&new_state, command{args[1], args[2:]})
	if err != nil {
		fmt.Println("Error")
		os.Exit(1)
	}
	//fmt.Printf("DbURL: %s\n", updatedCfg.Db_url)
	//fmt.Printf("User: %s\n", updatedCfg.Current_user_name)
}
