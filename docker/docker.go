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

// Package docker provides high-level utilities for interacting with the local
// docker daemon and pushing images to remote docker registries
package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/deciphernow/releaser/utils"
	"github.com/docker/docker/api/types"
	dc "github.com/docker/docker/client"
)

// PrepareDocker tags the source image as the final release images.
func PrepareDocker(source string, symver string, suffix string) ([]string, error) {
	source_base := strings.Split(source, ":")[0]

	major, minor, patch, err := utils.ParseSymver(symver)
	if err != nil {
		return []string{}, err
	}

	allTags := []string{
		strings.Join([]string{major, suffix}, ""),
		strings.Join([]string{strings.Join([]string{major, minor}, "."), suffix}, ""),
		strings.Join([]string{strings.Join([]string{major, minor, patch}, "."), suffix}, ""),
	}

	madeImages := []string{}
	for _, tag := range allTags {
		destination := strings.Join([]string{source_base, tag}, ":")
		_, err := tagImage(source, destination)
		if err != nil {
			fmt.Printf("Error tagging image as %s\n", destination)
		} else {
			madeImages = append(madeImages, destination)
		}
	}

	return madeImages, nil
}

// PushImage pushes docker images to a docker registry
func PushImage(source string, user string, password string) error {
	dockerCli, err := dc.NewEnvClient()
	if err != nil {
		return err
	}
	fmt.Println("Pushing next", source)

	authConfig := types.AuthConfig{
		Username: user,
		Password: password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	closer, err := dockerCli.ImagePush(context.Background(), source, types.ImagePushOptions{RegistryAuth: authStr})

	defer closer.Close()

	_, err = io.Copy(os.Stdout, closer)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	return nil
}

func tagImage(source string, destination string) (string, error) {
	dockerCli, err := dc.NewEnvClient()
	if err != nil {
		return "Could not establish docker client", err
	}

	err = dockerCli.ImageTag(context.Background(), source, destination)
	if err != nil {
		return fmt.Sprintf("Error tagging image as %s\n", destination), err
	}

	return "", nil
}
