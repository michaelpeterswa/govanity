package settings

import (
	"errors"
	"os"

	"github.com/michaelpeterswa/govanity/cmd/internal/structs"
	"gopkg.in/yaml.v2"
)

func InitSettings() (*structs.Settings, error) {
	fileSettings, err := os.ReadFile("config/settings.yaml")
	if err != nil {
		return nil, errors.New(err.Error())
	}

	var set structs.Settings

	err = yaml.Unmarshal(fileSettings, &set)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return &set, nil
}
