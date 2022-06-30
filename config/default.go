package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	BaseUrl string `mapstructure:"BASE_URL"`
	DBUri   string `mapstructure:"DATABASE_URL"`
	Port    string `mapstructure:"PORT"`
	Domain  string `mapstructure:"DOMAIN"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`

	EmailFrom string `mapstructure:"EMAIL_FROM"`
	SMTPHost  string `mapstructure:"SMTP_HOST"`
	SMTPPass  string `mapstructure:"SMTP_PASS"`
	SMTPPort  int    `mapstructure:"SMTP_PORT"`
	SMTPUser  string `mapstructure:"SMTP_USER"`
}

func LoadConfig() (config Config, err error) {
	config.BaseUrl = os.Getenv("BASE_URL")
	config.DBUri = os.Getenv("DATABASE_URL")
	config.Port = os.Getenv("PORT")
	config.Domain = os.Getenv("DOMAIN")

	config.AccessTokenPrivateKey = os.Getenv("ACCESS_TOKEN_PRIVATE_KEY")
	config.AccessTokenPublicKey = os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")
	config.RefreshTokenPrivateKey = os.Getenv("REFRESH_TOKEN_PRIVATE_KEY")
	config.RefreshTokenPublicKey = os.Getenv("REFRESH_TOKEN_PUBLIC_KEY")
	config.AccessTokenExpiresIn, err = time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRED_IN"))
	config.RefreshTokenExpiresIn, err = time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXPIRED_IN"))
	config.AccessTokenMaxAge, err = strconv.Atoi(os.Getenv("ACCESS_TOKEN_MAXAGE"))
	config.RefreshTokenMaxAge, err = strconv.Atoi(os.Getenv("REFRESH_TOKEN_MAXAGE"))

	config.EmailFrom = os.Getenv("EMAIL_FROM")
	config.SMTPHost = os.Getenv("SMTP_HOST")
	config.SMTPPass = os.Getenv("SMTP_PASS")
	config.SMTPPort, err = strconv.Atoi(os.Getenv("SMTP_PORT"))
	config.SMTPUser = os.Getenv("SMTP_USER")

	return
}
