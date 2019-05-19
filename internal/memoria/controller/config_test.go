package controller

import (
	// "encoding/json"
	// "net/http"
	// "net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	// defaultConfig
	defaultConfig := newDefaultConfig()
	// emptyConfig
	emptyConfig := &Config{}
	// emptyEndpointsConfig
	emptyEndpointsConfig := &Config{
		Endpoints: map[string]*EndpointConfig{},
	}
	// nilEndpointsConfig
	nilEndpointsConfig := &Config{
		Endpoints: map[string]*EndpointConfig{
			"NilEndpointConfigEndPoint01": nil,
			"NilEndpointConfigEndPoint02": nil,
			"NilEndpointConfigEndPoint03": nil,
		},
	}
	// nilIsDisabledEndpointsConfig
	nilIsDisabledEndpointsConfig := &Config{
		Endpoints: map[string]*EndpointConfig{
			"NoIsDisabledEndpointConfigEndPoint01": &EndpointConfig{},
			"NoIsDisabledEndpointConfigEndPoint02": &EndpointConfig{},
			"NoIsDisabledEndpointConfigEndPoint03": &EndpointConfig{},
		},
	}
	// disabledEndpointsConfig
	disabledEndpointsConfig := &Config{
		Endpoints: map[string]*EndpointConfig{
			"DisabledEndpoint01": &EndpointConfig{IsDisabled: true},
			"DisabledEndpoint02": &EndpointConfig{IsDisabled: true},
			"DisabledEndpoint03": &EndpointConfig{IsDisabled: true},
		},
	}
	// mixedEndpointsConfig
	mixedEndpointsConfig := &Config{
		Endpoints: map[string]*EndpointConfig{
			"DisabledEndpoint01":                   &EndpointConfig{IsDisabled: true},
			"NonDisabledEndpoint01":                &EndpointConfig{IsDisabled: false},
			"NilEndpointConfigEndPoint01":          nil,
			"NoIsDisabledEndpointConfigEndPoint01": &EndpointConfig{},
		},
	}
	// Tests.
	var tests = []struct {
		endpoint string
		config   *Config
		expected bool
	}{
		// defaultConfig
		{HealthEndpoint, defaultConfig, true},
		{LoginEndpoint, defaultConfig, true},
		{SignupEndpoint, defaultConfig, true},
		// emptyConfig
		{"UnConfiguredEndpoint01", emptyConfig, true},
		{"UnConfiguredEndpoint02", emptyConfig, true},
		{"UnConfiguredEndpoint03", emptyConfig, true},
		// emptyEndpointsConfig
		{"UnConfiguredEndpoint01", emptyEndpointsConfig, true},
		{"UnConfiguredEndpoint02", emptyEndpointsConfig, true},
		{"UnConfiguredEndpoint03", emptyEndpointsConfig, true},
		// nilEndpointsConfig
		{"NilEndpointConfigEndPoint01", nilEndpointsConfig, true},
		{"NilEndpointConfigEndPoint02", nilEndpointsConfig, true},
		{"NilEndpointConfigEndPoint02", nilEndpointsConfig, true},
		// nilIsDisabledEndpointsConfig
		{"NoIsDisabledEndpointConfigEndPoint01", nilIsDisabledEndpointsConfig, true},
		{"NoIsDisabledEndpointConfigEndPoint01", nilIsDisabledEndpointsConfig, true},
		{"NoIsDisabledEndpointConfigEndPoint01", nilIsDisabledEndpointsConfig, true},
		// Full disabled Config.
		{"DisabledEndpoint01", disabledEndpointsConfig, false},
		{"DisabledEndpoint02", disabledEndpointsConfig, false},
		{"DisabledEndpoint03", disabledEndpointsConfig, false},
		// // Mixed Config.
		{"DisabledEndpoint01", mixedEndpointsConfig, false},
		{"NonDisabledEndpoint01", mixedEndpointsConfig, true},
		{"NoEndpointConfigEndPoint01", mixedEndpointsConfig, true},
		{"NilEndpointConfigEndPoint01", mixedEndpointsConfig, true},
		{"NoIsDisabledEndpointConfigEndPoint01", mixedEndpointsConfig, true},
	}
	// Run tests.
	for _, test := range tests {
		assert.Equal(t, test.expected, test.config.isEnabled(test.endpoint))
		assert.Equal(t, test.expected, !test.config.isDisabled(test.endpoint))
	}
}
