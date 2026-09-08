package handlers

import (
	"context"
	"goflagsmith/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DecisionMaker interface {
	Route(ctx context.Context, req domain.UserContext) domain.RouteDecision
}

type RouteHandler struct {
	dm DecisionMaker
}

func NewRouteHandler(dm DecisionMaker) *RouteHandler {
	return &RouteHandler{dm: dm}
}

func (rh *RouteHandler) Handle(c *gin.Context) {

	var req domain.UserContext
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid request payload: " + err.Error()},
		)
		return
	}

	res := rh.dm.Route(c.Request.Context(), req)

	c.JSON(http.StatusOK, res)
}
