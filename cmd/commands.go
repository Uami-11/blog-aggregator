// Package cmd...
package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Uami-11/blog-aggregator/article"
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

func HandlerReset(s *State, comm Command) error {
	err := s.DB.DeleteUsers(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func HandlerAgg(s *State, comm Command) error {
	feed, err := article.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println("Title: ", feed.Channel.Title)
	fmt.Println("Link: ", feed.Channel.Link)
	fmt.Println("Description: ", feed.Channel.Description)
	fmt.Println("Items:")
	for _, item := range feed.Channel.Item {
		fmt.Println(item)
	}

	return nil
}

func HandlerUsers(s *State, comm Command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		fmt.Printf("• %s", user)
		if user == s.TheConfig.CurrentUserName {
			fmt.Print(" (current)")
		}

		fmt.Println()
	}

	return nil
}

func HandlerAddFeed(s *State, comm Command) error {
	if len(comm.Args) < 2 {
		return errors.New("feed needs two arguments: name and url")
	}

	var feeds database.CreateFeedParams

	feeds.ID = uuid.New()
	feeds.Name = comm.Args[0]
	feeds.Url = comm.Args[1]

	currentTime := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	feeds.CreatedAt = currentTime
	feeds.UpdatedAt = currentTime

	user, err := s.DB.GetUser(context.Background(), s.TheConfig.CurrentUserName)
	if err != nil {
		return errors.New("you are not logged in, can not create feed")
	}

	feeds.UserID = user.ID

	_, err = s.DB.CreateFeed(context.Background(), feeds)
	if err != nil {
		return errors.New("error in creating feed")
	}

	feedFollow := createFeedFollow(feeds.ID, user.ID)

	_, err = s.DB.CreateFeedFollow(context.Background(), feedFollow)
	if err != nil {
		return err
	}

	return nil
}

func HandlerFeeds(s *State, comm Command) error {
	feeds, err := s.DB.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		user, err := s.DB.GetFeedUser(context.Background(), feed.ID)
		if err != nil {
			return err
		}

		fmt.Printf("Name: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("By: %s\n", user)
	}

	return nil
}

func createFeedFollow(feedID, userID uuid.UUID) database.CreateFeedFollowParams {
	var feedFollow database.CreateFeedFollowParams

	feedFollow.ID = uuid.New()

	currentTime := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	feedFollow.CreatedAt = currentTime
	feedFollow.UpdatedAt = currentTime
	feedFollow.FeedID = feedID
	feedFollow.UserID = userID

	return feedFollow
}

func HandlerFollow(s *State, comm Command) error {
	if len(comm.Args) == 0 {
		return errors.New("follow command needs a url argument")
	}

	feedID, err := s.DB.FindFeedURL(context.Background(), comm.Args[0])
	if err != nil {
		return err
	}

	user, err := s.DB.GetUser(context.Background(), s.TheConfig.CurrentUserName)
	if err != nil {
		return err
	}

	feedFollow := createFeedFollow(feedID, user.ID)

	feedFollowInfo, err := s.DB.CreateFeedFollow(context.Background(), feedFollow)
	if err != nil {
		return err
	}

	fmt.Printf("%s followed the %s feed\n", feedFollowInfo[0].UserName, feedFollowInfo[0].FeedName)

	return nil
}

func HandlerFollowing(s *State, comm Command) error {
	user, err := s.DB.GetUser(context.Background(), s.TheConfig.CurrentUserName)
	if err != nil {
		return err
	}

	feedsFollowed, err := s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Printf("%s follows these feeds:\n", user.Name)
	for _, feed := range feedsFollowed {
		fmt.Println(feed.FeedName)
	}

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
	err = s.TheConfig.SetUser(user.Name)
	if err != nil {
		return errors.New("error in setting user")
	}

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
