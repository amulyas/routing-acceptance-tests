package helpers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cloudfoundry-incubator/cf-test-helpers/helpers"
	"github.com/nu7hatch/gouuid"
)

type RoutingConfig struct {
	helpers.Config
	RoutingApiUrl string       `json:"-"` //"-" is used for ignoring field
	Addresses     []string     `json:"addresses"`
	OAuth         *OAuthConfig `json:"oauth"`
}

type OAuthConfig struct {
	TokenEndpoint     string `json:"token_endpoint"`
	ClientName        string `json:"client_name"`
	ClientSecret      string `json:"client_secret"`
	Port              int    `json:"port"`
	SkipSSLValidation bool   `json:"skip_ssl_validation"`
}

func LoadConfig() RoutingConfig {
	loadedConfig := loadConfigJsonFromPath()
	loadedConfig.Config = helpers.LoadConfig()

	if loadedConfig.OAuth == nil {
		panic("missing configuration oauth")
	}

	if len(loadedConfig.Addresses) == 0 {
		panic("missing configuration 'addresses'")
	}

	if loadedConfig.AppsDomain == "" {
		panic("missing configuration apps_domain")
	}

	if loadedConfig.ApiEndpoint == "" {
		panic("missing configuration api")
	}

	loadedConfig.RoutingApiUrl = fmt.Sprintf("%s%s", loadedConfig.Protocol(), loadedConfig.ApiEndpoint)

	return loadedConfig
}

func loadConfigJsonFromPath() RoutingConfig {
	var config RoutingConfig

	path := configPath()

	configFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	return config
}

func configPath() string {
	path := os.Getenv("CONFIG")
	if path == "" {
		panic("Must set $CONFIG to point to an integration config .json file.")
	}

	return path
}

func (c RoutingConfig) Protocol() string {
	if c.UseHttp {
		return "http://"
	} else {
		return "https://"
	}
}

func RandomName() string {
	guid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return guid.String()
}
