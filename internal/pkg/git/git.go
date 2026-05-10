package git

import (
	"context"
	"net/http"
	"strings"

	"github.com/K-Phoen/semver-release-action/internal/pkg/action"
	"github.com/blang/semver/v4"
	"github.com/google/go-github/v85/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

func LatestTagCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "latest-tag [REPOSITORY] [GH_TOKEN]",
		Args: cobra.ExactArgs(2),
		Run:  executeLatestTag,
	}
}

func executeLatestTag(cmd *cobra.Command, args []string) {
	repository := args[0]
	githubToken := args[1]

	ctx := context.Background()

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
	client := github.NewClient(oauth2.NewClient(ctx, tokenSource))

	parts := strings.Split(repository, "/")
	owner := parts[0]
	repo := parts[1]

	latestRelease, response, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if response != nil && response.StatusCode == http.StatusNotFound {
		cmd.Print("v0.0.0")
		return
	}
	action.AssertNoError(cmd, err, "could not list git refs: %s", err)

	latestTag, err := semver.ParseTolerant(*latestRelease.TagName)

	if err != nil {
		cmd.Print("v0.0.0")
		return
	}

	cmd.Printf("v%s", latestTag)
}
