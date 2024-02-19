package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"main.go/internal/entity"
)

func (h *Handler) createItem(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "user userID not found")
		return
	}

	listID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid id parameter")
		return
	}

	var input entity.TimeslotItem

	if err = ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	itemID, err := h.services.TimeslotItem.Create(userID, listID, input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": itemID,
	})
}

func (h *Handler) getAllItems(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "user userID not found")
		return
	}

	listID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid list id parameter")
		return
	}

	items, err := h.services.TimeslotItem.GetAll(userID, listID)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, items)
}

func (h *Handler) getItemByID(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "user userID not found")
		return
	}

	itemID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid item id parameter")
		return
	}

	item, err := h.services.TimeslotItem.GetByID(userID, itemID)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (h *Handler) updateItem(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		return
	}

	itemID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid id parameter")
		return
	}

	var input entity.UpdateItemInput
	if err = ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.TimeslotItem.Update(userID, itemID, input); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) deleteItem(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "user userID not found")
		return
	}

	itemID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid item id parameter")
		return
	}

	err = h.services.TimeslotItem.Delete(userID, itemID)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
