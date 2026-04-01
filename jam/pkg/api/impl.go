package api

import (
	"github.com/gin-gonic/gin"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

// (GET /health)
func (s *Server) GetHealth(c *gin.Context) {

}
