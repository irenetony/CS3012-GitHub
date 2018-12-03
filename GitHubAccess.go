package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

//FetchRepo gets all the organisations that the user speccified is a member of.
func FetchRepo(username string) (int, error) {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "c3073e3bf454522a59eb849bce3fb8d58bdaa457"},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(ctx, username, opt)
		if err != nil {
			return 0, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	numOfRepos := len(allRepos)
	return numOfRepos, nil
}

//FetchFollowing gets the users that the user specified follows.
func FetchFollowing(username string) ([]*github.User, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "c3073e3bf454522a59eb849bce3fb8d58bdaa457"},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	opt := &github.ListOptions{
		PerPage: 100,
	}
	//s := ""
	var users []*github.User
	for {
		user, resp, err := client.Users.ListFollowing(ctx, username, opt)
		if err != nil {
			return nil, err
		}
		users = append(users, user...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage

	}

	return users, nil
}
func main() {
	r := gin.Default()
	userName := ""

	r.LoadHTMLFiles("tpl/website.html")
	r.GET("/web", func(c *gin.Context) {
		c.HTML(http.StatusOK, "website.html", gin.H{"title": "Main website"})
	})
	r.POST("/post", func(c *gin.Context) {
		userName = c.PostForm("name")
		followingUsers, err := FetchFollowing(userName)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		var allRepos []int

		for i := 0; i < len(followingUsers); i++ {
			repoNum, err := FetchRepo(followingUsers[i].GetLogin())
			//fmt.Println(followingUsers[i].GetLogin())
			allRepos = append(allRepos, repoNum)

			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		}
		//fmt.Println(len(followingUsers) + len(allRepos))
		var msg struct {
			Repo []int    `json:"repos"`
			User []string `json:"user"`
		}
		for i := 0; i < len(allRepos); i++ {
			msg.Repo = append(msg.Repo, allRepos[i])
			msg.User = append(msg.User, followingUsers[i].GetLogin())
		}

		c.JSON(http.StatusOK, msg)
	})

	// Open port for web server.
	// languages the repos are in
	port := ":9000"
	r.Run(port)

}
