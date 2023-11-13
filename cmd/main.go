package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"mk.io/controller/mopi/pkg/mongodb"
	"mk.io/controller/mopi/pkg/server"
	"mk.io/controller/mopi/pkg/startup"
)

type MopiConfig struct {
	Storage mongodb.MongoConnectionDetails
	App     struct {
		Port int
	}
}

type CLIArgs struct {
	Config string
}

func main() {
	args := GetCLIArgs()
	config := MopiConfig{}
	err := startup.LoadConfig[MopiConfig](args.Config, startup.YAML, &config)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to load config %s", err.Error()))
		os.Exit(1)
	}
	slog.Info(fmt.Sprintf("Running on port: %d", config.App.Port))

	uri, _ := config.Storage.GetConnectionURI()
	store, err := server.NewMongoStorer(uri)
	if err != nil {
		slog.Error(fmt.Sprintf("unable to connect to mongo: %s", err.Error()))
		os.Exit(1)
	}

	server := server.NewServer(config.App.Port, store)
	server.LoadResponsesFromStore()
	server.Run()

}

func GetCLIArgs() *CLIArgs {
	config := flag.String("config", "", "The config path to try and load from")
	flag.Parse()
	args := CLIArgs{}

	if config != nil {
		args.Config = *config
	}

	return &args

}
