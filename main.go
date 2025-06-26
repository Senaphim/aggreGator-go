package main

import (
	"fmt"
	"os"

	"github.com/senaphim/aggreGator-go/internal/config"
)

func main() {
	var st state
	conf, err := config.Read()
	if err != nil {
		fmt.Println(fmt.Errorf("Error reading config file:\n%v", err))
		os.Exit(1)
	}
	st.conf = &conf

	var cmds commands
	cmds.cmd = make(map[string]func(*state, command) error)
	cmds.register("login", handlerLogin)

	argSlice := os.Args
	if len(argSlice) < 2 {
		fmt.Println("No command entered, exiting")
		os.Exit(1)
	}

	var cmd command
	cmd.name = argSlice[1]
	cmd.args = argSlice[2:]
	if err := cmds.run(&st, cmd); err != nil {
		fmtErr := fmt.Errorf("Error: %v", err)
		fmt.Println(fmtErr)
		os.Exit(1)
	}

	return
}
