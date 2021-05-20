// Copyright 2017 Decipher Technology Studios LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package github provides high-level utilities for interacting with github
// for the purposes of performing releases and uploading assets.
package github

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/deciphernow/releaser/utils"
	gh_client "github.com/google/go-github/github"
)

// UploadReleaseAsset takes a file and uploads it to a specific github release.
func UploadReleaseAsset(client *gh_client.Client, releaseID int64, organization, repository, filename string) error {
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
		&gh_client.UploadOptions{Name: path.Base(filename)},
		file,
	)
	if err != nil {
		return err
	}

	return nil
}

// PrepareGithubRelease performs the entire github release process.
func PrepareGithubRelease(client *gh_client.Client, tag, organization, asset string) (string, error) {
	repository, err := utils.ParseRepoName()
	if err != nil {
		fmt.Println("failed to parse repo name:", err.Error())
		return "", err
	}

	gitHash, err := utils.ParseGitHash()
	if err != nil {
		fmt.Println("failed to parse git hash:", err.Error())
		return "", err
	}

	name := fmt.Sprintf("Release %s", tag)
	fmt.Printf("Creating release %s\n", name)
	release := gh_client.RepositoryRelease{
		TagName:         &tag,
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
		return "", err
	}

	if asset != "" {
		allAssets := strings.Split(asset, ",")

		for _, filename := range allAssets {
			UploadReleaseAsset(client, *releaseResp.ID, organization, repository, filename)
		}
	}

	return "", nil
}
