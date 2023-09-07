package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

type Config struct {
	Server Server
	DB     DB
	//Adapters Adapters
}

type Server struct {
	Name                string
	AppVersion          string
	Port                string
	BaseURI             string
	Mode                string
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	SSL                 bool
	CtxDefaultTimeout   time.Duration
	CSRF                bool
	Debug               bool
	MaxCountRequest     int           // max count of connections
	ExpirationLimitTime time.Duration //  expiration time of the limit
}

type DB struct {
	Mysql Mysql
}

type Mysql struct {
	Host            string
	Port            int
	UserName        string
	Password        string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	Options         string
}

// Get config path for local or docker
func getDefaultConfig() string {
	return "./config/config"
}

// Load config file from given path
func NewConfig() (*Config, error) {
	config := Config{}
	path := os.Getenv("cfgPath")
	if path == "" {
		path = getDefaultConfig()
	}

	v := viper.New()

	v.SetConfigName(path)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	err := v.Unmarshal(&config)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &config, nil
}
