package logging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDefaultLoggerConfigValidation test the config will correctly panic if validation fails.
func TestDefaultLoggerConfigValidation(t *testing.T) {
	// Tests.
	var tests = []struct {
		config   *Config
		expected string
	}{
		{	&Config{}, 
			"No default logger was defined.",
		},
		{
			&Config{
				DefaultLogger: &LogConfig{}, 
			},
			"The type of the default logger was not defined.",
		},
		{
			&Config{
				DefaultLogger: &LogConfig{Type: ""}, 
			},
			"The type of the default logger was not defined.",
		},
		{
			&Config{
				DefaultLogger: &LogConfig{Type: "unsupported"}, 
			},
			"The type of the default logger 'unsupported' was invalid.",
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