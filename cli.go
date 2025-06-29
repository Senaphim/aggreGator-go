package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/senaphim/aggreGator-go/internal/config"
	"github.com/senaphim/aggreGator-go/internal/database"
)

type state struct {
	conf *config.Config
	db   *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmd map[string]func(*state, command) error
}

func (cs *commands) run(s *state, c command) error {
	if _, ok := cs.cmd[c.name]; !ok {
		return errors.New("Command does not exist")
	}

	if err := cs.cmd[c.name](s, c); err != nil {
		fmtErr := fmt.Errorf("Error running command %v:\n%v", c.name, err)
		return fmtErr
	}

	return nil
}

func (cs *commands) register(name string, f func(*state, command) error) {
	if cs.cmd == nil {
		cs.cmd = make(map[string]func(*state, command) error)
	}
	cs.cmd[name] = f
}

func handlerLogin(s *state, c command) error {
	if len(c.args) != 1 {
		return errors.New("Incorrect number of arguments supplied. Expecting 1")
	}

	_, err := s.db.GetUser(context.Background(), c.args[0])
	if err != nil {
		fmtErr := fmt.Errorf("Error logging in, user not found:\n%v", err)
		return fmtErr
	}

	if err := s.conf.SetUser(c.args[0]); err != nil {
		fmtErr := fmt.Errorf("Error logging in:\n%v", err)
		return fmtErr
	}

	fmt.Println("Login successful!")

	return nil
}

func handlerRegister(s *state, c command) error {
	if len(c.args) != 1 {
		return errors.New("Incorrect number of arguments supplied. Expecting 1")
	}

	usr := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().Local(),
		UpdatedAt: time.Now().Local(),
		Name:      c.args[0],
	}

	newUsr, err := s.db.CreateUser(context.Background(), usr)
	if err != nil {
		fmtErr := fmt.Errorf("Error creating new user: %v", err)
		return fmtErr
	}

	s.conf.SetUser(newUsr.Name)

	fmt.Println("New user created successfully")
	fmt.Printf(
		"Id: %v, created at: %v, updated at: %v, name: %v\n",
		newUsr.ID,
		newUsr.CreatedAt,
		newUsr.UpdatedAt,
		newUsr.Name,
	)

	return nil
}

func handlerReset(s *state, c command) error {
	if len(c.args) != 0 {
		return errors.New("Incorrect number of arguments supplied. Expecting 0")
	}

	if err := s.db.DeleteAll(context.Background()); err != nil {
		fmtErr := fmt.Errorf("Error deleting users: %v", err)
		return fmtErr
	}

	return nil
}

func handlerUsers(s *state, c command) error {
	if len(c.args) != 0 {
		return errors.New("Incorrect number of arguments supplied. Expecting 0")
	}

	users, err := s.db.AllUsers(context.Background())
	if err != nil {
		fmtErr := fmt.Errorf("Error retrieving users: %v", err)
		return fmtErr
	}

	for _, user := range users {
		printString := fmt.Sprintf("%v", user.Name)
		if user.Name == *s.conf.CurrentUserName {
			printString += " (current)"
		}
		fmt.Println(printString)
	}

	return nil
}
