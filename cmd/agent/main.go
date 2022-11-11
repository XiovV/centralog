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

	checkAPIKey(repo)

	//containers := []string{"logserver1", "logserver2"}
	//
	//logBuffer := docker.NewLogBuffer(repo)
	//
	//for _, container := range containers {
	//	logWriter := docker.NewBackgroundLogWriter(logBuffer, container)
	//
	//	go dockerController.CollectLogs(container, logWriter, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Timestamps: true, Since: "0m"})
	//}

	srv := Server{
		Logger:     logger,
		Docker:     dockerController,
		Repository: repo,
	}

	logger.Info("running...", zap.String("port", os.Getenv("PORT")))

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
