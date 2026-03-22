package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func Read() Config {
	home_dir, err := os.UserHomeDir()
	file_path := home_dir + "/.gatorconfig.json"
	config, err := os.ReadFile(file_path)
	if err != nil {
		return Config{}
	}
	var config_struct Config
	err = json.Unmarshal(config, &config_struct)
	if err != nil {
		return Config{}
	}
	return config_struct
}

func (Config) SetUser(c Config, user_name string) {
	home_dir, err := os.UserHomeDir()
	file_path := home_dir + "/.gatorconfig.json"
	c.Current_user_name = user_name
	json_data, err := json.Marshal(c)
	if err != nil {
		return
	}
	err = os.WriteFile(file_path, json_data, 0644)
	fmt.Println("File written successfully.")
	//return nil
}

type state struct {
	config *Config
}

type command struct {
	name      string
	arguments []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("cmd slice empty")
	}
	s.config.Current_user_name = cmd.name
	fmt.Println("user set")
	return nil
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	c.commands[cmd.name](s, cmd)
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}

func Add_to_state(s *state, c *Config) {
	s.config = c
}
