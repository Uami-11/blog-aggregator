// Package cmd...
package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Uami-11/blog-aggregator/internal/config"
	"github.com/Uami-11/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type Command struct {
	Name string
	Args []string
}

type State struct {
	DB        *database.Queries
	TheConfig *config.Config
}

func HandlerLogin(s *State, comm Command) error {
	if len(comm.Args) == 0 {
		return errors.New("the login information should contain a username argument")
	}

	user, err := s.DB.GetUser(context.Background(), comm.Args[0])
	if err != nil {
		return errors.New("user does not exist")
	}

	err = s.TheConfig.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Println("User successfully logged in!")

	return nil
}

func HandlerRegister(s *State, comm Command) error {
	if len(comm.Args) == 0 {
		return errors.New("registration information should contain a username argument")
	}

	var user database.CreateUserParams

	user.ID = uuid.New()

	currentTime := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	user.CreatedAt = currentTime
	user.UpdatedAt = currentTime

	if _, err := s.DB.GetUser(context.Background(), comm.Args[0]); err == nil {
		return errors.New("that user already exists")
	}

	user.Name = comm.Args[0]

	_, err := s.DB.CreateUser(context.Background(), user)
	if err != nil {
		return errors.New("error creating user")
	}

	s.TheConfig.CurrentUserName = user.Name

	fmt.Printf("user %s has been created!\n", user.Name)
	fmt.Printf("UUID: %s\n", user.ID)
	fmt.Printf("Created At: %v\n", user.CreatedAt)
	fmt.Printf("Updated At: %v\n", user.UpdatedAt)

	return nil
}

type Commands struct {
	Comms map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, comm Command) error {
	handler, ok := c.Comms[comm.Name]
	if !ok {
		return errors.New("command does not exist")
	}

	err := handler(s, comm)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) Register(name string, f func(s *State, comm Command) error) {
	c.Comms[name] = f
}
