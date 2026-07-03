// Package cmd...
package cmd

import (
	"errors"
	"fmt"

	"github.com/Uami-11/blog-aggregator/internal/config"
)

type Command struct {
	Name string
	Args []string
}

type State struct {
	TheConfig *config.Config
}

func HandlerLogin(s *State, comm Command) error {
	if len(comm.Args) == 0 {
		return errors.New("the login information should contain a username argument")
	}

	err := s.TheConfig.SetUser(comm.Args[0])
	if err != nil {
		return err
	}

	fmt.Println("User successfully logged in!")

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
