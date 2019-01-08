package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/deciphernow/releaser/docker"
	"github.com/deciphernow/releaser/github"
	"github.com/deciphernow/releaser/utils"
	gh_client "github.com/google/go-github/github"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
)

func main() {
	app := cli.NewApp()
	app.Name = "Releaser"
	app.Usage = "Facilitate the release process of artifacts"
	app.Version = "0.2.2"

	// Declare flags common to commands, and pass them in Flags below.
	verFlag := cli.StringFlag{
		Name:  "semver",
		Value: "",
		Usage: "Semantic Version of the release to prepare",
	}

	dockerImage := cli.StringFlag{
		Name:  "image",
		Value: "",
		Usage: "Source Docker image to release",
	}

	dockerSuffix := cli.StringFlag{
		Name:  "suffix",
		Value: "",
		Usage: "String to be appended to the final docker tag. e.g. -alpine, -centos",
	}

	usernameFlag := cli.StringFlag{
		Name:  "username",
		Value: "",
		Usage: "Username for cmd operations",
	}

	passwordFlag := cli.StringFlag{
		Name:  "password",
		Value: "",
		Usage: "Password for cmd operations",
	}

	githubTokenFlag := cli.StringFlag{
		Name:  "token",
		Value: "",
		Usage: "Access token for github releases",
	}

	githubOrgFlag := cli.StringFlag{
		Name:  "organization",
		Value: "",
		Usage: "Organization for github releases",
	}

	assetFlag := cli.StringFlag{
		Name:  "asset",
		Value: "",
		Usage: "File[s] to be uploaded to the github release",
	}

	app.Commands = []cli.Command{
		{
			Name:  "docker",
			Usage: "Do the docker job",
			Flags: []cli.Flag{verFlag, dockerImage, usernameFlag, passwordFlag, dockerSuffix},
			Action: func(clictx *cli.Context) error {
				if utils.IsDevCommit(clictx.String("semver")) {
					fmt.Println("Dev tag found; exiting")
					return nil
				}

				images, err := docker.PrepareDocker(clictx.String("image"), clictx.String("semver"), clictx.String("suffix"))
				if err != nil {
					fmt.Println(images)
					fmt.Println(err)
				}

				for _, imageName := range images {
					docker.PushImage(imageName, clictx.String("username"), clictx.String("password"))
				}

				return nil
			},
		},
		{
			Name:  "github",
			Usage: "Do the github release",
			Flags: []cli.Flag{verFlag, githubTokenFlag, githubOrgFlag, usernameFlag, passwordFlag, assetFlag},
			Action: func(clictx *cli.Context) error {
				if utils.IsDevCommit(clictx.String("semver")) {
					fmt.Println("Dev tag found; exiting")
					return nil
				}

				client := gh_client.NewClient(&http.Client{})
				if clictx.String("token") != "" {
					fmt.Println("Using token auth")
					ctx := context.Background()
					ts := oauth2.StaticTokenSource(
						&oauth2.Token{AccessToken: clictx.String("githubToken")},
					)
					tc := oauth2.NewClient(ctx, ts)
					client = gh_client.NewClient(tc)
				} else {
					fmt.Println("Using Username/Password auth")
					auth := gh_client.BasicAuthTransport{
						Username: clictx.String("username"),
						Password: clictx.String("password"),
					}
					client = gh_client.NewClient(&http.Client{Transport: &auth})
				}

				msg, err := github.PrepareGithubRelease(
					*client,
					clictx.String("semver"),
					clictx.String("organization"),
					clictx.String("asset"),
				)

				if err != nil {
					fmt.Println("Error in github stuff", err)
					fmt.Println(msg)
				}

				return nil
			},
		},
	}

	// Global flags. Used when no "command" passed. Must be repeated above for commands.

	app.Flags = []cli.Flag{
		verFlag,
		dockerImage,
		usernameFlag,
		passwordFlag,
		githubTokenFlag,
		assetFlag,
	}

	// There is no "default" command.  Print help and exit.
	app.Action = func(clictx *cli.Context) error {
		fmt.Printf("Must specify command. Run `%s help` for info\n", app.Name)
		return nil
	}

	app.Run(os.Args)
}
