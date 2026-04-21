package api

import (
	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

// Create battle
// (POST /battles)
func (s *Server) CreateBattle(c *gin.Context) {}

// Finish battle
// (POST /battles/{id}/finish)
func (s *Server) FinishBattle(c *gin.Context, id openapi_types.UUID) {}

// Judge submit score
// (POST /battles/{id}/scores)
func (s *Server) SubmitScores(c *gin.Context, id openapi_types.UUID) {}

// Get balltes
// (GET /events/{id}/battles)
func (s *Server) GetEventBattles(c *gin.Context, id openapi_types.UUID) {}

// Health check
// (GET /health)
func (s *Server) Health(c *gin.Context) {}
