package config

import (
	"log"

	"github.com/spf13/viper"
)

var AppConfig Config

type Config struct {
	Hostname    string
	Port        int
	Environment string
	Debug       bool

	DBHostname string
	DBPort     int
	DBDatabase string
	DBUsername string
	DBPassword string

	RedisHostname string
	RedisPort     int
	RedisPassword string
	RedisDatabase int

	Domain    string
	JWTSecret string
}

func InitializeAppConfig() {
	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("[INIT] Unable to load configuration. %+v\n", err)
	}

	// Hostname
	if AppConfig.Hostname = viper.GetString("hostname"); AppConfig.Hostname == "" {
		log.Fatalln("[INIT] hostname is missing in config.yaml")
	}

	// Port
	AppConfig.Port = viper.GetInt("port")

	// Environment
	if AppConfig.Environment = viper.GetString("environment"); AppConfig.Environment == "" {
		log.Fatalln("[INIT] environment is missing in config.yaml")
	}

	// Debug
	AppConfig.Debug = viper.GetBool("debug")

	// DBHostname
	AppConfig.DBHostname = viper.GetString("db_hostname")

	// DBPort
	AppConfig.DBPort = viper.GetInt("db_port")

	// DBDatabase
	AppConfig.DBDatabase = viper.GetString("db_database")

	// DBUsername
	AppConfig.DBUsername = viper.GetString("db_username")

	// DBPassword
	AppConfig.DBPassword = viper.GetString("db_password")

	// RedisHostname
	AppConfig.RedisHostname = viper.GetString("redis_hostname")

	// RedisPort
	AppConfig.RedisPort = viper.GetInt("redis_port")

	// RedisPassword
	AppConfig.RedisPassword = viper.GetString("redis_password")

	// RedisDatabase
	AppConfig.RedisDatabase = viper.GetInt("redis_database")

	// Domain
	if AppConfig.Domain = viper.GetString("domain"); AppConfig.Domain == "" {
		log.Fatalln("[INIT] domain is missing in config.yaml")
	}

	// JWTSecret
	if AppConfig.JWTSecret = viper.GetString("secret"); AppConfig.JWTSecret == "" {
		log.Fatalln("[INIT] secret is missing in config.yaml")
	}

	log.Printf("[INIT] Configuration loaded from %s\n", viper.ConfigFileUsed())
}
