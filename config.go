package humcommon

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const defaultConfigFilePath = "/etc/https-user-management/config.conf"
const defaultTokenFilePath = "/etc/https-user-management/user.token"

type Config struct {
	TLS struct {
		CA   string
		Key  string
		Cert string
	}
	Debug     bool
	Logger    string
	URL       string
	TokenFile string
}

func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

var AppConfig *Config

func init() {
	var err error

	configFilePath := os.Getenv("HTTPS_USER_MANAGEMENT_CONFIG")
	if configFilePath == "" {
		configFilePath = defaultConfigFilePath
	}

	AppConfig, err = NewConfig(configFilePath)
	if err != nil {
		log.Fatalln("Unable to load config. Error: ", err)
	}
	if AppConfig.TokenFile == "" {
		AppConfig.TokenFile = defaultTokenFilePath
	}

	err = initLogger()
	if err != nil {
		log.Fatalln("Unable to init logger: ", err)
	}

	LogInfo("APP-CONFIG", fmt.Sprintf("%+v", AppConfig))

	err = initTLS()
	if err != nil {
		LogFatal("Unable to init tls", err)
	}
}
