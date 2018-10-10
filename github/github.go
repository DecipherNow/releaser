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

// UploadReleaseAsset
func uploadReleaseAsset(client gh_client.Client, releaseID int64, organization, repository, filename string) error {
	fmt.Printf("Uploading %s to %s/%s at release ID %d\n", filename, organization, repository, releaseID)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	_, _, err = client.Repositories.UploadReleaseAsset(
		context.Background(),
		organization,
		repository,
		releaseID,
		&github.UploadOptions{Name: path.Base(filename)},
		file,
	)
	if err != nil {
		return err
	}

	return nil
}

// PrepareGithubRelease does stuff
func PrepareGithubRelease(client gh_client.Client, symver, organization, asset string) (string, error) {
	repository, err := utils.ParseRepoName()

	gitHash, err := utils.ParseGitHash()

	major, minor, patch, err := utils.ParseSymver(symver)
	if err != nil {
		return "", err
	}

	name := fmt.Sprintf("Release %s", strings.Join([]string{major, minor, patch}, "."))
	fmt.Printf("Creating release %s\n", name)
	release := github.RepositoryRelease{
		TagName:         &symver,
		TargetCommitish: &gitHash,
		Name:            &name,
	}

	releaseResp, _, err := client.Repositories.CreateRelease(
		context.Background(),
		organization,
		repository,
		&release,
	)
	if err != nil {
		fmt.Printf("Error found in github release creation: %s", err)
	}

	if asset != "" {
		allAssets := strings.Split(asset, ",")

		for _, filename := range allAssets {
			uploadReleaseAsset(client, *releaseResp.ID, organization, repository, filename)
		}
	}

	return "", nil
}
