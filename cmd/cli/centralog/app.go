package centralog

import (
	pb "github.com/XiovV/centralog-agent/grpc"
	"github.com/XiovV/centralog-agent/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type App struct {
	centralogClient pb.CentralogClient
	repository      *repository.SQLite
}

func NewApp() *App {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	client := pb.NewCentralogClient(conn)

	repo := repository.New()

	return &App{
		centralogClient: client,
		repository:      repo,
	}
}

func newClient(target string) pb.CentralogClient {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	return pb.NewCentralogClient(conn)
}
