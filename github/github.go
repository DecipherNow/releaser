package github

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/deciphernow/releaser/utils"
	"github.com/google/go-github/github"
	gh_client "github.com/google/go-github/github"
)

//

// PrepareGithubRelease does stuff
func PrepareGithubRelease(client gh_client.Client, symver, organization, asset string) (string, error) {
	fmt.Println("github release process")
	repository, err := utils.ParseRepoName()

	gitHash, err := utils.ParseGitHash()
	fmt.Println(gitHash)

	major, minor, patch, err := utils.ParseSymver(symver)
	if err != nil {
		return "", err
	}

	name := fmt.Sprintf("Release %s", strings.Join([]string{major, minor, patch}, "."))
	fmt.Sprintf("Creating release %s", name)
	release := github.RepositoryRelease{
		TagName:         &symver,
		TargetCommitish: &gitHash,
		Name:            &name,
	}

	fmt.Println(release)

	releaseResp, resp, err := client.Repositories.CreateRelease(
		context.Background(),
		organization,
		repository,
		&release,
	)
	if err != nil {
		fmt.Sprintf("Error found in github release creation: %s", err)
	}
	fmt.Println(resp)

	file, err := os.Open(asset)
	if err != nil {
		return "Failed to open asset file", err
	}
	respAsset, response, err := client.Repositories.UploadReleaseAsset(
		context.Background(),
		organization,
		repository,
		*releaseResp.ID,
		&github.UploadOptions{Name: path.Base(asset)},
		file,
	)
	if err != nil {
		return "Failed to upload release asset", err
	}
	fmt.Println(respAsset)
	fmt.Println(response)

	return "", nil
}
