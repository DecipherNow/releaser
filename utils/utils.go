package utils

import (
	"bytes"
	"os/exec"
	"path"
	"strings"
)

func IsDevCommit(version string) bool {
	return strings.HasSuffix(version, "dev")
}

// ParseSymver is a really quick parser for symantic versionin
func ParseSymver(version string) (major, minor, patch string, err error) {
	if strings.HasPrefix(version, "v") {
		version = version[1:len(version)]
	}
	splitSymver := strings.Split(version, ".")

	return splitSymver[0], splitSymver[1], splitSymver[2], nil
}

// ParseGitHash gets the commit hash from the local directory
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

// ParseRepoName gets the name of the repo
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
