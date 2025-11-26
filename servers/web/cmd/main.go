package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"zotes/servers/web/internal/app"
	"zotes/shared/create"
	"zotes/shared/structs"

	"github.com/joho/godotenv"
)

func main() {
	logger := create.CreateLogger()

	var err error
	if os.Getenv("ENVIRONMENT") != "prod" {
		err = godotenv.Load("../../envs/web.env")
		if err != nil {
			logger.Error("Error loading .env file", slog.Any("error", err))
		}
		err = godotenv.Load("../../envs/environment.env")
		if err != nil {
			logger.Error("Error loading .env file", slog.Any("error", err))
		}
		err = godotenv.Load("../../envs/s3.env")
		if err != nil {
			logger.Error("Error loading .env file", slog.Any("error", err))
		}
		err = godotenv.Load("../../envs/pdfworker.env")
		if err != nil {
			logger.Error("Error loading .env file", slog.Any("error", err))
		}
		// err = godotenv.Load("../../envs/database.env")
		// if err != nil {
		// 	logger.Error("Error loading .env file", slog.Any("error", err))
		// }

	}

	if os.Getenv("ENVIRONMENT") == "dev" {
		logger.Debug("Reminder to start docker services with 'make on' at root, if database hasn't been setup use 'make up' to migrate.")
	}

	addr := flag.String(("addr"), ":4444", "HTTP network address")

	// dbauth := &env.DBENV{}
	// s3auth := &env.S3ENV{}
	// _ = &structs.S3Auth{
	// 	Bucket:    os.Getenv("AWS_BUCKET"),
	// 	Endpoint:  os.Getenv("AWS_ENDPOINT"),
	// 	Region:    os.Getenv("AWS_REGION"),
	// 	AccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
	// 	SecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	// }

	dbauth := &structs.DBAuth{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		Schema:   os.Getenv("DB_SCHEMA"),
	}

	pool := create.PGPool(dbauth)
	app := &app.App{
		DB:                pool,
		Logger:            logger,
		Session:           create.PGSessionManager(pool, "zotes-client"),
		OidcAuthenticator: create.OidcAuthenticator("CLIENT_UI"),
		// JwtIssuer:         create.JwtIssuer(""),
		// JwtValidator:      create.JwtValidator(""),
		Rbac: create.Casbin(dbauth),
	}

	defer pool.Close()
	routes := GetRoutes(app)
	// InsertRoutes(routes, deps)
	err = http.ListenAndServe(*addr, routes)

	app.Logger.Error("Server stopped unexpectedly" + err.Error())

	// REMEMBER TO USE JIRA
	app.Logger.Error("Server stopped unexpectedly" + err.Error())

	// s3, err := create.GetS3Manager(s3auth)
	// if err != nil {
	// 	logger.Error("Error getting s3 manager", slog.Any("error", err))
	// 	return
	// }
	//
	// pool, err := create.GetPGPool(dbauth)
	// if err != nil {
	// 	logger.Error("Error getting pool", slog.Any("error", err))
	// 	return
	// }
	// enforcer, err := create.GetCasbin(dbauth)
	// if err != nil {
	// 	logger.Error("Error getting casbin", slog.Any("error", err))
	// 	return
	// }
	// /
	// oidc, err := create.GetOidcAuthenticator("CLIENT_UI")
	// if err != nil {
	// 	logger.Error("Error getting authenticator", slog.Any("error", err))
	// 	return
	// }
	//
	// jwtissuer, err := create.GetJwtIssuer("")
	// // if err != nil {
	// // 	logger.Error("Error getting jwt issuer", slog.Any("error", err))
	// // }
	//
	// // this is expected to be nil
	// jwtvalidator, _ := create.GetJwtValidator("")
	// if err != nil {
	// 	logger.Error("Error getting jwt validator", slog.Any("error", err))
	// a

	// deps := routes.NewDependencies(
	// 	pool,
	// 	func(tx pgx.Tx) *db.Queries {
	// 		return db.New(tx)
	// 	},
	// 	logger,
	// 	create.GetPGSessionManager(pool, "zote"),
	// 	create.GetSanitizer(),
	// 	oidc,
	// 	jwtissuer,
	// 	jwtvalidator,
	// 	enforcer,
	// 	s3,
	// 	// create.GetCoreApiTarget(),
	// )
	// defer deps.Close()
	//
	// deps.Logger().Info("Starting server", slog.String("addr", *addr))
	//
	// routes := GetRoutes(deps)
	// InsertRoutes(routes, deps)
	// err = http.ListenAndServe(*addr, routes)
	//
	// deps.Logger().Error("Server stopped unexpectedly" + err.Error())
	//
	// // REMEMBER TO USE JIRA
	// logger.Error("Server stopped unexpectedly" + err.Error())
}
