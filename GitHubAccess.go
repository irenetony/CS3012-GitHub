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
		&oauth2.Token{AccessToken: "token"},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	repos, _, err := client.Repositories.List(ctx, username, nil)
	numOfRepos := len(repos)
	return numOfRepos, err
}

//FetchFollowing gets the users that the user specified follows.
func FetchFollowing(username string) ([]*github.User, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "token"},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	//s := ""
	users, _, err := client.Users.ListFollowing(ctx, username, nil)
	return users, err
}
func main() {

	followingUsers, err := FetchFollowing("irenetony")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	var allRepos []int

	for i := 0; i < len(followingUsers); i++ {
		repoNum, err := FetchRepo(followingUsers[i].GetLogin())
		allRepos = append(allRepos, repoNum)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	}

	r := gin.Default()

	r.LoadHTMLFiles("tpl/website.html")
	r.GET("/web", func(c *gin.Context) {
		c.HTML(http.StatusOK, "website.html", gin.H{"title": "Main website"})
	})
	r.GET("/dataJSON", func(c *gin.Context) {

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
