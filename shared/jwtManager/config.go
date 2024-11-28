package jwtManager

import (
	"flag"
	"os"
	"time"

	"github.com/ACamaraLara/K8sBlockChainDemo/shared/config"
)

const (
	DefaultExpirationTimeHours int    = 4
	DeffaultJWTIssuer          string = "account-service"
)

type Config struct {
	SecretKey         string
	AccessTokenExpiry time.Duration
	Issuer            string
}

func (cfg *Config) AddFlagsParams() {
	var timeExpHours int
	flag.IntVar(&timeExpHours, "jwt-expiration",
		config.GetEnvironIntWithDefault("JWT_EXPIRATION_HOURS", DefaultExpirationTimeHours), "JWT expiration in hours (JWT_EXPIRATION_HOURS).")
	cfg.AccessTokenExpiry = time.Duration(timeExpHours) * time.Hour
	flag.StringVar(&cfg.Issuer, "jwt-issuer", config.GetEnvironWithDefault("JWT_ISSUER", DeffaultJWTIssuer), "JWT service issuer (JWT_ISSUER).")
	// Secret key doesn't have to be deffault value written in code. it has to be configured in service configuration.
	secretKey := os.Getenv("JWT_SECRET_KEY")
	cfg.SecretKey = secretKey
}
