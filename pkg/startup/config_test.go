package startup

import (
	"testing"
)

func TestLoadJson(t *testing.T) {
	type TestJsonConfig struct {
		Name         string `json:"name"`
		OmittedValue int    `json:"value"`
		IsItATest    bool   `json:"test"`
	}

	tjc := TestJsonConfig{}
	err := LoadConfig[TestJsonConfig]("./test/test_config.json", JSON, &tjc)
	if err != nil {
		t.Error(err)
	}
	if tjc.Name != "test_config" {
		t.Errorf("Loaded config name should have been test_config, was %s", tjc.Name)
	}
	if tjc.OmittedValue != 0 {
		t.Errorf("Loaded config OmittedValue should have been left at its zero value, was %d", tjc.OmittedValue)
	}
	if tjc.IsItATest != true {
		t.Errorf("Loaded config IsItATest should have been true, was %t", tjc.IsItATest)
	}
}

func TestLoadYaml(t *testing.T) {
	type TestYamlConfig struct {
		Name         string `yaml:"name"`
		OmittedValue int    `yaml:"value"`
		IsItATest    bool   `yaml:"test"`
	}

	tyc := TestYamlConfig{}
	err := LoadConfig[TestYamlConfig]("./test/test_config.json", YAML, &tyc)
	if err != nil {
		t.Error(err)
	}
	if tyc.Name != "test_config" {
		t.Errorf("Loaded config name should have been test_config, was %s", tyc.Name)
	}
	if tyc.OmittedValue != 0 {
		t.Errorf("Loaded config OmittedValue should have been left at its zero value, was %d", tyc.OmittedValue)
	}
	if tyc.IsItATest != true {
		t.Errorf("Loaded config IsItATest should have been true, was %t", tyc.IsItATest)
	}
}

func TestInvalidFileType(t *testing.T) {
	type TestBsonConfig struct {
		Name string `bson:"name"`
	}

	tbc := TestBsonConfig{}
	err := LoadConfig[TestBsonConfig]("./test/test_config.json", 66, &tbc)
	if err == nil {
		t.Errorf("An error was expected, got none")
	}
	if err.Error() != "config files of the supplied type are not yet supported" {
		t.Errorf("Expected \"not supported\" error, got %s", err.Error())
	}
}

func TestInvalidFilePath(t *testing.T) {
	type EmptyConfig struct{}
	ec := EmptyConfig{}
	err := LoadConfig[EmptyConfig]("invalid_path.json", JSON, &ec)
	if err == nil {
		t.Errorf("No error was returned when passing an invalid file path")
	}
	expectedErrorMsg := "open invalid_path.json: no such file or directory"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Error message should have been \"%s\", was \"%s\"", expectedErrorMsg, err.Error())
	}
}
