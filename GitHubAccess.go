package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

//FetchOrganizations gets all the organisations that the user speccified is a member of.
func FetchOrganizations(username string) ([]*github.Repository, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "token"},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	//orgs, _, err := client.Organizations.List(context.Background(), username, nil)
	repos, _, err := client.Repositories.List(ctx, username, nil)
	return repos, err
}
func main() {

	organizations, err := FetchOrganizations("irenetony")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	for i, organization := range organizations {
		fmt.Printf("%v. %v\n", i+1, organization.GetName())
	}
}
