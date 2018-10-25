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

// Package utils provides simple general-purpose tooling that has no better
// place in this package.
package utils

import (
	"bytes"
	"os/exec"
	"path"
	"strings"
)

// IsDevCommit determins if the version string represents a development commit.
//
// The check is simply against the suffix of `dev`.  If the version string
// contains the suffix `dev`, this function will return true. e.g.
// 1.0.0-dev
// test.dev
// 4.x.x.3-dev
// will all be considered dev tags.  Any other string will be considered a release.
func IsDevCommit(version string) bool {
	return strings.HasSuffix(version, "dev")
}

// ParseSymver is a really quick parser for symantic versioning.
func ParseSymver(version string) (major, minor, patch string, err error) {
	if strings.HasPrefix(version, "v") {
		version = version[1:len(version)]
	}
	splitSymver := strings.Split(version, ".")

	return splitSymver[0], splitSymver[1], splitSymver[2], nil
}

// ParseGitHash gets the commit hash from the git repo in the current directory.
func ParseGitHash() (hash string, err error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()

	if err != nil {
		return "", err
	}

	gitHash := out.String()

	return strings.TrimSpace(gitHash), nil
}

// ParseRepoName gets the name of the repo from the current directory.
func ParseRepoName() (hash string, err error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()

	if err != nil {
		return "", err
	}

	fullPath := out.String()
	repoName := path.Base(fullPath)

	return strings.TrimSpace(repoName), nil
}
