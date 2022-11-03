package server

import (
	"github.com/XiovV/centralog-agent/pkg/docker"
	"github.com/XiovV/centralog-agent/repository"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	LOCAL_ENV      = "LOCAL"
	STAGING_ENV    = "STAGING"
	PRODUCTION_ENV = "PROD"
)

type Server struct {
	Logger     *zap.Logger
	Docker     *docker.Controller
	Repository *repository.Repository
}

func (s *Server) Serve() *gin.Engine {

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
