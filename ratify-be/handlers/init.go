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
	AuthAuthenticate(credentials datatransfers.UserLogin) (user models.User, err error)
	AuthRegister(credentials datatransfers.UserSignup) (userSubject string, err error)
	AuthVerify(token string) (err error)

	// user
	UserGetOneBySubject(subject string) (user models.User, err error)
	UserGetOneByUsername(username string) (user models.User, err error)
	UserGetOneByEmail(email string) (user models.User, err error)
	UserUpdate(subject string, user datatransfers.UserUpdate) (err error)
	UserUpdatePassword(subject, oldPassword, newPassword string) (err error)

	// application
	ApplicationGetOneByClientID(clientID string) (application models.Application, err error)
	ApplicationGetOneByOwnerSubject(ownerSubject string) (applications []models.Application, err error)
	ApplicationGetAll() (applications []models.Application, err error)
	ApplicationRegister(application datatransfers.ApplicationInfo, ownerSubject string) (clientID, clientSecret string, err error)
	ApplicationUpdate(application datatransfers.ApplicationInfo) (err error)
	ApplicationRenewClientSecret(clientID string) (clientSecret string, err error)
	ApplicationDelete(clientID string) (err error)

	// oauth
	OAuthGenerateAuthorizationCode(authRequest datatransfers.AuthorizationRequest, subject, sessionID string) (authorizationCode string, err error)
	OAuthValidateAuthorizationCode(application models.Application, authorizationCode string) (sessionID, subject, scope string, err error)
	OAuthGenerateAccessToken(tokenRequest datatransfers.TokenRequest, sessionID, subject string, withRefresh bool) (accessToken, refreshToken string, err error)
	OAuthGenerateIDToken(clientID, subject string, scope []string) (idToken string, err error)
	OAuthIntrospectAccessToken(accessToken string) (tokenInfo datatransfers.TokenIntrospection, err error)
	OAuthStoreCodeChallenge(authorizationCode string, pkce datatransfers.PKCEAuthFields) (err error)
	OAuthValidateCodeVerifier(authorizationCode string, pkce datatransfers.PKCETokenFields) (err error)
	OAuthRevokeAccessToken(accessToken string) (err error)
	OAuthRevokeAllTokens(userSubject, clientID string, global bool) (err error)

	// session
	SessionInitialize(subject string, userAgent datatransfers.UserAgent) (sessionID string, err error)
	SessionInfo(sessionID string) (session datatransfers.SessionInfo, err error)
	SessionRevoke(sessionID string) (err error)
	SessionAddChild(sessionID, accessToken string) (err error)

	// mailer
	MailerSendEmailVerification(user models.User) (err error)

	// mfa
	MFAEnableTOTP(user models.User) (uri string, err error)
	MFAConfirmTOTP(otp string, user models.User) (err error)
	MFADisableTOTP(user models.User) (err error)
	MFACheckTOTP(otp string, user models.User) (valid bool)

	// log
	LogGetAllActivity(subject string) (logs []models.Log, err error)
	LogGetAllAdmin() (logs []models.Log, err error)
	LogInsertLogin(user models.User, application models.Application, success bool, detail datatransfers.LogDetail)
	LogInsertUser(user models.User, success bool, detail datatransfers.LogDetail)
	LogInsertApplication(user models.User, application models.Application, action bool, detail datatransfers.LogDetail)
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
