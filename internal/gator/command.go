// Package gator has the commands of this application
package gator

import (
	"errors"
	"fmt"

	"github.com/Uami-11/blog-aggregator/internal/config"
)

type Command struct {
	Name      string
	Arguments []string
}

type Commands struct {
	Comms map[string]func(*config.State, Command) error
}

func HandlerLogin(state *config.State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return errors.New("usage: gator login <username>")
	}

	err := state.Conf.SetUser(cmd.Arguments[0])
	if err != nil {
		return err
	}

	fmt.Println("User has been set!")
	return nil
}

func (comms *Commands) Run(state *config.State, cmd Command) error {
	fun, ok := comms.Comms[cmd.Name]
	if ok {
		return fun(state, cmd)
	}
	return errors.New("could not find command")
}

func (comms *Commands) Register(name string, f func(*config.State, Command) error) {
	comms.Comms[name] = f
}

func run(state *config.State, cmd Command) error {
	if state.Conf.DBURL != "" {
		return errors.New("gator failed to run")
	}

	fmt.Println("Gator ran!")
	return nil
}
