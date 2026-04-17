package main

import (
	"fmt"
	"os"

	"github.com/Uami-11/blog-aggregator/internal/config"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read the configuration file: %v\n", err)
	}

	conf.CurrentUserName = "Uami"

	err = conf.SetUser()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set you as the current user: %v\n", err)
	}

	conf, err = config.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read the configuration file: %v\n", err)
	}

	fmt.Printf("database url: %s\ncurrent user: %s\n", conf.DBURL, conf.CurrentUserName)
}
