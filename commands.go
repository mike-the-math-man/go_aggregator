package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mike-the-math-man/internal/config"
	"github.com/mike-the-math-man/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

type command struct {
	name      string
	arguments []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		fmt.Println("not enough arguments in handlerlogin")
		return fmt.Errorf("Nooooo")
	}
	user, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		fmt.Println("user does not exist - pls register")
		return err
	}
	s.config.SetUser(*s.config, user.Name)

	//s.config.Current_user_name = cmd.name
	fmt.Println("user set")
	return nil
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	err := c.commands[cmd.name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		fmt.Println("not enough arguments in register - supply name")
		return fmt.Errorf("Nooooo")
	}
	user_id := uuid.New()
	name := cmd.arguments[0]
	_, err := s.db.GetUser(context.Background(), name)
	params := database.CreateUserParams{ID: user_id, CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: name}
	if err != nil {
		user, err := s.db.CreateUser(context.Background(), params)
		if err != nil {
			fmt.Println("error creating user", err)
			return err
		}
		s.config.SetUser(*s.config, user.Name)
		fmt.Println("User created")
		return nil
	} else {
		return err
	}
}

func handlerReset(s *state, cmd command) error {
	err := s.db.TruncateUsers(context.Background())
	if err != nil {
		fmt.Println("error deleting users", err)
		return err
	}
	fmt.Println("Users Deleted")
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		fmt.Println("error sql users", err)
		return err
	}
	current_user := s.config.Current_user_name
	for _, john := range users {
		if john.Name == current_user {
			fmt.Printf("* %s (current)\n", john.Name)
		} else {
			fmt.Printf("* %s\n", john.Name)
		}
	}
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	request, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		fmt.Println("error creating request", err)
		return nil, err
	}
	request.Header.Set("User-Agent", "gator")
	var john http.Client
	result, err := john.Do(request)
	if err != nil {
		fmt.Println("error doing request", err)
		return nil, err
	}
	read_result, err := io.ReadAll(result.Body)
	if err != nil {
		fmt.Println("error reading result from request", err)
		return nil, err
	}
	var feed RSSFeed

	err = xml.Unmarshal(read_result, &feed)
	if err != nil {
		fmt.Println("error unmarshalling xml", err)
		return nil, err
	}
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}
	return &feed, nil
}

func aggregator_list(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		fmt.Println("not enough arguments in aggregator_list")
		return fmt.Errorf("Nooooo")
	}
	time_between_reqs, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		fmt.Println("error parsing time between requests", err)
		return err
	}
	fmt.Printf("Collecting feeds every %v\n", time_between_reqs)
	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
	return nil
}

func add_feed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) <= 1 {
		fmt.Println("not enough arguments in add_feed")
		return fmt.Errorf("Nooooo")
	}
	feed_id := uuid.New()
	name := cmd.arguments[0]
	url := cmd.arguments[1]
	//_, err = s.db.GetFeed(context.Background(), name)
	params := database.CreateFeedParams{ID: feed_id, CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: name, Url: url, UserID: user.ID}
	//if err != nil {
	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		fmt.Println("error creating feed", err)
		return err
	}
	fmt.Println("feed added")
	fmt.Printf("%+v\n", feed)
	var new_command []string
	new_command = append(new_command, cmd.arguments[1])
	cmd.arguments = new_command
	//fmt.Println(cmd)
	err = follow_feed(s, cmd, user)
	if err != nil {
		fmt.Println("problem following feed", err)
	}
	return nil
}

func get_feeds_list(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		fmt.Println("error sql feeds", err)
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("%+v\n", feed)
	}
	return nil
}

func follow_feed(s *state, cmd command, user database.User) error {
	follow_feed_id := uuid.New()
	url := cmd.arguments[0]

	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		fmt.Println("feed not followed, please add feed before following", err)
	}
	params := database.CreateFeedFollowParams{ID: follow_feed_id, CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: user.ID, FeedID: feed.ID}
	feed_follow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		fmt.Println("error following feed", err)
		return err
	}
	fmt.Printf("Feed: %s followed by user: %s\n", feed_follow.FeedName, feed_follow.UserName)
	return nil
}

func follow_feeds_list(s *state, cmd command, user database.User) error {
	feeds_list, err := s.db.GetFeedFollowsForUser(context.Background(), s.config.Current_user_name)
	if err != nil {
		fmt.Println("error sql get feed follows for user", err)
		return err
	}

	for _, feed := range feeds_list {
		fmt.Printf("%s\n", feed.FeedName)
	}
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		user, err := s.db.GetUser(context.Background(), s.config.Current_user_name)
		if err != nil {
			fmt.Println("error getting user", err)
			return err
		}
		return handler(s, c, user)
	}
}

func unfollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		fmt.Println("not enough arguments in unfollowFeed")
		return fmt.Errorf("Nooooo")
	}
	url := cmd.arguments[0]

	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		fmt.Println("feed not followed, please add feed before following", err)
	}
	params := database.UnfollowFeedParams{UserID: user.ID, FeedID: feed.ID}
	s.db.UnfollowFeed(context.Background(), params)
	return nil
}

func scrapeFeeds(s *state) {
	feed_row, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Println("error getting next feed to follow", err)
		return
	}
	err = s.db.MarkFeedFetched(context.Background(), feed_row.ID)
	if err != nil {
		fmt.Println("error marking feed fetched", err)
		return
	}
	feed, err := fetchFeed(context.Background(), feed_row.Url)
	if err != nil {
		fmt.Println("error fetching feed", err)
		return
	}
	for _, item := range feed.Channel.Item {
		Pubdate, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			fmt.Println("error parsing time", err)
		}
		params := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: Pubdate,
			FeedID:      feed_row.ID,
		}
		_, err = s.db.CreatePost(context.Background(), params)
		if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			fmt.Println(err)
		}
	}
}

func browse(s *state, cmd command, user database.User) error {
	limit := int32(2)
	if len(cmd.arguments) != 0 {
		i64, err := strconv.ParseInt(cmd.arguments[0], 10, 32)
		if err != nil {
			fmt.Println("error parsing input to int", err)
			return err
		}
		limit = int32(i64)
	}
	params := database.GetPostsForUserParams{UserID: user.ID, Limit: limit}
	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		fmt.Println("error getting posts for user")
		return err
	}
	for _, post := range posts {
		fmt.Printf("%s\n", post.Title)
		fmt.Printf("Feed: %s\n", post.FeedName)
		fmt.Printf("URL: %s\n", post.Url)
		re := regexp.MustCompile("<[^>]*>")
		clean := re.ReplaceAllString(post.Description, "")
		fmt.Printf("Description: %s\n", clean)
		fmt.Println()
	}
	return nil
}
