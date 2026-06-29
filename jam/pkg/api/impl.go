package api

import (
	"jam/internal/application"
	"net/http"

	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

var _ ServerInterface = (*Server)(nil)

type Server struct {
	jamService *application.JamService
}

func NewServer(jamService *application.JamService) *Server {
	return &Server{
		jamService: jamService,
	}
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
func (s *Server) CreateJams(c *gin.Context) {
	var createJam CreateJamsJSONRequestBody
	if err := c.BindJSON(&createJam); err != nil {
		return
	}

	params := application.CreateJamParams{
		Name:      createJam.Name,
		Location:  createJam.Location,
		Date:      createJam.Date,
		CreatedBy: createJam.CreatedBy.String(),
	}
	err := s.jamService.CreateJam(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, nil)
}

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
