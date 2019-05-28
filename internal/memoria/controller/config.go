package controller

//-------------------------------------------------------------------------------------------------
// Models

// Config denotes the configuration of the server controller.
type Config struct {
	Endpoints map[string]*EndpointConfig `json:"endpoints" binding:"optional"`
}

// EndpointConfig defines the configurations of an Endpoint.
type EndpointConfig struct {
	// IsDisabled denotes it the Endpoint should be enabled. No value denotes false.
	IsDisabled bool `json:"isDisabled" binding:"required"`
}

//-------------------------------------------------------------------------------------------------
// Public

// NewDefaultConfig creates a default configuration.
func NewDefaultConfig() *Config {
	return &Config{
		Endpoints: map[string]*EndpointConfig{
			HealthEndpoint: &EndpointConfig{IsDisabled: false},
			LoginEndpoint:  &EndpointConfig{IsDisabled: false},
			SignupEndpoint: &EndpointConfig{IsDisabled: false},
		},
	}
}

//-------------------------------------------------------------------------------------------------
// Private

// isDisabled checks if an endpoint is disabled. If not configured then it is considered enabled.
func (c *Config) isDisabled(endpoint string) bool {
	if val, ok := c.Endpoints[endpoint]; ok {
		if val != nil {
			return val.IsDisabled
		}
	}
	return false
}

// isEnabled checks if an endpoint is enabled. If not configured then it is considered enabled.
func (c *Config) isEnabled(endpoint string) bool {
	return !c.isDisabled(endpoint)
}

