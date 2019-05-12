package framework

import (
	"fmt"
	
	"github.com/docker/docker/client"
)

// Framework is used in the execution of e2e tests.
type Framework struct {
	Name string
	client *client.Client
	mongoDB *MongoDB

}

// NewFramework constructs a new e2e test Framework with default options.
func NewFramework(name string) *Framework {
	client, err := client.NewEnvClient()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Docker daemon: %v", err))
	}

	mongoDB := NewMongoDB(client)

	return &Framework{
		Name: name,
		client: client,
		mongoDB: mongoDB,
	}
}
// BeforeEach creates resources before a test.
func (fw *Framework) BeforeEach() {
	fw.mongoDB.Create()
}

// AfterEach cleans up resources after a test.
func (fw *Framework) AfterEach() {
	fw.mongoDB.Stop()
	fw.mongoDB.Clean()
}
