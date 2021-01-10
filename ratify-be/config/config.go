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

	SMTPServer   string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string

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

	if AppConfig.Hostname = viper.GetString("hostname"); AppConfig.Hostname == "" {
		log.Fatalln("[INIT] hostname is missing in config.yaml")
	}
	AppConfig.Port = viper.GetInt("port")
	if AppConfig.Environment = viper.GetString("environment"); AppConfig.Environment == "" {
		log.Fatalln("[INIT] environment is missing in config.yaml")
	}
	AppConfig.Debug = viper.GetBool("debug")

	AppConfig.DBHostname = viper.GetString("db_hostname")
	AppConfig.DBPort = viper.GetInt("db_port")
	AppConfig.DBDatabase = viper.GetString("db_database")
	AppConfig.DBUsername = viper.GetString("db_username")
	AppConfig.DBPassword = viper.GetString("db_password")

	AppConfig.RedisHostname = viper.GetString("redis_hostname")
	AppConfig.RedisPort = viper.GetInt("redis_port")
	AppConfig.RedisPassword = viper.GetString("redis_password")
	AppConfig.RedisDatabase = viper.GetInt("redis_database")

	AppConfig.SMTPServer = viper.GetString("smtp_server")
	AppConfig.SMTPPort = viper.GetInt("smtp_port")
	AppConfig.SMTPUsername = viper.GetString("smtp_username")
	AppConfig.SMTPPassword = viper.GetString("smtp_password")

	if AppConfig.Domain = viper.GetString("domain"); AppConfig.Domain == "" {
		log.Fatalln("[INIT] domain is missing in config.yaml")
	}
	if AppConfig.JWTSecret = viper.GetString("secret"); AppConfig.JWTSecret == "" {
		log.Fatalln("[INIT] secret is missing in config.yaml")
	}

	log.Printf("[INIT] Configuration loaded from %s\n", viper.ConfigFileUsed())
}
