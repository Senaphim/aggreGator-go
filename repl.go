package main

import (
	"errors"
	"fmt"

	"github.com/senaphim/aggreGator-go/internal/config"
)

type state struct {
	conf *config.Config
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
	if len(c.args) == 0 {
		return errors.New("Please provide a username when logging in")
	}

	if err := s.conf.SetUser(c.args[0]); err != nil {
		fmtErr := fmt.Errorf("Error logging in:\n%v", err)
		return fmtErr
	}

	fmt.Println("Login successful!")

	return nil
}

