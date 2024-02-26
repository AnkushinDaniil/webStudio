package rest

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userID"
)

func (h *Handlers) userIdentity(ctx *gin.Context) {
	header := ctx.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(ctx, http.StatusUnauthorized, "empty authorization header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(ctx, http.StatusUnauthorized, "invalid authorization header")
		return
	}

	if len(headerParts[1]) == 0 {
		newErrorResponse(ctx, http.StatusUnauthorized, "token is empty")
		return
	}
	userID, err := h.AuthorizationHandler.service.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
	}

	ctx.Set(userCtx, userID)
}

func getUserID(ctx *gin.Context) (int, error) {
	userID, exists := ctx.Get(userCtx)
	if !exists {
		newErrorResponse(ctx, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}

	idInt, exists := userID.(int)
	if !exists {
		newErrorResponse(ctx, http.StatusInternalServerError, "the user ID is of an invalid type")
		return 0, errors.New("the user ID is of an invalid type")
	}
	return idInt, nil
}
