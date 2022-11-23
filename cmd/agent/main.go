package main

import (
	"fmt"
	"github.com/XiovV/centralog-agent/docker"
	"github.com/XiovV/centralog-agent/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	appEnv = os.Getenv("APP_ENV")
)

func main() {
	logger, err := initLogger()
	if err != nil {
		log.Fatalln("couldn't init logger: ", err)
	}

	defer logger.Sync()

	repo := repository.New()
	dockerController := docker.New(repo)

	rateLimiter := NewLimiter(5, 60*time.Second)

	checkAPIKey(repo)

	srv := Server{
		Logger:     logger,
		Docker:     dockerController,
		Repository: repo,
		LogBuffer:  docker.NewLogBuffer(repo),
		RateLimit:  rateLimiter,
	}

	logger.Info("initialising log listener...")
	err = srv.ListenForLogs()
	if err != nil {
		logger.Error("couldn't read config", zap.Error(err))
	}

	logger.Info("server is listening for requests...", zap.String("port", os.Getenv("PORT")))

	srv.Serve()
}

func checkAPIKey(r *repository.SQLite) {
	key := r.GetAPIKey()
	if len(key) == 0 {
		newKey, keyStr := generateHashedAPIKey()

		r.StoreAPIKey(newKey)

		fmt.Println("Your new API key is:", keyStr)
	}
}

func generateHashedAPIKey() ([]byte, string) {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, KEY_LENGTH)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	keyStr := string(s)

	hashedKey, _ := bcrypt.GenerateFromPassword([]byte(keyStr), bcrypt.DefaultCost)

	return hashedKey, keyStr
}

func initLogger() (*zap.Logger, error) {
	if appEnv == LOCAL_ENV || appEnv == STAGING_ENV {
		logger, err := zap.NewDevelopment()

		if err != nil {
			return nil, err
		}

		return logger, nil
	}

	logger, err := zap.NewProduction()

	if err != nil {
		return nil, err
	}

	return logger, nil
}
