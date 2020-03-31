package app

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/oauth2"
)

type Config interface {
	Environment() EnvMode
	Database() *Database
	GitHub() *GitHub
	Logger() *zap.SugaredLogger
	Addr() string
}

type EnvMode string

const (
	Production  EnvMode = "production"
	Development EnvMode = "development"
	Test        EnvMode = "test"
)

type Database struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type GitHub struct {
	AccessToken   string
	OAuthClientID string
	OAuthSecret   string
	Client        *github.Client
}

func (db *Database) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.Name)

}

type Env struct {
	mode   EnvMode
	logger *zap.SugaredLogger
	db     *Database
	gh     *GitHub
}

var _ Config = (*Env)(nil)

func (e *Env) Environment() EnvMode {
	return e.mode
}

func (e *Env) Logger() *zap.SugaredLogger {
	return e.logger
}

func (e *Env) Database() *Database {
	return e.db
}

func (e *Env) GitHub() *GitHub {
	return e.gh
}

func (e *Env) Addr() string {
	return ":5000"
}

func FromEnv(deploy string) (*Env, error) {

	// load from .env but skip if not found
	if err := godotenv.Load(); err != nil {
		fmt.Fprintf(os.Stdout, "SKIP: loading .env failed: %s", err)
	}

	mode := Environment()
	var err error

	var log *zap.SugaredLogger
	if log, err = initLogger(mode); err != nil {
		return nil, err
	}

	log.With("name", "app").Infof("in %q mode ", mode)

	env := &Env{mode: mode, logger: log}

	if env.db, err = initDB(); err != nil {
		return nil, err
	}

	// fetch only for api deployment
	if deploy == "api" {
		if env.gh, err = initGithub(); err != nil {
			return nil, err
		}
	}

	return env, nil
}

func env(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("NO %q environment variable defined", key)
	}
	return val, nil
}

func Environment() EnvMode {
	mode := "production"
	if val, ok := os.LookupEnv("ENVIRONMENT"); ok {
		mode = val
	}

	switch strings.ToLower(mode) {
	case "development":
		return Development
	case "test":
		return Test
	default:
		return Production
	}
}

func initDB() (*Database, error) {
	var err error

	db := &Database{}
	if db.Host, err = env("POSTGRESQL_HOST"); err != nil {
		return nil, err
	}
	if db.Port, err = env("POSTGRESQL_PORT"); err != nil {
		return nil, err
	}
	if db.Name, err = env("POSTGRESQL_DATABASE"); err != nil {
		return nil, err
	}

	if db.User, err = env("POSTGRESQL_USER"); err != nil {
		return nil, err
	}
	if db.Password, err = env("POSTGRESQL_PASSWORD"); err != nil {
		return nil, err
	}

	return db, nil
}

func initGithub() (*GitHub, error) {
	var err error
	gh := &GitHub{}
	if gh.AccessToken, err = env("GITHUB_TOKEN"); err != nil {
		return nil, err
	}
	if gh.OAuthClientID, err = env("CLIENT_ID"); err != nil {
		return nil, err
	}
	if gh.OAuthSecret, err = env("CLIENT_SECRET"); err != nil {
		return nil, err
	}

	token := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: gh.AccessToken})
	client := oauth2.NewClient(context.Background(), token)
	gh.Client = github.NewClient(client)
	return gh, nil
}

func initLogger(mode EnvMode) (*zap.SugaredLogger, error) {

	var log *zap.Logger
	var err error

	switch mode {
	case Production:
		log, err = zap.NewProduction()

	default:
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		log, err = config.Build()
	}

	if err != nil {
		return nil, err
	}
	return log.Sugar(), nil
}
