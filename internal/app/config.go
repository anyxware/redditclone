package app

type SignerConfig struct {
	SigningKey string
}

type ApiConfig struct {
	Host string `yaml:"-"`
	Port string `yaml:"-"`
}

type MySQLConfig struct {
	Host               string `yaml:"-"`
	Port               string `yaml:"-"`
	Username           string `yaml:"-"`
	Password           string `yaml:"-"`
	DBName             string `yaml:"-"`
	MaxOpenConnections int    `yaml:"max_open_connections"`
}

type MongoConfig struct {
	Host           string `yaml:"-"`
	Port           string `yaml:"-"`
	Username       string `yaml:"-"`
	Password       string `yaml:"-"`
	DBName         string `yaml:"dbname"`
	CollectionName string `yaml:"collection_name"`
}

type RedisConfig struct {
	Host string `yaml:"-"`
	Port string `yaml:"-"`
}

type Config struct {
	ApiConfig    ApiConfig    `yaml:"api"`
	MySQLConfig  MySQLConfig  `yaml:"mysql"`
	MongoConfig  MongoConfig  `yaml:"mongo"`
	RedisConfig  RedisConfig  `yaml:"redis"`
	SignerConfig SignerConfig `yaml:"-"`
}
