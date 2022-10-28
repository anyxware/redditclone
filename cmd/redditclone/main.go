package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"redditclone/internal/app"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	// disable http logger
	log.SetOutput(ioutil.Discard)

	ymlFile, err := ioutil.ReadFile("configs/config.yml")
	if err != nil {
		logrus.Fatalln(err)
	}

	var cfg app.Config
	if err = yaml.Unmarshal(ymlFile, &cfg); err != nil {
		logrus.Fatalln(err)
	}

	if err = godotenv.Load(".env"); err != nil {
		logrus.Fatalln(err)
	}

	cfg.ApiConfig.Host = os.Getenv("API_HOST")
	cfg.ApiConfig.Port = os.Getenv("API_PORT")

	cfg.MySQLConfig.Host = os.Getenv("MYSQL_HOST")
	cfg.MySQLConfig.Port = os.Getenv("MYSQL_PORT")
	cfg.MySQLConfig.Username = os.Getenv("MYSQL_USER")
	cfg.MySQLConfig.Password = os.Getenv("MYSQL_PASSWORD")
	cfg.MySQLConfig.DBName = os.Getenv("MYSQL_DATABASE")

	cfg.MongoConfig.Host = os.Getenv("MONGO_HOST")
	cfg.MongoConfig.Port = os.Getenv("MONGO_PORT")
	cfg.MongoConfig.Username = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	cfg.MongoConfig.Password = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	cfg.MongoConfig.DBName = os.Getenv("MONGO_DATABASE")

	cfg.RedisConfig.Host = os.Getenv("REDIS_HOST")
	cfg.RedisConfig.Port = os.Getenv("REDIS_PORT")

	cfg.SignerConfig.SigningKey = os.Getenv("SIGNING_KEY")

	app.Run(cfg)
}
