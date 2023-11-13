package mongodb

import (
	"errors"
	"fmt"
)

type MongoConnectionDetails struct {
	Url      string `yaml:"url"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type MongoConnectionDetailsOmitAuth struct {
	Url      string `yaml:"url"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
}

func (mcd *MongoConnectionDetails) GetConnectionURI() (string, error) {
	// Return the registered URL or build one up from the component parts
	if mcd.Url != "" {
		return mcd.Url, nil
	} else {
		if mcd.Host == "" {
			return "", errors.New("missing mongo host")
		}
		if mcd.Port == 0 {
			return "", errors.New("missing mongo port")
		}
		if mcd.Database == "" {
			return "", errors.New("missing mongo database name")
		}
		return fmt.Sprintf("mongodb://%s:%d/%s", mcd.Host, mcd.Port, mcd.Database), nil
	}
}
