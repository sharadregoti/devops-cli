package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func imageCheck(imageName string) bool {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Search for the image by name
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{All: true})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Check if the image exists
	exists := false
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == imageName {
				exists = true
				break
			}
		}
	}

	if exists {
		return true
	} else {
		return false
	}
}

func imagePullCheck(imageName string) bool {
	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Check if the image can be pulled from the repository
	// repository := "docker.io"
	// authConfig := types.AuthConfig{}
	_, err = cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	if err != nil {
		return false
		// fmt.Println("Image cannot be pulled from the repository.")
		// fmt.Println(err)
		// os.Exit(1)
	}

	return true
}
