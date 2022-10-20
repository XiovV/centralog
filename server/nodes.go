package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) getNodeInfo(c *gin.Context) {
	var response struct {
		Instances []string `json:"instance_name"`
	}

	response.Instances = []string{"instance1", "instance2", "instance3"}

	c.JSON(http.StatusOK, response)
}
