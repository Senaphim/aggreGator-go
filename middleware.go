package main

import (
	"context"
	"fmt"

	"github.com/senaphim/aggreGator-go/internal/database"
)

func loggedIn(
	handler func(s *state, c command, user database.User) error,
) func(*state, command) error {

	return func(s *state, c command) error {
		user, err := s.db.GetUserByName(context.Background(), *s.conf.CurrentUserName)
		if err != nil {
			fmtErr := fmt.Errorf("Error getting user from database:\n%v", err)
			return fmtErr
		}
		return handler(s, c, user)
	}
}
