package humcommon

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"time"
)

var transport *http.Transport

func GetHTTPClient() *http.Client {
	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}
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
		logger.Info("Unable to init client certificates: either cert or key missing")
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

	return nil
}
