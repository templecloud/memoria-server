package persistence

import (
	"fmt"
)

const (
	// Localhost is the localhost domain.
	Localhost = "localhost"
	// MongoDBProtocol is the MongoDB protocol.
	MongoDBProtocol = "mongodb"
	// MongoDBDefaultPort is the MongoDB default port.
	MongoDBDefaultPort = 27017
)

//-------------------------------------------------------------------------------------------------
// Models

// Config denotes the configuration of the server controller.
type Config struct {
	Connection *ConnectionConfig `json:"connection" binding:"optional"`
}

// ConnectionConfig denotes the configuration of the datastore.
type ConnectionConfig struct {
	Protocol string `json:"protocol" binding:"required"` 
	Host  string `json:"host" binding:"required"`
	Port int    `json:"port" binding:"required"`
}

//-------------------------------------------------------------------------------------------------
// Public Functions

// NewDefaultConfig returns a default MongoDB configuration.
func NewDefaultConfig() *Config {
	return &Config{
		Connection: &ConnectionConfig{
			Protocol: MongoDBProtocol,
			Host:  Localhost,
			Port: MongoDBDefaultPort,
		},
	}
}

// Get the datasource connection URI. e.g. "mongodb://localhost:27017"
func (c *Config) uri() string {
	c.validate()
	return fmt.Sprintf("%s://%s:%d", c.Connection.Protocol, c.Connection.Host, c.Connection.Port)
}

// Panic if the configuration is invalid.
func (c *Config) validate() {
	if c.Connection == nil {
		panic(fmt.Sprintf("The datasource connection was not defined."))
	}	
	if c.Connection.Protocol == "" {
		panic(fmt.Sprintf("The datasource connection protocol was not defined."))
	}
	if c.Connection.Host == "" {
		panic(fmt.Sprintf("The datasource connection host was not defined."))
	}
	if c.Connection.Port < 1000 {
		panic(fmt.Sprintf("The datasource connection port '%d' was less than 1000.", c.Connection.Port))
	}
}
