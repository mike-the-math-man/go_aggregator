package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/mike-the-math-man/internal/config"
	"github.com/mike-the-math-man/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

type command struct {
	name      string
	arguments []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		fmt.Println("not enough arguments in handlerlogin")
		return fmt.Errorf("Nooooo")
	}
	user, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		fmt.Println("user does not exist - pls register")
		return err
	}
	s.config.SetUser(*s.config, user.Name)

	//s.config.Current_user_name = cmd.name
	fmt.Println("user set")
	return nil
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	err := c.commands[cmd.name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		fmt.Println("not enough arguments in register - supply name")
		return fmt.Errorf("Nooooo")
	}
	user_id := uuid.New()
	name := cmd.arguments[0]
	_, err := s.db.GetUser(context.Background(), name)
	params := database.CreateUserParams{ID: user_id, CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: name}
	if err != nil {
		user, err := s.db.CreateUser(context.Background(), params)
		if err != nil {
			fmt.Println("error creating user", err)
			return err
		}
		s.config.SetUser(*s.config, user.Name)
		fmt.Println("User created")
		return nil
	} else {
		os.Exit(1)
		return err
	}
}
