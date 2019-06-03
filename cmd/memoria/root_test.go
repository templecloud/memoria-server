package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHandleConfigFlag tests handling a 'configFlag' string.
func TestHandleConfigFlag(t *testing.T) {
	// Tests.
	var tests = []struct {
		config   string
		expectedPath string
		expectedName string
	}{
		{"/some/path/name.yaml", "/some/path/", "name"},
		{"/some/path/name.toml", "/some/path/", "name"},
		{"/some/path/name.json", "/some/path/", "name"},
		{"/some/path/name.other.yaml", "/some/path/", "name.other"},
		{"/some/path/name.", "/some/path/", "name"},
		{"/some/path/name", "/some/path/", "name"},

		{"./some/path/name.yaml", "./some/path/", "name"},
		{"./some/path/name.toml", "./some/path/", "name"},
		{"./some/path/name.json", "./some/path/", "name"},
		{"./some/path/name.other.yaml", "./some/path/", "name.other"},
		{"./some/path/name.", "./some/path/", "name"},
		{"./some/path/name", "./some/path/", "name"},

		{"../some/path/name.yaml", "../some/path/", "name"},
		{"../some/path/name.toml", "../some/path/", "name"},
		{"../some/path/name.json", "../some/path/", "name"},
		{"../some/path/name.other.yaml", "../some/path/", "name.other"},
		{"../some/path/name.", "../some/path/", "name"},
		{"../some/path/name", "../some/path/", "name"},

		{"name.yaml", "", "name"},
		{"name.toml", "", "name"},
		{"name.json", "", "name"},
		{"name.other.yaml", "", "name.other"},
		{"name.", "", "name"},
		{"name", "", "name"},

	}
	// Run tests.
	for _, test := range tests {
		path, name := handleConfigFlag(test.config)
		assert.Equal(t, test.expectedPath, path)
		assert.Equal(t, test.expectedName, name)
	}
}