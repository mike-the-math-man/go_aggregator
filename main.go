package main

import (
	"fmt"

	"github.com/mike-the-math-man/internal/config"
)

func main() {
	config_file := config.Read()
	config_file.SetUser(config_file, "Mike")
	updatedCfg := config.Read()
	fmt.Printf("DbURL: %s\n", updatedCfg.Db_url)
	fmt.Printf("User: %s\n", updatedCfg.Current_user_name)
}
