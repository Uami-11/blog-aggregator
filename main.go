package main

import (
	"fmt"
	"os"

	"github.com/Uami-11/blog-aggregator/cmd"
	"github.com/Uami-11/blog-aggregator/internal/config"
)

func main() {
	currentConfig := config.Read()

	var stateConfig cmd.State

	stateConfig.TheConfig = &currentConfig

	var comms cmd.Commands

	comms.Comms = make(map[string]func(*cmd.State, cmd.Command) error)

	comms.Register("login", cmd.HandlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("no command given")
		fmt.Println("")
		os.Exit(1)
	}

	var comm cmd.Command

	comm.Name = os.Args[1]
	comm.Args = os.Args[2:]

	fmt.Println("Current:")
	fmt.Println(stateConfig.TheConfig.CurrentUserName)
	fmt.Println(stateConfig.TheConfig.DBURL)

	err := comms.Run(&stateConfig, comm)
	if err != nil {
		fmt.Printf("command error: %v", err)
		fmt.Println("")
		os.Exit(1)
	}

	currentConfig = config.Read()
	stateConfig.TheConfig = &currentConfig

	fmt.Println("Now:")
	fmt.Println(stateConfig.TheConfig.CurrentUserName)
	fmt.Println(stateConfig.TheConfig.DBURL)
}
