package api

import (
	app "battle/internal/App"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type Server struct {
	battleService app.BattleService
}

func NewServer(battleService *app.BattleService) ServerInterface {
	return &Server{}
}

// Create battle
// (POST /battles)
func (s *Server) CreateBattle(c *gin.Context) {
	var createBattle CreateBattleJSONBody
	if err := c.BindJSON(&createBattle); err != nil {
		return
	}

	params := app.CreateBattleData{
		EventID:   createBattle.EventId,
		Dancer1ID: createBattle.Dancer1Id,
		Dancer2ID: createBattle.Dancer2Id,
	}
	err := s.battleService.CreateBattle(c.Request.Context(), params)
	if err != nil {
		responseError := Error{
			Error: err.Error(),
		}
		c.JSON(s.httpCodeFromError(err), responseError)
		return
	}
	c.Status(http.StatusCreated)
}

// Finish battle
// (POST /battles/{id}/finish)
func (s *Server) FinishBattle(c *gin.Context, id string) {
	result, err := s.battleService.FinishBattle(id)
	if err != nil {
		responseError := Error{
			Error: err.Error(),
		}
		c.JSON(s.httpCodeFromError(err), responseError)
		return
	}
	response := BattleResult{
		BattleId:          &result.BattleID,
		WinnerId:          result.WinnerID,
		Dancer1TotalScore: &result.Dancer1TotalScore,
		Dancer2TotalScore: &result.Dander2TotalScore,
		Status:            (*BattleResultStatus)(&result.Status),
		FinishedAt:        result.FinishedAt,
	}
	c.JSON(http.StatusOK, response)
}

// Judge submit score
// (POST /battles/{id}/scores)
func (s *Server) SubmitScores(c *gin.Context, id string) {

}

// Get balltes
// (GET /events/{id}/battles)
func (s *Server) GetEventBattles(c *gin.Context, id string) {}

// Health check
// (GET /health)
func (s *Server) Health(c *gin.Context) {}

func (s *Server) httpCodeFromError(err error) int {
	switch {
	case errors.Is(err, app.ErrEventNotFound):
		return http.StatusNotFound
	case errors.Is(err, app.ErrDancerNoDancerRole):
	case errors.Is(err, app.ErrNoSubmittedBattle):
		return http.StatusUnprocessableEntity
	case errors.Is(err, app.ErrDanceMultipleBattle):
		return http.StatusConflict
	}
	return http.StatusInternalServerError
}

func (s *Server) uuidFromString(id string) openapi_types.UUID {
	uuid, _ := uuid.FromBytes([]byte(id))
	return uuid
}
