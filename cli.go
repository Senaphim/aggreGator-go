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

	_, err := s.db.GetUserByName(context.Background(), c.args[0])
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
		fmtErr := fmt.Errorf("Error creating new user:\n%v", err)
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
		fmtErr := fmt.Errorf("Error deleting users:\n%v", err)
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
		fmtErr := fmt.Errorf("Error retrieving users:\n%v", err)
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

func handlerAgg(_ *state, c command) error {
	if len(c.args) != 0 {
		return errors.New("Incorrect number of arguments supplied. Expecting 0")
	}

	feed, err := FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		fmtErr := fmt.Errorf("Error aggregating:\n%v", err)
		return fmtErr
	}

	fmt.Printf("%v\n", feed)
	return nil
}

func handlerAddFeed(s *state, c command, user database.User) error {
	if len(c.args) != 2 {
		return errors.New("Incorrect number of arguments supplied. Expecting 2")
	}

	fd := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().Local(),
		UpdatedAt: time.Now().Local(),
		Name:      c.args[0],
		Url:       c.args[1],
		UserID:    user.ID,
	}

	newFeed, err := s.db.CreateFeed(context.Background(), fd)
	if err != nil {
		fmtErr := fmt.Errorf("Error creating new feed:\n%v", err)
		return fmtErr
	}

	fmt.Printf("%v\n", newFeed)

	handlerFollow(s, command{"follow", []string{c.args[1]}}, user)

	return nil
}

func handlerFeeds(s *state, c command) error {
	if len(c.args) != 0 {
		return errors.New("Incorrect number of arguments supplied. Expecting 0")
	}

	fds, err := s.db.AllFeeds(context.Background())
	if err != nil {
		fmtErr := fmt.Errorf("Error fetching feeds:\n%v", err)
		return fmtErr
	}

	for _, fd := range fds {
		fmt.Printf("%v\n", fd.UserID)
		usr, err := s.db.GetUserByUuid(context.Background(), fd.UserID)
		if err != nil {
			fmtErr := fmt.Errorf("Error fetching feed %v creation user:\n%v", fd.Name, err)
			return fmtErr
		}
		fmt.Printf("Name: %v, URL: %v, Created By: %v\n", fd.Name, fd.Url, usr.Name)
	}

	return nil
}

func handlerFollow(s *state, c command, usr database.User) error {
	if len(c.args) != 1 {
		return errors.New("Incorrect number of arguments supplied. Expecting 1")
	}

	fd, err := s.db.GetFeedByUrl(context.Background(), c.args[0])
	if err != nil {
		fmtErr := fmt.Errorf("Error fetching feed:\n%v", err)
		return fmtErr
	}

	follow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().Local(),
		UpdatedAt: time.Now().Local(),
		UserID:    usr.ID,
		FeedID:    fd.ID,
	}
	newFollow, err := s.db.CreateFeedFollow(context.Background(), follow)
	if err != nil {
		fmtErr := fmt.Errorf("Error following feed:\n%v", err)
		return fmtErr
	}

	fmt.Printf("Followed feed: %v User: %v\n", newFollow.Name_2, newFollow.Name)

	return nil
}

func handlerFollowing(s *state, c command, usr database.User) error {
	if len(c.args) != 0 {
		return errors.New("Incorrect number of arguments supplied. Expecting 0")
	}

	following, err := s.db.GetFeedFollowsForUser(context.Background(), usr.ID)
	if err != nil {
		fmtErr := fmt.Errorf("Error fetching followed feeds:\n%v", err)
		return fmtErr
	}

	for _, follow := range following {
		fmt.Printf("%v\n", follow.Name)
	}
	return nil
}
