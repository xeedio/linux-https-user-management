package humcommon

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const defaultConfigFilePath = "/etc/https-user-management/config.yaml"
const defaultTokenFilePath = "/etc/https-user-management/user.token"

var AppConfig *Config
var ConfigError bool

type Config struct {
	TLS struct {
		CA   string
		Key  string
		Cert string
	}
	Debug     bool
	URL       string
	TokenFile string
}

func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}

func init() {
	var err error

	initLogger()

	configFilePath := os.Getenv("HTTPS_USER_MANAGEMENT_CONFIG")
	if configFilePath == "" {
		configFilePath = defaultConfigFilePath
	}

	AppConfig, err = NewConfig(configFilePath)
	if err != nil {
		logger.Warnf("Unable to load config. Error: %v", err)
		ConfigError = true
		return
	}

	if AppConfig.TokenFile == "" {
		AppConfig.TokenFile = defaultTokenFilePath
	}

	logger.Infof("AppConfig: %+v", *AppConfig)

	if AppConfig.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	err = initTLS()
	if err != nil {
		ConfigError = true
		logger.Warnf("Unable to init tls: %v", err)
		return
	}
}
