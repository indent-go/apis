package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/crunchydata/crunchy-proxy/config"
	"go.indent.com/apis/pkg/providers/postgres/pgproxy"
	"sigs.k8s.io/yaml"
)

func init() {
	flag.Parse()
}

var ConfigFilePath = "./config.yaml"

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	cfg, err := readConfig(ConfigFilePath)
	if err != nil {
		logger.Printf("Unable to read config from '%s': %v", ConfigFilePath, err)
		cfg = pgproxy.DefaultConfig
	}

	// add origin servers provided in environment
	if pwd, ok := os.LookupEnv("IPA_ORIGIN_PASSWORD"); ok {
		cfg.OriginPassword = pwd
	}

	if host, ok := os.LookupEnv("IPA_ORIGIN_HOST"); ok {
		cfg.OriginHost = host
	}

	// TODO: remove
	config.SetConfigPath("./origin-config.yaml")
	config.ReadConfig()

	server := &pgproxy.Server{
		Config: cfg,
		Logger: logger,
	}

	if err = server.Setup(); err != nil {
		logger.Fatalf("Could not setup server: %v", err)
	}

	stopCh := make(chan struct{})
	if err = server.Start(stopCh); err != nil {
		logger.Fatalf("Error running server: %v", err)
	}
}

func readConfig(filename string) (cfg pgproxy.Config, err error) {
	if len(flag.Arg(0)) != 0 {
		filename = flag.Arg(0)
	}

	cfgData, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(cfgData, &cfg)
	return
}
