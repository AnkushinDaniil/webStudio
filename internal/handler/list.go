package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/internal/entity"
)

func (h *Handler) createList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "user userId not found")
		return
	}

	var input entity.TimeslotList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TimeslotList.Create(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllLists(c *gin.Context) {
}

func (h *Handler) getListBuId(c *gin.Context) {
}

func (h *Handler) updateList(c *gin.Context) {
}

func (h *Handler) deleteList(c *gin.Context) {
}
