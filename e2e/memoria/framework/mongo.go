package framework

import (
	"context"
	"fmt"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const mongoContainerImage = "mongo:4.1.9-bionic"
const mongoDefaultPort = 27017

// MongoDB defines a containerized MongoDB instance.
type MongoDB struct {
	client *client.Client
	containerID string
	Image string
	HostPort int
	ContainerPort int
}

// NewMongoDB creates a manageable containerized MongoDB instance.
func NewMongoDB() *MongoDB {
	client, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Unable to create docker client.")
		panic(err)
	}
	return &MongoDB{
		client: client,
		Image: mongoContainerImage,
		HostPort: mongoDefaultPort,
		ContainerPort: mongoDefaultPort,
	}
}

// Create a new MongoDB instance.
func (mdb *MongoDB) Create() (*MongoDB, error) {
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: strconv.Itoa(mdb.HostPort),
	}
	containerPort, err := nat.NewPort("tcp", strconv.Itoa(mdb.ContainerPort))
	if err != nil {
		panic("Unable to get the containerPort.")
	}

	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}
	cont, err := mdb.client.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: mdb.Image,
		},
		&container.HostConfig{
			PortBindings: portBinding,
		}, nil, "")
	if err != nil {
		panic(err)
	}

	mdb.client.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	mdb.containerID = cont.ID
	fmt.Printf("MongoDB container %s is started.\n", mdb.containerID)
	// return mdb.containerID, nil
	return mdb, nil
}

// Stop the specified MongoDB.
func (mdb *MongoDB) Stop() error {
	err := mdb.client.ContainerStop(context.Background(), mdb.containerID, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("MongoDB container %s is stopped.\n", mdb.containerID)
	return err
}

// Clean removes the specified MongoDB container resources.
func (mdb *MongoDB) Clean() error {
	opts := types.ContainerRemoveOptions{RemoveVolumes: true, RemoveLinks: false, Force: true}
	err := mdb.client.ContainerRemove(context.Background(), mdb.containerID, opts)
	if err != nil {
		panic(err)
	}
	fmt.Printf("MongoDB container %s is removed.\n", mdb.containerID)
	mdb.containerID = ""
	return err
}
