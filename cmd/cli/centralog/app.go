package centralog

import (
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

func (a *App) newClient(target string) (pb.CentralogClient, error) {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return pb.NewCentralogClient(conn), nil
}
