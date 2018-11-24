package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/thedevsaddam/renderer"
	"golang.org/x/oauth2"
)

var rnd *renderer.Render

func init() {
	opts := renderer.Options{
		ParseGlobPattern: "tpl/*.html",
	}

	rnd = renderer.New(opts)
}

func website(w http.ResponseWriter, r *http.Request) {
	usr := struct {
		Name string
		ID   int
	}{
		Name: "John",
		ID:   001,
	}
	err := rnd.HTML(w, http.StatusOK, "website", usr)
	if err != nil {
		log.Fatal(err)
	}
}

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
	for i := 0; i < len(users); i++ {
		fmt.Printf("%v", users[i].GetLogin())
	}

	return users, err
}
func main() {
	mux := http.NewServeMux()
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

	// Call html page handler.
	mux.HandleFunc("/", website)

	r := gin.Default()

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
	r.Run(":9000")
	port := ":8000"
	log.Println("Listening on port", port)
	http.ListenAndServe(port, mux)

}
