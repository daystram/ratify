package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/daystram/ratify/ratify-be/config"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/models"
)

var Handler HandlerFunc

type HandlerFunc interface {
	AuthenticateUser(credentials datatransfers.UserLogin) (token string, err error)
	RegisterUser(credentials datatransfers.UserSignup) (err error)
	RetrieveUserBySubject(subject string) (user models.User, err error)
	RetrieveUserByUsername(username string) (user models.User, err error)
	RetrieveUserByEmail(email string) (user models.User, err error)
	UpdateUser(id string, user datatransfers.UserUpdate) (err error)

	RetrieveApplication(clientID string) (application models.Application, err error)
	RetrieveOwnedApplications(ownerSubject string) (applications []models.Application, err error)
	RegisterApplication(application datatransfers.ApplicationInfo, ownerSubject string) (clientID string, err error)
	UpdateApplication(application datatransfers.ApplicationInfo) (err error)

	GenerateAuthorizationCode(application models.Application) (authorizationCode string, err error)
	ValidateAuthorizationCode(application models.Application, authorizationCode string) (err error)
	GenerateAccessRefreshToken(application models.Application) (accessToken, refreshToken string, err error)
	StoreCodeChallenge(authorizationCode string, pkce datatransfers.PKCEAuthFields) (err error)
	ValidateCodeVerifier(authorizationCode string, pkce datatransfers.PKCETokenFields) (err error)
}

type module struct {
	db *dbEntity
	rd *redis.Client
}

type dbEntity struct {
	conn             *gorm.DB
	applicationOrmer models.ApplicationOrmer
	userOrmer        models.UserOrmer
}

func InitializeHandler() {
	var err error

	// Initialize DB
	var db *gorm.DB
	db, err = gorm.Open(postgres.Open(
		fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
			config.AppConfig.DBHostname, config.AppConfig.DBPort, config.AppConfig.DBDatabase,
			config.AppConfig.DBUsername, config.AppConfig.DBPassword),
	), &gorm.Config{})
	if err != nil {
		log.Fatalf("[INIT] Failed connecting to PostgreSQL Database at %s:%d. %v\n",
			config.AppConfig.DBHostname, config.AppConfig.DBPort, err)
	}
	log.Printf("[INIT] Successfully connected to PostgreSQL Database\n")

	//Initialize Redis
	var rd *redis.Client
	rd = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.AppConfig.RedisHostname, config.AppConfig.RedisPort),
		Password: config.AppConfig.RedisPassword,
		DB:       config.AppConfig.RedisDatabase,
	})
	if err = rd.Info(context.Background()).Err(); err != nil {
		log.Fatalf("[INIT] Failed connecting to Redis at %s:%d. %v\n",
			config.AppConfig.RedisHostname, config.AppConfig.RedisPort, err)
	}
	log.Printf("[INIT] Successfully connected to Redis\n")

	// Compose handler modules
	Handler = &module{
		db: &dbEntity{
			conn:             db,
			applicationOrmer: models.NewApplicationOrmer(db),
			userOrmer:        models.NewUserOrmer(db),
		},
		rd: rd,
	}
}
