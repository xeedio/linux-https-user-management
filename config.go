package humcommon

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"log/syslog"
	"net/http"
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

var logPrefix = "LHUM "
var logger *log.Logger
var transport *http.Transport

func GetHTTPClient() *http.Client {
	return &http.Client{Transport: transport}
}

func SetLogPrefix(p string) {
	logPrefix = p
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

	fmt.Printf("AppConfig: %+v\n", AppConfig)

	err = initLogger()
	if err != nil {
		log.Fatalln("Unable to init logger: ", err)
	}

	err = initTLS()
	if err != nil {
		LogFatal("Unable to init tls: ", err)
	}
}

func initLogger() error {
	var err error
	if AppConfig.Logger == "syslog" {
		logger, err = syslog.NewLogger(syslog.LOG_INFO, log.Ltime)
		if err != nil {
			logger = log.New(os.Stderr, logPrefix, log.Ltime)
			logger.Printf("syslog.Open() err: %v", err)
		}
	} else {
		logger = log.New(os.Stdout, logPrefix, log.Ltime)
	}

	return nil
}

func initTLS() error {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	certPath := AppConfig.TLS.Cert
	keyPath := AppConfig.TLS.Key
	if certPath != "" && keyPath != "" {
		mainCert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			return err
		}
		tlsConfig.Certificates = []tls.Certificate{mainCert}
		tlsConfig.BuildNameToCertificate()
	} else {
		LogInfo("TLS", "Unable to init client certificates: either cert or key missing")
	}

	if AppConfig.TLS.CA != "" {
		caCert, err := ioutil.ReadFile(AppConfig.TLS.CA)
		if err != nil {
			return err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	transport = &http.Transport{TLSClientConfig: tlsConfig}

	LogDebug("TLS", fmt.Sprintf("client_cert=%v,insecureSkipVerify=%v", len(tlsConfig.Certificates) > 0, tlsConfig.InsecureSkipVerify))

	return nil
}

func LogDebug(module string, data ...interface{}) {
	if !AppConfig.Debug {
		return
	}
	module = "DEBUG-" + module
	LogInfo(module, data)
}

func LogInfo(module string, data ...interface{}) {
	s := fmt.Sprintf("[%s] %s", module, fmt.Sprint(data...))
	logger.Println(s)
}

func LogFatal(module string, err error) {
	s := fmt.Sprintf("[%s] Fatal: %v", module, err)
	logger.Fatal(s)
}

type User struct {
	Admin    bool   `json:"is_staff"`
	Email    string `json:"email"`
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type TokenUser struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func Authenticate(user, password string) (*TokenUser, error) {
	LogDebug("API-AUTHENTICATE", fmt.Sprintf("Making request to %s", AppConfig.URL))

	authStruct := struct {
		User     string `json:"username"`
		Password string `json:"password"`
	}{
		user,
		password,
	}

	b, err := json.Marshal(authStruct)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(b)

	client := GetHTTPClient()
	resp, err := client.Post(AppConfig.URL, "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	LogDebug("API-AUTHENTICATE", fmt.Sprintf("Response: code=%d(%s),length=%d,content-type=%s", resp.StatusCode, resp.Status, resp.ContentLength, resp.Header.Get("Content-Type")))

	if resp.StatusCode != 200 {
		b, _ = ioutil.ReadAll(resp.Body)
		LogInfo("API-AUTHENTICATE", fmt.Sprintf("Invalid response body: %s", b))
		return nil, errors.New(resp.Status)
	}

	tokenUser := TokenUser{}
	err = json.NewDecoder(resp.Body).Decode(&tokenUser)
	if err != nil {
		return nil, err
	}

	return &tokenUser, nil
}
