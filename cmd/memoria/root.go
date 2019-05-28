package main

import (
	"fmt"
	"os"

	// "github.com/spf13/cobra"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"

	"github.com/templecloud/memoria-server/internal/memoria/boot"
)

var verbose bool

func init() {
	// cobra.OnInitialize(initConfig, initLoggingConfig)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Display verbose info during initialisation.")
}

func loadInitConfig() *boot.Config {
	var cfg boot.Config
	loader := newLoader("memoria")
	loader.Unmarshal(&cfg)
	return &cfg
}

func newLoader(filename string) *viper.Viper {
	// Config name (without extension)
	loader := viper.New()
	loader.SetConfigName(filename)
	
	// Search paths...
	loader.AddConfigPath(".")
	loader.AddConfigPath("./config")
	loader.AddConfigPath("$MEMORIA_HOME/config")
	loader.AddConfigPath("$HOME/.memoria")

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

func toYAML(loader *viper.Viper) string {
    c := loader.AllSettings()
	bs, err := yaml.Marshal(c)
	if err != nil {
        fmt.Println("unable to marshal config to YAML: ", err)
		os.Exit(1)
    }
	return string(bs)
}

