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
func FetchRepo(username string) ([]*github.Repository, error) {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "tocken"},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	//orgs, _, err := client.Organizations.List(context.Background(), username, nil)
	repos, _, err := client.Repositories.List(ctx, username, nil)
	return repos, err
}
func main() {
	mux := http.NewServeMux()
	repos, err := FetchRepo("irenetony")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Call html page handler.
	mux.HandleFunc("/", website)

	r := gin.Default()

	r.GET("/dataJSON", func(c *gin.Context) {

		var msg struct {
			Repo []string `json:"repos"`
		}
		for i := 0; i < len(repos); i++ {
			msg.Repo = append(msg.Repo, repos[i].GetName())
		}

		c.JSON(http.StatusOK, msg)
	})
	// Open port for web server.
	r.Run(":9000")
	port := ":8000"
	log.Println("Listening on port", port)
	http.ListenAndServe(port, mux)

}
