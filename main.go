package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/senaphim/aggreGator-go/internal/config"
	"github.com/senaphim/aggreGator-go/internal/database"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Println(fmt.Errorf("Error reading config file:\n%v", err))
		os.Exit(1)
	}

	db, err := sql.Open("postgres", *conf.DbUrl)
	if err != nil {
		fmtErr := fmt.Errorf("Error connecting to database:\n%v", err)
		fmt.Println(fmtErr)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	st := &state{
		db:   dbQueries,
		conf: &conf,
	}

	var cmds commands
	cmds.cmd = make(map[string]func(*state, command) error)
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", loggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", loggedIn(handlerFollow))
	cmds.register("following", loggedIn(handlerFollowing))

	argSlice := os.Args
	if len(argSlice) < 2 {
		fmt.Println("No command entered, exiting")
		os.Exit(1)
	}

	var cmd command
	cmd.name = argSlice[1]
	cmd.args = argSlice[2:]
	if err := cmds.run(st, cmd); err != nil {
		fmtErr := fmt.Errorf("Error:\n%v", err)
		fmt.Println(fmtErr)
		os.Exit(1)
	}
}
