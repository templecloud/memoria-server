package persistence

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMongoDBConfigURI ensure the 'happy path' works as expected when generating a URI.
func TestMongoDBConfigURI(t *testing.T) {
	// defaultConfig
	defaultMongoDBConfig := NewDefaultConfig()
	// Config01
	mongoDBConfig01 := NewDefaultConfig()
	somehost01 := "somehost"
	mongoDBConfig01.Connection.Host = somehost01
	// Config02
	mongoDBConfig02 := NewDefaultConfig()
	somePort02 := 27018
	mongoDBConfig02.Connection.Port = somePort02
	// Config03
	mongoDBConfig03 := NewDefaultConfig()
	mongoDBConfig03.Connection.Host = somehost01
	mongoDBConfig03.Connection.Port = somePort02
	// Tests.
	var tests = []struct {
		config   *Config
		expected string
	}{
		{defaultMongoDBConfig, fmt.Sprintf("%s://%s:%d", MongoDBProtocol, Localhost, MongoDBDefaultPort)},
		{mongoDBConfig01, fmt.Sprintf("%s://%s:%d", MongoDBProtocol, somehost01, MongoDBDefaultPort)},
		{mongoDBConfig02, fmt.Sprintf("%s://%s:%d", MongoDBProtocol, Localhost, somePort02)},
		{mongoDBConfig03, fmt.Sprintf("%s://%s:%d", MongoDBProtocol, somehost01, somePort02)},
	}
	// Run tests.
	for _, test := range tests {
		assert.Equal(t, test.expected, test.config.uri())
	}
}

// TestConnectionConfigValidation test the config will correctly panic if validation fails.
func TestConnectionConfigValidation(t *testing.T) {
	// Tests.
	var tests = []struct {
		config   *Config
		expected string
	}{
		{	&Config{}, 
			"The datasource connection was not defined.",
		},
		{
			&Config{
				Connection: &ConnectionConfig{}, 
			},
			"The datasource connection protocol was not defined.",
		},
		{
			&Config{
				Connection: &ConnectionConfig{Protocol: MongoDBProtocol}, 
			},
			"The datasource connection host was not defined.",
		},
		{
			&Config{
				Connection: &ConnectionConfig{Protocol: MongoDBProtocol, Host: Localhost}, 
			},
			"The datasource connection port '0' was less than 1000.",
		},
		{
			&Config{
				Connection: &ConnectionConfig{Protocol: MongoDBProtocol, Host: Localhost, Port: -1}, 
			},
			"The datasource connection port '-1' was less than 1000.",
		},		
	}
	// Run tests.
	for _, test := range tests {
		f := func() {
    		test.config.validate()
		}
		assert.Panics(t, f)
		assert.PanicsWithValue(t, test.expected, f)
	}
}