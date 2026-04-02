package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/mike-the-math-man/internal/config"
	"github.com/mike-the-math-man/internal/database"
)

func main() {
	updatedCfg := config.Read()

	db, err := sql.Open("postgres", updatedCfg.Db_url)
	dbQueries := database.New(db)
	new_commands := commands{
		map[string]func(*state, command) error{},
	}
	new_state := state{
		dbQueries,
		&updatedCfg,
	}
	new_commands.register("login", handlerLogin)
	new_commands.register("register", handlerRegister)
	new_commands.register("reset", handlerReset)
	new_commands.register("users", handlerListUsers)
	new_commands.register("agg", aggregator_list)
	new_commands.register("addfeed", middlewareLoggedIn(add_feed))
	new_commands.register("feeds", middlewareLoggedIn(get_feeds_list))
	new_commands.register("follow", middlewareLoggedIn(follow_feed))
	new_commands.register("following", middlewareLoggedIn(follow_feeds_list))
	new_commands.register("unfollow", middlewareLoggedIn(unfollowFeed))
	args := os.Args
	if len(args) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	err = new_commands.run(&new_state, command{args[1], args[2:]})
	if err != nil {
		fmt.Println("Error")
		os.Exit(1)
	}
	//fmt.Printf("DbURL: %s\n", updatedCfg.Db_url)
	//fmt.Printf("User: %s\n", updatedCfg.Current_user_name)
}
