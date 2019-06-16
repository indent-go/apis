package pgproxy

// DefaultConfig contains sane values for running the proxy.
var DefaultConfig = Config{
	ServerHost:    "0.0.0.0:5432",
	ServerTLSKey:  "./certs/server.key",
	ServerTLSCert: "./certs/server.pem",

	ClaimSigningKey: "replace in staging and production",

	OriginHost:     "localhost:5433",
	OriginUsername: "demo",
	OriginDatabase: "demo",
}

// Config dictates how the proxy handles clients and connects to upstream Postgres servers.
type Config struct {
	// ServerHost is the address listened to by the proxy.
	ServerHost    string `json:"serverHost"`
	ServerTLSCert string `json:"serverTLSCert"`
	ServerTLSKey  string `json:"serverTLSKey"`

	// ClaimSigningKey is the key used for signing access request claims.
	ClaimSigningKey string `json:"claimSigningKey"`

	// OriginHost with port where the connection should be made, ie. 'localhost:5238'
	OriginHost     string `json:"originHost"`
	OriginUsername string `json:"originUsername"`
	OriginPassword string `json:"originPassword"`
	OriginDatabase string `json:"originDatabase"`
}
