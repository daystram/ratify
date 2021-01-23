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
	viper.SetConfigName(".env") // for development
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/")
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	if AppConfig.Hostname = viper.GetString("HOSTNAME"); AppConfig.Hostname == "" {
		log.Fatalln("[INIT] HOSTNAME is not set")
	}
	AppConfig.Port = viper.GetInt("PORT")
	if AppConfig.Environment = viper.GetString("ENVIRONMENT"); AppConfig.Environment == "" {
		log.Fatalln("[INIT] ENVIRONMENT is not set")
	}
	AppConfig.Debug = viper.GetBool("DEBUG")

	AppConfig.DBHostname = viper.GetString("DB_HOSTNAME")
	AppConfig.DBPort = viper.GetInt("DB_PORT")
	AppConfig.DBDatabase = viper.GetString("DB_DATABASE")
	AppConfig.DBUsername = viper.GetString("DB_USERNAME")
	AppConfig.DBPassword = viper.GetString("DB_PASSWORD")

	AppConfig.RedisHostname = viper.GetString("REDIS_HOSTNAME")
	AppConfig.RedisPort = viper.GetInt("REDIS_PORT")
	AppConfig.RedisPassword = viper.GetString("REDIS_PASSWORD")
	AppConfig.RedisDatabase = viper.GetInt("REDIS_DATABASE")

	AppConfig.SMTPServer = viper.GetString("SMTP_SERVER")
	AppConfig.SMTPPort = viper.GetInt("SMTP_PORT")
	AppConfig.SMTPUsername = viper.GetString("SMTP_USERNAME")
	AppConfig.SMTPPassword = viper.GetString("SMTP_PASSWORD")

	if AppConfig.Domain = viper.GetString("DOMAIN"); AppConfig.Domain == "" {
		log.Fatalln("[INIT] DOMAIN is not set")
	}
	if AppConfig.JWTSecret = viper.GetString("SECRET"); AppConfig.JWTSecret == "" {
		log.Fatalln("[INIT] SECRET is not set")
	}

	log.Printf("[INIT] Configuration loaded!")
}
