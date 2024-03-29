package centralog

import (
	"github.com/AlecAivazis/survey/v2"
	pb "github.com/XiovV/centralog-agent/grpc"
	"github.com/XiovV/centralog-agent/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	centralogClient pb.CentralogClient
	repository      *repository.SQLite
}

func NewApp() *App {
	return &App{
		repository: repository.New(),
	}
}

func (a *App) initClient(target string) error {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	a.centralogClient = pb.NewCentralogClient(conn)

	return nil
}

func (a *App) confirmationPrompt(text string) bool {
	var confirmation bool
	prompt := &survey.Confirm{
		Message: text,
	}
	survey.AskOne(prompt, &confirmation)

	return confirmation
}
