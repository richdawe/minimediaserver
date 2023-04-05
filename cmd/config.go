package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/richdawe/minimediaserver/services/catalog"
	"github.com/richdawe/minimediaserver/services/storage"
	"github.com/spf13/viper"
)

type StorageServiceConfig struct {
	Type string `mapstructure:"type"`
	Path string `mapstructure:"path"`
}

type Config struct {
	Addr            string // Server IP + port
	StorageServices []StorageServiceConfig
}

func setLoadConfigOptions() {
	viper.SetConfigName(".minimediaserver")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
}

func loadConfig() (Config, error) {
	// Set some defaults
	viper.SetDefault("host", "127.0.0.1")
	viper.SetDefault("port", "1323")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// TODO: info log
			fmt.Println("no configuration file found - using defaults")
		} else {
			return Config{}, err
		}
	}

	var config Config

	host := viper.GetString("host")
	if host == "*" {
		host = ""
	}
	port := viper.GetString("port")
	config.Addr = host + ":" + port

	err := viper.UnmarshalKey("storageServices", &config.StorageServices)
	if err != nil {
		return Config{}, err
	}

	for _, css := range config.StorageServices {
		fmt.Printf("%+v\n", css)
	}

	return config, nil
}

func buildCatalog(config Config) (*catalog.CatalogService, error) {
	catalogService, err := catalog.New()
	if err != nil {
		return nil, err
	}

	// Always provide at least one storage service.
	if len(config.StorageServices) == 0 {
		config.StorageServices = []StorageServiceConfig{
			{
				Type: "nullStorage",
			},
		}
	}

	for _, css := range config.StorageServices {
		var ss storage.StorageService
		var err error

		switch css.Type {
		case "nullStorage":
			ss, err = storage.NewNullStorage()
		case "diskStorage":
			path := css.Path
			if path == "" {
				path = "."
			}
			path = strings.Replace(path, "$HOME", os.Getenv("HOME"), -1)
			ss, err = storage.NewDiskStorage(path)
		default:
			err = fmt.Errorf("unknown storage service %s", css.Type)
		}
		if err != nil {
			return nil, err
		}

		err = catalogService.AddStorage(ss)
		if err != nil {
			return nil, err
		}
	}

	return catalogService, nil
}
