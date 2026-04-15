package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

var _ ServerInterface = (*Server)(nil)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

// (GET /health)
func (s *Server) Health(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

// Get list of all jams
// (GET /jams)
func (s *Server) JamsList(c *gin.Context) {}

// Create new jam
// (POST /jams)
func (s *Server) CreateJams(c *gin.Context) {}

// find jam by ID
// (GET /jams/{id})
func (s *Server) FindJamById(c *gin.Context, id openapi_types.UUID) {}

// remove jam participant
// (DELETE /jams/{id}/participants)
func (s *Server) DeleteJamParticipant(c *gin.Context, id openapi_types.UUID) {}

// Get list of all jams
// (GET /jams/{id}/participants)
func (s *Server) JamParticipantList(c *gin.Context, id openapi_types.UUID) {}

// Create new jam participant
// (POST /jams/{id}/participants)
func (s *Server) CreateJamParticipant(c *gin.Context, id openapi_types.UUID) {
}
