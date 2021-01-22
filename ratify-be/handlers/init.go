package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/go-gomail/gomail"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/daystram/ratify/ratify-be/config"
	"github.com/daystram/ratify/ratify-be/datatransfers"
	"github.com/daystram/ratify/ratify-be/models"
)

// Handler is the global singleton to reference the handler.
var Handler handlerFunc

type handlerFunc interface {
	// auth
	AuthenticateUser(credentials datatransfers.UserLogin) (user models.User, err error)
	RegisterUser(credentials datatransfers.UserSignup) (userSubject string, err error)
	VerifyUser(token string) (err error)

	// user
	RetrieveUserBySubject(subject string) (user models.User, err error)
	RetrieveUserByUsername(username string) (user models.User, err error)
	RetrieveUserByEmail(email string) (user models.User, err error)
	UpdateUser(subject string, user datatransfers.UserUpdate) (err error)
	UpdateUserPassword(subject, oldPassword, newPassword string) (err error)

	// application
	RetrieveApplication(clientID string) (application models.Application, err error)
	RetrieveOwnedApplications(ownerSubject string) (applications []models.Application, err error)
	RetrieveAllApplications() (applications []models.Application, err error)
	RegisterApplication(application datatransfers.ApplicationInfo, ownerSubject string) (clientID, clientSecret string, err error)
	UpdateApplication(application datatransfers.ApplicationInfo) (err error)
	RenewApplicationClientSecret(clientID string) (clientSecret string, err error)
	DeleteApplication(clientID string) (err error)

	// oauth
	GenerateAuthorizationCode(authRequest datatransfers.AuthorizationRequest, subject, sessionID string) (authorizationCode string, err error)
	ValidateAuthorizationCode(application models.Application, authorizationCode string) (sessionID, subject, scope string, err error)
	GenerateAccessRefreshToken(tokenRequest datatransfers.TokenRequest, sessionID, subject string, withRefresh bool) (accessToken, refreshToken string, err error)
	GenerateIDToken(clientID, subject string, scope []string) (idToken string, err error)
	IntrospectAccessToken(accessToken string) (tokenInfo datatransfers.TokenIntrospection, err error)
	StoreCodeChallenge(authorizationCode string, pkce datatransfers.PKCEAuthFields) (err error)
	ValidateCodeVerifier(authorizationCode string, pkce datatransfers.PKCETokenFields) (err error)
	RevokeTokens(userSubject, clientID string, global bool) (err error)

	// session
	InitializeSession(subject string) (sessionID string, err error)
	CheckSession(sessionID string) (user models.User, newSessionID string, err error)
	ClearSession(sessionID string) (err error)

	// mailer
	SendVerificationEmail(user models.User) (err error)

	// mfa
	EnableTOTP(user models.User) (uri string, err error)
	ConfirmTOTP(otp string, user models.User) (err error)
	DisableTOTP(user models.User) (err error)
	CheckTOTP(otp string, user models.User) (valid bool)

	// log
	RetrieveActivityLogs(subject string) (logs []models.Log, err error)
	RetrieveAdminLogs() (logs []models.Log, err error)
	LogLogin(user models.User, application models.Application, success bool, detail datatransfers.LogDetail)
	LogUser(user models.User, success bool, detail datatransfers.LogDetail)
	LogApplication(user models.User, application models.Application, action bool, detail datatransfers.LogDetail)
}

type module struct {
	db     *dbEntity
	rd     *redis.Client
	mailer *gomail.Dialer
}

type dbEntity struct {
	conn             *gorm.DB
	applicationOrmer models.ApplicationOrmer
	userOrmer        models.UserOrmer
	logOrmer         models.LogOrmer
}

// InitializeHandler initializes application components and handler bundle.
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

	// Initialize Redis
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

	// Initialize Mailer
	var mailer *gomail.Dialer
	mailer = gomail.NewDialer(config.AppConfig.SMTPServer, config.AppConfig.SMTPPort, config.AppConfig.SMTPUsername, config.AppConfig.SMTPPassword)
	if _, err = mailer.Dial(); err != nil {
		log.Fatalf("[INIT] Failed authenticating to SMTP Server at %s:%d. %v\n",
			config.AppConfig.SMTPServer, config.AppConfig.SMTPPort, err)
	}
	log.Printf("[INIT] Successfully authenticathed to SMTP Server\n")

	// Compose handler modules
	Handler = &module{
		db: &dbEntity{
			conn:             db,
			applicationOrmer: models.NewApplicationOrmer(db),
			userOrmer:        models.NewUserOrmer(db),
			logOrmer:         models.NewLogOrmer(db),
		},
		rd:     rd,
		mailer: mailer,
	}

	log.Printf("[INIT] Initialization complete!\n")
}
