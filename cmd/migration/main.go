package main

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"redditclone/internal/app"
	"time"
)

type MigrationConfig struct {
	Attempts int `yaml:"attempts"`
	Timeout  int `yaml:"timeout"`
}

type Config struct {
	app.MySQLConfig `yaml:"mysql"`
	MigrationConfig `yaml:"migration"`
}

func main() {
	ymlFile, err := ioutil.ReadFile("configs/config.yml")
	if err != nil {
		logrus.Fatalln(err)
	}

	var cfg Config
	if err = yaml.Unmarshal(ymlFile, &cfg); err != nil {
		logrus.Fatalln(err)
	}

	if err = godotenv.Load(".env"); err != nil {
		logrus.Fatalln(err)
	}

	cfg.Host = os.Getenv("MYSQL_HOST")
	cfg.Port = os.Getenv("MYSQL_PORT")
	cfg.Username = os.Getenv("MYSQL_USER")
	cfg.Password = os.Getenv("MYSQL_PASSWORD")
	cfg.DBName = os.Getenv("MYSQL_DATABASE")

	// migrate -path ./migrations/ -database 'mysql://user:password@tcp(localhost:3306)/redditclone' up 1
	databaseURL := fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	var (
		attempts  = cfg.Attempts
		migration *migrate.Migrate
	)

	for attempts > 0 {
		migration, err = migrate.New("file://migrations", databaseURL)
		if err == nil {
			break
		}

		logrus.Infof("migrate: mysql is trying to connect, attempts left: %d", attempts)
		time.Sleep(time.Duration(cfg.Timeout) * time.Second)
		attempts--
	}

	if err != nil {
		logrus.Fatalf("migrate: mysql connect error: %s", err)
	}

	err = migration.Up()
	defer migration.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logrus.Fatalf("migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		logrus.Infof("migrate: no change")
		return
	}

	logrus.Infof("migrate: up success")
}
