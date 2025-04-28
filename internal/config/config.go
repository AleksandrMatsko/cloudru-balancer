// Defines the configuration for server
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Read given configFileName and fill config.
func Read(configFileName string, config interface{}) error {
	bytes, err := os.ReadFile(configFileName)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return fmt.Errorf("failed to parse config file: %s, err: %w", configFileName, err)
	}

	return nil
}

// Print given config to stdout.
func Print(config interface{}) {
	bytes, _ := yaml.Marshal(config)
	fmt.Println(string(bytes))
}
