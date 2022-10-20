package server

import (
	"fmt"
	"github.com/XiovV/centralog-agent/pkg/docker"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	LOCAL_ENV      = "LOCAL"
	STAGING_ENV    = "STAGING"
	PRODUCTION_ENV = "PROD"
)

type Server struct {
	Logger *zap.Logger
	Docker *docker.Controller
}

type Logs struct {
}

func (l Logs) Write(p []byte) (int, error) {
	fmt.Println(p)

	return 0, nil
}

func (s *Server) New() {
	out := s.Docker.GetLogs("d64552a3f96c", types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})

	stdcopy.StdCopy(docker.StdOut, docker.StdOut, out)
}

func (s *Server) Serve() *gin.Engine {
	s.New()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), s.CORS())

	v1 := router.Group("/v1")
	nodeGroup := v1.Group("/nodes")
	nodeGroup.Use(s.Auth())
	{
		nodeGroup.GET("/info", s.getNodeInfo)
	}

	return router
}
