package handler

import (
	"effectiveMobile/pkg/db"
	"effectiveMobile/pkg/domain/people"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// @Summary Get info about a person
// @Description Get info about a person by passport series and number
// @Tags People
// @Produce  json
// @Param passportSerie query string true "Passport Series"
// @Param passportNumber query string true "Passport Number"
// @Success 200 {object} people.Request
// @Failure 400 gin.H map[string]string
// @Failure 404 gin.H map[string]string
// @Failure 500 gin.H map[string]string
// @Router /info [get]
func (h *Handler) InfoPeople(c *gin.Context) {
	passportSerie := c.Query("passportSerie")
	passportNumber := c.Query("passportNumber")
	result, err := h.service.InfoPeople(c.Request.Context(), passportSerie, passportNumber)
	if err != nil {
		switch err.Error() {
		case db.ErrValidate.Error():
			c.JSON(400, gin.H{"message": "params error validate serie or number"})
			log.Error(err.Error())
			break
		case db.ErrPassportSerie.Error():
			c.JSON(400, gin.H{"error": err.Error()})
			log.Error(err.Error())
			break
		case db.ErrPassportNumber.Error():
			c.JSON(400, gin.H{"error": err.Error()})
			log.Error(err.Error())
			break
		case db.ErrNotExist.Error():
			c.JSON(404, gin.H{"error": err.Error()})
			log.Error(err.Error())
			break
		default:
			c.JSON(500, gin.H{"error": err.Error()})
			log.Error(err.Error())
		}
		return
	}

	c.JSON(200, result)
	log.Info("Success InfoPeople %v", result)
	return
}

// @Summary Get list of people
// @Description Get list of people with optional filters and pagination
// @Tags People
// @Produce  json
// @Param filter query people.Filter false "Filter parameters"
// @Param pagination query people.Pagination false "Pagination parameters"
// @Success 200 {array} people.People
// @Failure 400 gin.H map[string]string
// @Router /people [get]
func (h *Handler) GetPeople(c *gin.Context) {
	var filter people.Filter
	if err := c.BindQuery(&filter); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Error(err.Error())
		return
	}

	var pagination people.Pagination
	if err := c.BindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error bind Get pagination": err.Error()})
		log.Error(err.Error())
		return
	}

	peoples, err := h.service.GetPeople(c.Request.Context(), &filter, &pagination)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Error(err.Error())
		return
	}
	c.JSON(200, peoples)
	log.Info("Success GetPeople %v", peoples)
	return
}

// @Summary Update a person
// @Description Update person information
// @Tags People
// @Accept  json
// @Produce  json
// @Param people body people.Info true "Update person info"
// @Success 200 {object} people.Info
// @Failure 400 gin.H map[string]string
// @Failure 401 gin.H map[string]string
// @Failure 409 gin.H map[string]string
// @Failure 500 gin.H map[string]string
// @Router /people [put]
func (h *Handler) PutPeople(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(401, gin.H{"error": "Authorization required"})
		log.Errorf("Auth error: Authorization required")
		return
	}
	id, ok := userId.(string)
	if !ok {
		c.JSON(401, gin.H{"error": "Invalid user ID"})
		log.Errorf("Invalid user ID")
		return
	}

	var updatePeople people.Info
	if err := c.BindJSON(&updatePeople); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	result, err := h.service.PutPeople(c.Request.Context(), id, updatePeople)
	if err != nil {
		switch err.Error() {
		case db.ErrParamNotFound.Error():
			c.JSON(http.StatusBadRequest, gin.H{"error": "Account ID is required"})
			log.Error(err.Error())
		case db.ErrDuplicate.Error():
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			log.Error(err.Error())
		case db.ErrUpdateFailed.Error():
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update people"})
			log.Error(err.Error())
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Error(err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, result)
	log.Info("Success PutPeople %v", result)
	return
}

// @Summary Delete a person
// @Description Delete person information
// @Tags People
// @Success 200 gin.H map[string]string
// @Failure 401 gin.H map[string]string
// @Failure 403 gin.H map[string]string
// @Failure 500 gin.H map[string]string
// @Router /people [delete]
func (h *Handler) DeletePeople(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(401, gin.H{"error": "Authorization required"})
		log.Errorf("Auth error: Authorization required")
		return
	}
	id, ok := userId.(string)
	if !ok {
		c.JSON(401, gin.H{"error": "Invalid user ID"})
		log.Errorf("Invalid user ID")
		return
	}

	err := h.service.DeletePeople(c.Request.Context(), id)
	if err != nil {
		switch err.Error() {
		case db.ErrParamNotFound.Error():
			c.JSON(400, gin.H{"error": "problems with param"})
			log.Error(err.Error())
			break
		case db.ErrDeleteFailed.Error():
			c.JSON(403, gin.H{"error": err.Error()})
			log.Error(err.Error())
			break
		default:
			c.JSON(500, gin.H{"error": err.Error()})
			log.Error(err.Error())
		}
		return
	}
	c.JSON(200, gin.H{"id": id})
	log.Printf("Success DeletePeople %v", id)
	return
}
