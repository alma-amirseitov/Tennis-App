package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// Server
	Port        int    `envconfig:"PORT" default:"8080"`
	Environment string `envconfig:"ENVIRONMENT" default:"development"`

	// Database
	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`

	// Redis
	RedisURL string `envconfig:"REDIS_URL" required:"true"`

	// JWT
	JWTSecret     string        `envconfig:"JWT_SECRET" required:"true"`
	JWTAccessTTL  time.Duration `envconfig:"JWT_ACCESS_TTL" default:"15m"`
	JWTRefreshTTL time.Duration `envconfig:"JWT_REFRESH_TTL" default:"720h"`

	// SMS
	SMSProvider    string `envconfig:"SMS_PROVIDER" default:"mock"`
	SMSAPIURL      string `envconfig:"SMS_API_URL"`
	SMSAPILogin    string `envconfig:"SMS_API_LOGIN"`
	SMSAPIPassword string `envconfig:"SMS_API_PASSWORD"`

	// Storage (S3/MinIO)
	S3Endpoint  string `envconfig:"S3_ENDPOINT"`
	S3AccessKey string `envconfig:"S3_ACCESS_KEY"`
	S3SecretKey string `envconfig:"S3_SECRET_KEY"`
	S3Bucket    string `envconfig:"S3_BUCKET" default:"tennisapp"`
	S3PublicURL string `envconfig:"S3_PUBLIC_URL"`

	// Firebase
	FirebaseCredentials string `envconfig:"FIREBASE_CREDENTIALS"`

	// Sentry
	SentryDSN string `envconfig:"SENTRY_DSN"`

	// Dev
	DevOTPCode string `envconfig:"DEV_OTP_CODE"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	return &cfg, nil
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func (c *Config) Addr() string {
	return fmt.Sprintf(":%d", c.Port)
}
