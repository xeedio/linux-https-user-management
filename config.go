package humcommon

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	TLS struct {
		CA   string
		Key  string
		Cert string
	}
	Debug  bool
	Logger string
	URL    string
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
	configFile := "/etc/https-user-management.conf"

	var err error

	AppConfig, err = NewConfig(configFile)
	if err != nil {
		log.Fatalln("Unable to load config. Error: ", err)
	}

	err = initLogger()
	if err != nil {
		log.Fatalln("Unable to init logger: ", err)
	}

	err = initTLS()
	if err != nil {
		LogFatal("Unable to init tls", err)
	}
}
