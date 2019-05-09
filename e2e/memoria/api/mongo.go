package identity

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// CreateNewContainer creates a new container instance from the image.
func CreateNewContainer(image string) (string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}

	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: "27017",
	}
	containerPort, err := nat.NewPort("tcp", "27017")
	if err != nil {
		panic("Unable to get the port")
	}

	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}
	cont, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: image,
		},
		&container.HostConfig{
			PortBindings: portBinding,
		}, nil, "")
	if err != nil {
		panic(err)
	}

	cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	fmt.Printf("Container %s is started\n", cont.ID)
	return cont.ID, nil
}

// StopContainer stops the specified container.
func StopContainer(containerID string) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	
	err = cli.ContainerStop(context.Background(), containerID, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Container %s is stopped\n", containerID)
	return err
}

// RemoveContainer removes the specified container.
func RemoveContainer(containerID string) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	
	opts := types.ContainerRemoveOptions{RemoveVolumes: true, RemoveLinks: false, Force: true}
	err = cli.ContainerRemove(context.Background(), containerID, opts)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Container %s is removed\n", containerID)
	return err
}
