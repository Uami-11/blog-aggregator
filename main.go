package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Uami-11/blog-aggregator/internal/config"
)

func main() {
	var userName string

	fmt.Print("Enter your username: ")

	scanner := bufio.NewScanner(os.Stdin)

	if scanner.Scan() {
		userName = scanner.Text()
	}

	currentConfig := config.Read()

	fmt.Println("Current:")
	fmt.Println(currentConfig.CurrentUserName)
	fmt.Println(currentConfig.DBURL)
	err := config.SetUser(userName)
	if err != nil {
		fmt.Printf("Write error: %v\n", err)
	}

	currentConfig = config.Read()
	fmt.Println("Now:")
	fmt.Println(currentConfig.CurrentUserName)
	fmt.Println(currentConfig.DBURL)
}
