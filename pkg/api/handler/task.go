package handler

import (
	"effectiveMobile/pkg/db"
	"effectiveMobile/pkg/domain/task"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// @Summary Start a task for a person
// @Description Start a task for a person
// @Tags Tasks
// @Accept  json
// @Produce  json
// @Param task body task.Task true "Task info"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /people/task/start [post]
func (h *Handler) StartTask(c *gin.Context) {
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

	var currTask task.Task
	if err := c.BindJSON(&currTask); err != nil {
		c.JSON(400, gin.H{"error bind Task people": err.Error()})
		log.Error("error bind Task people %v", err.Error())
		return
	}

	result, err := h.service.TaskStart(c.Request.Context(), id, currTask)
	if err != nil {
		switch err.Error() {
		case db.ErrDuplicate.Error():
			c.JSON(409, "Email already exist")
			log.Error("Register email failed %v", err.Error())
			break
		default:
			c.JSON(500, gin.H{"error Task people": err.Error()})
			log.Error("error service Task people %v", err.Error())
		}
		return
	}

	c.JSON(201, gin.H{"data": result})
	log.Info("Success registration: %v", result)
}

// @Summary Finish a task for a person
// @Description Finish a task for a person
// @Tags Tasks
// @Produce  json
// @Param taskId path string true "Task ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /people/task/finish/{taskId} [post]
func (h *Handler) FinishTask(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(401, gin.H{"error": "Authorization required"})
		log.Errorf("Auth error: Authorization required")
		return
	}
	_, ok := userId.(string)
	if !ok {
		c.JSON(401, gin.H{"error": "Invalid user ID"})
		log.Errorf("Invalid user ID")
		return
	}

	taskId := c.Param("taskId")
	result, err := h.service.TaskFinish(c.Request.Context(), taskId)
	if err != nil {
		switch err.Error() {
		case db.ErrParamNotFound.Error():
			c.JSON(403, gin.H{"error": "Task not found"})
			log.Error("Register task failed %v", err.Error())
			return
		case db.ErrUpdateFailed.Error():
			c.JSON(400, gin.H{"error": "Update failed"})
			log.Error(err.Error())
			return
		default:
			c.JSON(500, gin.H{"error": "Internal server error on FinishTask"})
			log.Error(err.Error())
			return
		}
	}
	c.JSON(200, gin.H{"data": result})
	log.Info("Success finish task: %v", result)
}

// @Summary Get tasks for a person
// @Description Get tasks for a person within a time range
// @Tags Tasks
// @Produce  json
// @Param startTime query string true "Start Time"
// @Param endTime query string true "End Time"
// @Success 200 {array} task.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /people/task/ [get]
func (h *Handler) GetTask(c *gin.Context) {
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	result, err := h.service.GetTask(c.Request.Context(), startTime, endTime)
	if err != nil {
		switch err.Error() {
		case db.ErrParamNotFound.Error():
			c.JSON(403, gin.H{"error": "Task not found"})
			log.Error("Register task failed %v", err.Error())
			return
		case db.ErrNotExist.Error():
			c.JSON(404, gin.H{"error": "Task not found"})
			log.Error("Get task %v", err.Error())
			return
		default:
			c.JSON(500, gin.H{"error Task people": err.Error()})
			log.Error("error service Task people %v", err.Error())
			return
		}
	}
	c.JSON(200, gin.H{"data": result})
	log.Info("Success get task: %v", result)
	return
}

// @Summary Get list of tasks
// @Description Get list of all tasks
// @Tags Tasks
// @Produce  json
// @Success 200 {array} task.Task
// @Failure 500 {object} map[string]string
// @Router /tasks [get]
func (h *Handler) GetAllTask(c *gin.Context) {
	result, err := h.service.GetAllTask(c.Request.Context())
	if err != nil {
		switch err.Error() {
		case db.ErrParamNotFound.Error():
			c.JSON(403, gin.H{"error": "Tasks not found"})
			log.Error("Register task failed %v", err.Error())
			return
		case db.ErrNotExist.Error():
			c.JSON(404, gin.H{"error": "Tasks not found"})
			log.Error("Get task %v", err.Error())
			return
		default:
			c.JSON(500, gin.H{"error Tasks people": err.Error()})
			log.Error("error service Tasks people %v", err.Error())
			return
		}
	}
	c.JSON(200, gin.H{"data": result})
	log.Info("Success get tasks: %v", result)
	return
}

// @Summary Delete a task for a person
// @Description Delete a task for a person
// @Tags Tasks
// @Param taskId path string true "Task ID"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /people/task/{taskId} [delete]
func (h *Handler) DeleteTask(c *gin.Context) {
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

	err := h.service.DeleteTask(c.Request.Context(), id)
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
