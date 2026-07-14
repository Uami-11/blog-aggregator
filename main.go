package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Uami-11/blog-aggregator/cmd"
	"github.com/Uami-11/blog-aggregator/internal/config"
	"github.com/Uami-11/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	currentConfig := config.Read()

	var stateConfig cmd.State

	stateConfig.TheConfig = &currentConfig

	db, err := sql.Open("postgres", stateConfig.TheConfig.DBURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "db connection: %v", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	stateConfig.DB = dbQueries

	var comms cmd.Commands

	comms.Comms = make(map[string]func(*cmd.State, cmd.Command) error)

	comms.Register("login", cmd.HandlerLogin)
	comms.Register("register", cmd.HandlerRegister)
	comms.Register("reset", cmd.HandlerReset)
	comms.Register("users", cmd.HandlerUsers)
	comms.Register("agg", cmd.HandlerAgg)
	comms.Register("addfeed", cmd.MiddlewareLoggedIn(cmd.HandlerAddFeed))
	comms.Register("feeds", cmd.HandlerFeeds)
	comms.Register("follow", cmd.MiddlewareLoggedIn(cmd.HandlerFollow))
	comms.Register("following", cmd.MiddlewareLoggedIn(cmd.HandlerFollowing))
	comms.Register("unfollow", cmd.MiddlewareLoggedIn(cmd.HandlerUnfollow))
	comms.Register("browse", cmd.MiddlewareLoggedIn(cmd.HandlerBrowse))

	if len(os.Args) < 2 {
		fmt.Println("no command given")
		fmt.Println("")
		os.Exit(1)
	}

	var comm cmd.Command

	comm.Name = os.Args[1]
	comm.Args = os.Args[2:]

	fmt.Printf("Running command %s:\n", comm.Name)
	err = comms.Run(&stateConfig, comm)
	if err != nil {
		fmt.Printf("command error: %v", err)
		fmt.Println("")
		os.Exit(1)
	}

	currentConfig = config.Read()
	stateConfig.TheConfig = &currentConfig
}
