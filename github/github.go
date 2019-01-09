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
	"github.com/google/go-github/github"
	gh_client "github.com/google/go-github/github"
)

// UploadReleaseAsset takes a file and uploads it to a specific github release.
func UploadReleaseAsset(client gh_client.Client, releaseID int64, organization, repository, filename string) error {
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

// PrepareGithubRelease performs the entire github release process.
func PrepareGithubRelease(client gh_client.Client, semver, organization, asset string) (string, error) {
	repository, err := utils.ParseRepoName()

	gitHash, err := utils.ParseGitHash()

	major, minor, patch, err := utils.Parsesemver(semver)
	if err != nil {
		return "", err
	}

	name := fmt.Sprintf("Release %s", strings.Join([]string{major, minor, patch}, "."))
	fmt.Printf("Creating release %s\n", name)
	release := github.RepositoryRelease{
		TagName:         &semver,
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
			UploadReleaseAsset(client, *releaseResp.ID, organization, repository, filename)
		}
	}

	return "", nil
}
