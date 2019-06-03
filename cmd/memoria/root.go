package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"

	"github.com/templecloud/memoria-server/internal/memoria/boot"
)

//-------------------------------------------------------------------------------------------------
// Flags

var verbose bool
var config string

func init() {
	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "", 
		"Define the configuration file path.")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, 
		"Display verbose info during initialisation.")
}

//-------------------------------------------------------------------------------------------------
// Functions

func loadInitConfig() *boot.Config {
	var cfg boot.Config
	loader := newLoader(config)
	loader.Unmarshal(&cfg)
	return &cfg
}

// newLoader creates and loads the configuration specified by the config flag; or the default.
func newLoader(configFlag string) *viper.Viper {
	// Create Config - NB: 'ConfigName' does not require filetype extension.
	loader := viper.New()
	if config == "" {
		loader.SetConfigName("memoria")
		loader.AddConfigPath(".")
		loader.AddConfigPath("./config")
		loader.AddConfigPath("$MEMORIA_HOME/config")
		loader.AddConfigPath("$HOME/.memoria")
	} else {
		configPath, configName := handleConfigFlag(configFlag)
		loader.SetConfigName(configName)
		loader.AddConfigPath(configPath)
	}

	// Load Config
	if err := loader.ReadInConfig(); err != nil {
		fmt.Println("unable to read config:", err)
		os.Exit(1)
	}
	// Display Loaded config as YAML
	if verbose {
		fmt.Printf("loaded config >>>\n%s\n", toYAML(loader))
	}

	return loader
}

// handleConfigFlag processes the 'config' flag to return the directory and the extension stripped 
// filename.
func handleConfigFlag(configFlag string) (string, string) {
	return filepath.Split(strings.TrimSuffix(configFlag, filepath.Ext(configFlag)))
}

// toYAML converts the specified loaders config to a string.
func toYAML(loader *viper.Viper) string {
    c := loader.AllSettings()
	bs, err := yaml.Marshal(c)
	if err != nil {
        fmt.Println("unable to marshal config to YAML: ", err)
		os.Exit(1)
    }
	return string(bs)
}

