package config

// Credentials contains user and password
type Credentials struct {
	User     string
	Password string
}

// TLSCertificate contains CA certificate, public keu (certificate) and private key
type TLSCertificate struct {
	CA          []byte
	Certificate []byte
	Key         []byte
}

// Secret is interface to get secret value
type Secret interface {
	Config

	// GetCredentials returns configuration value as Credentials
	// Calling with parameter key "cred" will get configuration values from "cred.user" and "cred.password"
	GetCredentials(string) Credentials

	// GetTLSCertificate returns configuration value as TLSCertificate
	// Calling with parameter key "tls" will get configuration values from "tls.ca", "tls.cert" and "tls.key"
	GetTLSCertificate(string) TLSCertificate
}
