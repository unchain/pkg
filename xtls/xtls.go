package xtls

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"

	"github.com/pkg/errors"
)

// KeyPairConfig contains the private key and certificate for TLS encryption
type KeyPairConfig struct {
	Key  *KeyConfig
	Cert *KeyConfig
}

func (cfg *KeyPairConfig) X509KeyPair() (tls.Certificate, error) {
	certBytes, err := cfg.Cert.Bytes()

	if err != nil {
		return tls.Certificate{}, err
	}

	keyBytes, err := cfg.Key.Bytes()

	if err != nil {
		return tls.Certificate{}, err
	}

	clientCert, err := tls.X509KeyPair(certBytes, keyBytes)

	if err != nil {
		return tls.Certificate{}, errors.Wrap(err, "failed to parse TLS keys")
	}

	return clientCert, nil
}

// KeyConfig manages the source of a TLS key, coming either from a file path or a specified PEM formatted string
type KeyConfig struct {
	Path string
	Pem  string
}

// String returns the tls certificate as a string by loading it either from the embedded Pem or Path
func (cfg *KeyConfig) String() (string, error) {
	bytes, err := cfg.Bytes()

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// Bytes returns the tls certificate as a byte array by loading it either from the embedded Pem or Path
func (cfg *KeyConfig) Bytes() ([]byte, error) {
	var bytes []byte
	var err error

	if cfg.Pem != "" {
		bytes = []byte(cfg.Pem)
	} else if cfg.Path != "" {
		bytes, err = ioutil.ReadFile(cfg.Path)

		if err != nil {
			return nil, errors.Wrapf(err, "failed to load pem bytes from path %s", cfg.Path)
		}
	}

	return bytes, nil
}

// Bytes returns the tls certificate as a byte array by loading it either from the embedded Pem or Path
func (cfg *KeyConfig) DecodedBytes() ([]byte, error) {
	bytes, err := cfg.Bytes()

	if err != nil {
		return nil, err
	}

	data, _ := pem.Decode(bytes)

	return data.Bytes, nil
}

func (cfg *KeyConfig) ParsePKCS1PrivateKey() (*rsa.PrivateKey, error) {
	bytes, err := cfg.DecodedBytes()

	if err != nil {
		return nil, err
	}

	rsaPriv, err := x509.ParsePKCS1PrivateKey(bytes)

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse private key")
	}

	return rsaPriv, nil
}

func (cfg *KeyConfig) ParsePKCS1PublicKey2() (*rsa.PublicKey, error) {
	bytes, err := cfg.DecodedBytes()

	if err != nil {
		return nil, err
	}

	publicKeyImported, err := x509.ParsePKIXPublicKey(bytes)

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse public key")
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

	if !ok {
		return nil, errors.New("incorrect public key type")
	}

	return rsaPub, nil
}

func (cfg *KeyConfig) ParsePKCS1PublicKey() (*rsa.PublicKey, error) {
	bytes, err := cfg.DecodedBytes()

	if err != nil {
		return nil, err
	}

	rsaPub, err := x509.ParsePKCS1PublicKey(bytes)

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse public key")
	}

	return rsaPub, nil
}

// TLSCert returns the tls certificate as a *x509.Certificate by loading it either from the embedded Pem or Path
func (cfg *KeyConfig) ParseCertificate() (*x509.Certificate, error) {
	bytes, err := cfg.DecodedBytes()

	if err != nil {
		return nil, err
	}

	if len(bytes) > 0 {
		pub, err := x509.ParseCertificate(bytes)

		if err != nil {
			return nil, errors.Wrap(err, "certificate parsing failed")
		}

		return pub, nil
	}

	// return an error with an error code for clients to test against status.EmptyCert code
	return nil, errors.New("pem data missing")
}
