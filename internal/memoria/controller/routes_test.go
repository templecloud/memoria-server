package controller

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)


type Endpoint struct {
	method string
	path string
}

func hasEndpoint(engine *gin.Engine, endpoint *Endpoint) bool {
	for _, route := range engine.Routes() {
		if route.Method == endpoint.method && route.Path == endpoint.path {
			return true
		}
	}
	return false
}

func TestControllerRouteConfiguration(t *testing.T) {
	// defaultConfig
	defaultConfig := newDefaultConfig()
	// disabledEndpointsConfig
	disabledEndpointsConfig := &Config{
		Endpoints: map[string]*EndpointConfig{
			HealthEndpoint: &EndpointConfig{IsDisabled: true},
			LoginEndpoint: &EndpointConfig{IsDisabled: true},
			SignupEndpoint: &EndpointConfig{IsDisabled: true},
		},
	}
	// disabledHealthEndpointConfig
	disabledHealthEndpointConfig := &Config{
		Endpoints: map[string]*EndpointConfig{
			HealthEndpoint: &EndpointConfig{IsDisabled: true},
		},
	}
	// Tests
	var tests = []struct {
		config   *Config
		expectedEnabled []Endpoint
		expectedDisabled []Endpoint
	}{
		{defaultConfig,
			[]Endpoint{
				{"GET", APIv1+HealthPath}, 
				{"POST", APIv1+LoginPath}, 
				{"POST", APIv1+SignupPath},
			},
			[]Endpoint{},
		},
		{disabledEndpointsConfig,
			[]Endpoint{},
			[]Endpoint{
				{"GET", APIv1+HealthPath}, 
				{"POST", APIv1+LoginPath}, 
				{"POST", APIv1+SignupPath},
			},
		},
		{disabledHealthEndpointConfig,
			[]Endpoint{{"POST", APIv1+LoginPath}, {"POST", APIv1+SignupPath},},
			[]Endpoint{{"GET", APIv1+HealthPath}},
		},
	}		
	// Run tests.
	for _, test := range tests {
		engine := gin.Default()
		engine = ConfigureEndpoints(engine, nil, test.config)
		for _, endpoint := range test.expectedEnabled {
			assert.True(t, hasEndpoint(engine, &endpoint))
		}
		for _, endpoint := range test.expectedDisabled {
			assert.False(t, hasEndpoint(engine, &endpoint))
		}
	}

}
