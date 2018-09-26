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

	// Declare flags common to commands, and pass them in Flags below.
	verFlag := cli.StringFlag{
		Name:  "symver",
		Value: "",
		Usage: "Symantic Version of the release to prepare",
	}

	dockerImage := cli.StringFlag{
		Name:  "dockerImage",
		Value: "",
		Usage: "Source Docker image to release",
	}

	dockerUserFlag := cli.StringFlag{
		Name:  "dockerUser",
		Value: "",
		Usage: "Username for docker push/pulls",
	}

	dockerPasswordFlag := cli.StringFlag{
		Name:  "dockerPassword",
		Value: "",
		Usage: "Password for docker push/pulls",
	}

	githubTokenFlag := cli.StringFlag{
		Name:  "githubToken",
		Value: "",
		Usage: "Access token for github releases",
	}

	githubUserFlag := cli.StringFlag{
		Name:  "githubUser",
		Value: "",
		Usage: "Username for github releases",
	}

	githubPasswordFlag := cli.StringFlag{
		Name:  "githubPassword",
		Value: "",
		Usage: "Password for github releases",
	}

	githubOrgFlag := cli.StringFlag{
		Name:  "githubOrg",
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
			Flags: []cli.Flag{verFlag, dockerImage, dockerUserFlag, dockerPasswordFlag},
			Action: func(clictx *cli.Context) error {
				if utils.IsDevCommit(clictx.String("symver")) {
					fmt.Println("Dev tag found; exiting")
					return nil
				}

				images, err := docker.PrepareDocker(clictx.String("dockerImage"), clictx.String("symver"))
				if err != nil {
					fmt.Println(images)
					fmt.Println(err)
				}

				for _, imageName := range images {
					docker.PushImage(imageName, clictx.String("dockerUser"), clictx.String("dockerPassword"))
				}

				return nil
			},
		},
		{
			Name:  "github",
			Usage: "Do the github release",
			Flags: []cli.Flag{verFlag, githubTokenFlag, githubOrgFlag, githubUserFlag, githubPasswordFlag, assetFlag},
			Action: func(clictx *cli.Context) error {
				if utils.IsDevCommit(clictx.String("symver")) {
					fmt.Println("Dev tag found; exiting")
					return nil
				}

				client := gh_client.NewClient(&http.Client{})
				if clictx.String("githubToken") != "" {
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
						Username: clictx.String("githubUser"),
						Password: clictx.String("githubPassword"),
					}
					client = gh_client.NewClient(&http.Client{Transport: &auth})
				}

				msg, err := github.PrepareGithubRelease(
					*client,
					clictx.String("symver"),
					clictx.String("githubOrg"),
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
		dockerUserFlag,
		dockerPasswordFlag,
		githubTokenFlag,
		githubUserFlag,
		githubPasswordFlag,
		assetFlag,
	}

	// There is no "default" command.  Print help and exit.
	app.Action = func(clictx *cli.Context) error {
		fmt.Printf("Must specify command. Run `%s help` for info\n", app.Name)
		return nil
	}

	app.Run(os.Args)
}
