package startup

import (
	"encoding/json"
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type ConfigFileType int

const (
	YAML ConfigFileType = iota
	JSON
)

// Generic function to load config into a struct that is passed in.
// Declare a struct for the config you want to receive with the appropriate
// markup tags matching the fileType you are using.
func LoadConfig[T interface{}](filePath string, fileType ConfigFileType, configType *T) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	switch fileType {
	case YAML:
		if parseErr := yaml.Unmarshal([]byte(file), configType); parseErr != nil {
			return fmt.Errorf("failed to parse yaml config from file: %s. %+v", filePath, parseErr)
		}
	case JSON:
		if parseErr := json.Unmarshal([]byte(file), configType); err != nil {
			return fmt.Errorf("failed to parse json config from file: %s. %+v", filePath, parseErr)
		}
	default:
		return fmt.Errorf("config files of the supplied type are not yet supported")
	}

	return nil
}
