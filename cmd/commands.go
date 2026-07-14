// Package cmd...
package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
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

func HandlerBrowse(s *State, comm Command) error {
	var limit int32

	if len(comm.Args) != 0 {
		limitInt, err := strconv.Atoi(comm.Args[0])
		if err != nil {
			return err
		}
		limit = int32(limitInt)
	} else {
		limit = 2
	}

	var param database.GetPostsForUserParams

	user, _ := s.DB.GetUser(context.Background(), s.TheConfig.CurrentUserName)

	param.ID = user.ID

	param.Limit = limit
	posts, err := s.DB.GetPostsForUser(context.Background(), param)
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Printf("Description:\n%s\n", post.Description)
	}

	return nil
}

func ScrapeFeeds(s *State) error {
	feed, err := s.DB.GetNextFeedToFetched(context.Background())
	if err != nil {
		return err
	}

	currentTime := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	var markParams database.MarkFeedFetchedParams
	markParams.ID = feed.ID
	markParams.LastFetchedAt = currentTime

	err = s.DB.MarkFeedFetched(context.Background(), markParams)
	if err != nil {
		return err
	}

	theFeed, err := article.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	var postParams database.CreatePostParams

	fmt.Println("Title: ", theFeed.Channel.Title)
	fmt.Println("Link: ", theFeed.Channel.Link)
	fmt.Println("Description: ", theFeed.Channel.Description)
	for _, item := range theFeed.Channel.Item {
		postParams.Title = item.Title
		postParams.Url = item.Link
		postParams.Description = item.Description
		pubTime, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			return err
		}
		postParams.PublishedAt = sql.NullTime{
			Time:  pubTime,
			Valid: true,
		}

	}

	currentTime = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	postParams.ID = uuid.New()
	postParams.CreatedAt = currentTime
	postParams.UpdatedAt = currentTime
	postParams.FeedID = feed.ID

	_, err = s.DB.CreatePost(context.Background(), postParams)
	if err != nil {
		return err
	}

	return nil
}

func HandlerAgg(s *State, comm Command) error {
	if len(comm.Args) == 0 {
		return errors.New("agg command needs a time argument")
	}

	timeBetweenReqs, err := time.ParseDuration(comm.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v", timeBetweenReqs)
	ticker := time.NewTicker(timeBetweenReqs)

	defer ticker.Stop()

	for ; ; <-ticker.C {
		err = ScrapeFeeds(s)
		fmt.Println("Scraped, on to the next")
		if err != nil {
			return err
		}
	}
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

// Middleware add

func MiddlewareLoggedIn(handler func(s *State, comm Command) error) func(*State, Command) error {
	return func(s *State, comm Command) error {
		_, err := s.DB.GetUser(context.Background(), s.TheConfig.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, comm)
	}
}

func HandlerAddFeed(s *State, comm Command) error {
	if len(comm.Args) < 2 {
		return errors.New("feed needs two arguments: name and url")
	}

	user, _ := s.DB.GetUser(context.Background(), s.TheConfig.CurrentUserName)

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

	feeds.UserID = user.ID

	_, err := s.DB.CreateFeed(context.Background(), feeds)
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

func HandlerUnfollow(s *State, comm Command) error {
	if len(comm.Args) == 0 {
		return errors.New("follow command needs a url argument")
	}

	user, _ := s.DB.GetUser(context.Background(), s.TheConfig.CurrentUserName)

	feedID, err := s.DB.FindFeedURL(context.Background(), comm.Args[0])
	if err != nil {
		return err
	}

	var unfollow database.UnfollowFeedParams

	unfollow.FeedID = feedID
	unfollow.UserID = user.ID

	err = s.DB.UnfollowFeed(context.Background(), unfollow)
	if err != nil {
		return err
	}

	fmt.Println("Successfully unfollowed")

	return nil
}

func HandlerFollow(s *State, comm Command) error {
	if len(comm.Args) == 0 {
		return errors.New("follow command needs a url argument")
	}

	user, _ := s.DB.GetUser(context.Background(), s.TheConfig.CurrentUserName)

	feedID, err := s.DB.FindFeedURL(context.Background(), comm.Args[0])
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
	user, _ := s.DB.GetUser(context.Background(), s.TheConfig.CurrentUserName)

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
