package main

import (
	"fmt"
	"os"

	"github.com/Uami-11/blog-aggregator/internal/config"
	"github.com/Uami-11/blog-aggregator/internal/gator"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read the configuration file: %v\n", err)
	}

	err = conf.SetUser("Uami")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set you as the current user: %v\n", err)
	}

	conf, err = config.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read the configuration file: %v\n", err)
	}

	state := &config.State{
		Conf: &conf,
	}

	com := gator.Commands{
		Comms: make(map[string]func(*config.State, gator.Command) error),
	}

	com.Register("login", gator.HandlerLogin)

	if len(os.Args) < 2 {
		fmt.Print("Only two arguments?!")
		os.Exit(1)
	}

	login := gator.Command{
		Name:      os.Args[1],
		Arguments: os.Args[2:],
	}

	err = com.Run(state, login)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run command: %v\n", err)
		os.Exit(1)
	}
}
