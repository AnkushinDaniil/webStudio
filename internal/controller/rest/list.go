package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"main.go/internal/entity"
)

//go:generate mockgen -source=list.go -destination=mocks/listMock.go
type TimeslotListService interface {
	Create(userID int, list entity.TimeslotsList) (int, error)
	GetAll(userID int) ([]entity.TimeslotsList, error)
	GetByID(userID, listID int) (entity.TimeslotsList, error)
	Delete(userID, listID int) error
	Update(userID, listID int, input entity.UpdateListInput) error
}

type TimeslotListHandler struct {
	service TimeslotListService
}

func NewTimeslotListHandler(service TimeslotListService) *TimeslotListHandler {
	return &TimeslotListHandler{service: service}
}

// @Summary Create timeslot list
// @Security ApiKeyAuth
// @Tags lists
// @Description create timeslot list
// @ID create-list
// @Accept  json
// @Produce  json
// @Param input body entity.TimeslotsList true "list info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists [post].
func (h *TimeslotListHandler) createList(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "user userID not found")
		return
	}

	var input entity.TimeslotsList
	if err = ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	listID, err := h.service.Create(userID, input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": listID,
	})
}

type getAllListsResponse struct {
	Data []entity.TimeslotsList `json:"data"`
}

// @Summary Get All Lists
// @Security ApiKeyAuth
// @Tags lists
// @Description get all lists
// @ID get-all-lists
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllListsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists [get].
func (h *TimeslotListHandler) getAllLists(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "user userID not found")
		return
	}

	lists, err := h.service.GetAll(userID)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, getAllListsResponse{Data: lists})
}

// @Summary Get List By Id
// @Security ApiKeyAuth
// @Tags lists
// @Description get list by id
// @ID get-list-by-id
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.ListsItem
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/:id [get].
func (h *TimeslotListHandler) getListByID(ctx *gin.Context) {
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

	list, err := h.service.GetByID(userID, listID)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, list)
}

// @Summary Update list
// @Security ApiKeyAuth
// @Tags lists
// @Description update list
// @ID update-list
// @Accept  json
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/:id [put].
func (h *TimeslotListHandler) updateList(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		return
	}

	listID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid id parameter")
		return
	}

	var input entity.UpdateListInput
	if err = ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.service.Update(userID, listID, input); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, statusResponse{"ok"})
}

// @Summary Delete list
// @Security ApiKeyAuth
// @Tags lists
// @Description update list
// @ID delete-list
// @Accept  json
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/:id [delete].
func (h *TimeslotListHandler) deleteList(ctx *gin.Context) {
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

	err = h.service.Delete(userID, listID)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
