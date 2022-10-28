package app

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"os"
	"os/signal"
	"redditclone/internal/handler"
	"redditclone/internal/repository/mongorepo"
	"redditclone/internal/repository/mysqlrepo"
	"redditclone/internal/service"
	"redditclone/pkg/cookie"
	"redditclone/pkg/token"
	"syscall"
)

func initMySQL(cfg MySQLConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		if err := db.Close(); err != nil {
			return nil, err
		}
		return nil, err
	}
	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	return db, nil
}

func initMongoDB(cfg MongoConfig) (*mongo.Client, error) {
	mongoURL := fmt.Sprintf("mongodb://%s:%s",
		cfg.Host,
		cfg.Port,
	)
	credential := options.Credential{
		Username: cfg.Username,
		Password: cfg.Password,
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL).SetAuth(credential))
	if err != nil {
		return nil, err
	}
	if err = client.Connect(context.TODO()); err != nil {
		return nil, err
	}
	if err = client.Ping(context.TODO(), nil); err != nil {
		if err := client.Disconnect(context.TODO()); err != nil {
			logrus.Errorln(err)
		}
		return nil, err
	}
	return client, nil
}

func Run(cfg Config) {
	// init MySQL
	db, err := initMySQL(cfg.MySQLConfig)
	if err != nil {
		logrus.Fatalln(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logrus.Errorln(err)
		}
		logrus.Infoln("connection with mysql closed")
	}()
	logrus.Infoln("connected to mysql")

	// init MongoDB
	client, err := initMongoDB(cfg.MongoConfig)
	if err != nil {
		logrus.Fatalln(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			logrus.Errorln(err)
		}
		logrus.Infoln("connection with mongodb closed")
	}()
	logrus.Infoln("connected to mongo")
	collection := client.Database(cfg.MongoConfig.DBName).Collection(cfg.MongoConfig.CollectionName)

	// init Redis
	redisAddress := fmt.Sprintf("%s:%s", cfg.RedisConfig.Host, cfg.RedisConfig.Port)
	conn, err := redis.Dial("tcp", redisAddress)
	defer func() {
		if err := conn.Close(); err != nil {
			logrus.Errorln(err)
		}
		logrus.Infoln("connection with redis closed")
	}()
	logrus.Infoln("connected to redis")

	// declare app objects
	usersRepo := mysqlrepo.NewUsersRepo(db)
	postsRepo := mongorepo.NewPostsRepo(collection)
	//usersRepo := slicerepo.NewUsersRepo()
	//postsRepo := slicerepo.NewPostsRepo()

	services := service.NewService(usersRepo, postsRepo)

	cookieStorage := cookie.NewRedisStorage(conn)
	sessions := cookie.NewManager(cookieStorage)
	signer := token.NewSigner(cfg.SignerConfig.SigningKey)
	handlers := handler.NewHandler(signer, sessions, services)

	router := handlers.CreateRouter()
	apiAddress := fmt.Sprintf("%s:%s", cfg.ApiConfig.Host, cfg.ApiConfig.Port)
	server := &http.Server{Addr: apiAddress, Handler: router}

	// run application
	go func() {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalln(err)
		}
	}()
	logrus.Infof("http server started: %s", apiAddress)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Infoln("app is shutting down")
	if err = server.Shutdown(context.Background()); err != nil {
		logrus.Errorln(err)
	}
}
