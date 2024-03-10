package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"main.go/internal/entity"
)

//go:generate mockgen -source=item.go -destination=mocks/itemMock.go
type TimeslotItemService interface {
	Create(userID, listID int, input entity.TimeslotItem) (int, error)
	GetAll(userID, listID int) ([]entity.TimeslotItem, error)
	GetByID(userID, itemID int) (entity.TimeslotItem, error)
	Delete(userID, itemID int) error
	Update(userID, itemID int, input entity.UpdateItemInput) error
	GetByRange(input entity.ItemsByRange) ([]entity.TimeslotItemWithUsername, error)
}

type TimeslotItemHandler struct {
	service TimeslotItemService
}

func NewTimeslotItemHandler(service TimeslotItemService) *TimeslotItemHandler {
	return &TimeslotItemHandler{service: service}
}

// @Summary Create timeslot item
// @Security ApiKeyAuth
// @Tags items
// @Description create timeslot item
// @ID create-item
// @Accept  json
// @Produce  json
// @Param input body entity.TimeslotItem true "item info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/items [post].
func (h *TimeslotItemHandler) createItem(ctx *gin.Context) {
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

	itemID, err := h.service.Create(userID, listID, input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": itemID,
	})
}

type getAllItemsResponse struct {
	Data []entity.TimeslotItem `json:"data"`
}

// @Summary Get All Items
// @Security ApiKeyAuth
// @Tags items
// @Description get all items
// @ID get-all-items
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllItemsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/items [get].
func (h *TimeslotItemHandler) getAllItems(ctx *gin.Context) {
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

	items, err := h.service.GetAll(userID, listID)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, getAllItemsResponse{Data: items})
}

// @Summary Get Item By Id
// @Security ApiKeyAuth
// @Tags items
// @Description get item by id
// @ID get-item-by-id
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.TimeslotItem
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/items/:id [get].
func (h *TimeslotItemHandler) getItemByID(ctx *gin.Context) {
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

	item, err := h.service.GetByID(userID, itemID)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, item)
}

type getItemsByRangeResponse struct {
	Data []entity.TimeslotItemWithUsername `json:"data"`
}

// @Summary Get Item By Id
// @Security ApiKeyAuth
// @Tags items
// @Description get item by id
// @ID get-item-by-id
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.TimeslotItem
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/items/:id [get].
func (h *TimeslotItemHandler) getItemsByRange(ctx *gin.Context) {
	_, err := getUserID(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "user userID not found")
		return
	}

	var input entity.ItemsByRange

	if err = ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	items, err := h.service.GetByRange(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, getItemsByRangeResponse{Data: items})
}

// @Summary Update Item
// @Security ApiKeyAuth
// @Tags items
// @Description update list
// @ID update-list
// @Accept  json
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/items/:id [put].
func (h *TimeslotItemHandler) updateItem(ctx *gin.Context) {
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

	if err = h.service.Update(userID, itemID, input); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, statusResponse{"ok"})
}

// @Summary Delete Item
// @Security ApiKeyAuth
// @Tags imens
// @Description update item
// @ID delete-item
// @Accept  json
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/items/:id [delete].
func (h *TimeslotItemHandler) deleteItem(ctx *gin.Context) {
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

	err = h.service.Delete(userID, itemID)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
