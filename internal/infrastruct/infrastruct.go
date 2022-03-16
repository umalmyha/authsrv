package infrastruct

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	valueobj "github.com/umalmyha/authsrv/internal/business/value-object"
	"github.com/umalmyha/authsrv/pkg/database/rdb"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapProductionLogger(service string) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.DisableStacktrace = true
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.CallerKey = "src"
	config.InitialFields = map[string]interface{}{
		"service": service,
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

func NewCliZapLogger() (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.DisableCaller = true
	config.DisableStacktrace = true
	config.Encoding = "console"
	config.EncoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {}
	config.EncoderConfig.EncodeLevel = func(l zapcore.Level, pae zapcore.PrimitiveArrayEncoder) {}
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

func LoadEnv() error {
	if os.Getenv("APP_ENV") != "production" { // TODO: add normal handling later
		return godotenv.Load()
	}
	return nil
}

func JwtConfig() (valueobj.JwtConfig, error) {
	var cfg valueobj.JwtConfig

	privateKeyFile := os.Getenv("AUTHSRV_JWT_PRIVATE_KEY_FILE")
	if privateKeyFile == "" {
		return cfg, errors.New("private key file is not specified")
	}

	privateKey, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return cfg, err
	}

	algorithm := os.Getenv("AUTHSRV_JWT_ALGORITHM")
	issuer := os.Getenv("AUTHSRV_JWT_ISSUER")

	ttlStr := os.Getenv("AUTHSRV_JWT_TTL_MINUTES")
	ttl, err := time.ParseDuration(fmt.Sprintf("%sm", ttlStr))
	if err != nil {
		return cfg, err
	}

	return valueobj.NewJwtConfig(algorithm, issuer, string(privateKey), ttl)
}

func PasswordConfig() (valueobj.PasswordConfig, error) {
	minLengthStr := os.Getenv("AUTHSRV_PASSWORD_MIN_LENGTH")
	if minLengthStr == "" {
		minLengthStr = "0"
	}

	minLength, err := strconv.Atoi(minLengthStr)
	if err != nil {
		return valueobj.PasswordConfig{}, err
	}

	maxLengthStr := os.Getenv("AUTHSRV_PASSWORD_MAX_LENGTH")
	if maxLengthStr == "" {
		maxLengthStr = "0"
	}

	maxLength, err := strconv.Atoi(maxLengthStr)
	if err != nil {
		return valueobj.PasswordConfig{}, err
	}

	hasDigit := false
	if os.Getenv("AUTHSRV_PASSWORD_MUST_HAVE_DIGIT") != "" {
		hasDigit = true
	}

	hasUppercase := false
	if os.Getenv("AUTHSRV_PASSWORD_MUST_HAVE_UPPERCASE") != "" {
		hasUppercase = true
	}

	return valueobj.NewPasswordConfig(minLength, maxLength, hasDigit, hasUppercase)
}

func RefreshTokenConfig() (valueobj.RefreshTokenConfig, error) {
	ttlStr := os.Getenv("AUTHSRV_REFRESH_TOKEN_TTL_HOURS")
	ttl, err := time.ParseDuration(fmt.Sprintf("%sh", ttlStr))
	if err != nil {
		return valueobj.RefreshTokenConfig{}, err
	}

	maxCountStr := os.Getenv("AUTHSRV_REFRESH_TOKEN_MAX_COUNT")
	maxCount, err := strconv.Atoi(maxCountStr)
	if err != nil {
		return valueobj.RefreshTokenConfig{}, err
	}

	cookieName := os.Getenv("AUTHSRV_REFRESH_TOKEN_COOKIE_NAME")

	return valueobj.NewRefreshTokenConfig(ttl, maxCount, cookieName)
}

func ConnectToDb() (*sqlx.DB, error) {
	dbConfig := rdb.NewConfig(
		rdb.DatabasePostgres,
		rdb.WithUser(os.Getenv("AUTHSRV_DB_USERNAME")),
		rdb.WithPassword(os.Getenv("AUTHSRV_DB_PASSWORD")),
		rdb.WithDatabase(os.Getenv("AUTHSRV_DB_DBNAME")),
		rdb.WithHost(os.Getenv("AUTHSRV_DB_HOST")),
		rdb.WithParams(
			rdb.Param("sslmode", "disable"),
		),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := rdb.Connect(ctx, dbConfig)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("database connection error: %s", err.Error()))
	}

	return db, nil
}