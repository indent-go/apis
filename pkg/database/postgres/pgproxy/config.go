package pgproxy

// DefaultConfig contains sane values for running PGA.
var DefaultConfig = Config{
	ServerHost:    "0.0.0.0:5432",
	ServerTLSKey:  "./certs/server.key",
	ServerTLSCert: "./certs/server.pem",

	OriginHost:     "localhost:5433",
	OriginUsername: "demo",
	OriginDatabase: "demo",
}

// Config dictates how PGA handles clients and connects to upstream Postgres servers.
type Config struct {
	// ServerHost is the address listened to by PGA.
	ServerHost    string `json:"serverHost"`
	ServerTLSCert string `json:"serverTLSCert"`
	ServerTLSKey  string `json:"serverTLSKey"`

	// OriginHost with port where the connection should be made, ie. 'localhost:5238'
	OriginHost     string `json:"originHost"`
	OriginUsername string `json:"originUsername"`
	OriginPassword string `json:"originPassword"`
	OriginDatabase string `json:"originDatabase"`
}
