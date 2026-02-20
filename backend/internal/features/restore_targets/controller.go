package restore_targets

import (
	users_middleware "databasus-backend/internal/features/users/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RestoreTargetController struct {
	restoreTargetService *RestoreTargetService
}

func NewRestoreTargetController(
	restoreTargetService *RestoreTargetService,
) *RestoreTargetController {
	return &RestoreTargetController{
		restoreTargetService: restoreTargetService,
	}
}

func (c *RestoreTargetController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/restore-targets", c.GetRestoreTargets)
	router.POST("/restore-targets", c.CreateRestoreTarget)
	router.GET("/restore-targets/:id", c.GetRestoreTarget)
	router.PUT("/restore-targets/:id", c.UpdateRestoreTarget)
	router.DELETE("/restore-targets/:id", c.DeleteRestoreTarget)
}

// CreateRestoreTarget
// @Summary Create a restore target
// @Description Create a new restore target for quick restore
// @Tags restore-targets
// @Accept json
// @Produce json
// @Param request body CreateRestoreTargetWithDBRequest true "Restore target data"
// @Success 200 {object} RestoreTarget
// @Failure 400
// @Failure 401
// @Router /api/v1/restore-targets [post]
func (c *RestoreTargetController) CreateRestoreTarget(ctx *gin.Context) {
	user, ok := users_middleware.GetUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CreateRestoreTargetWithDBRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	target, err := c.restoreTargetService.CreateRestoreTarget(ctx, user, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, target)
}

// GetRestoreTargets
// @Summary Get all restore targets
// @Description Get all restore targets for the current workspace
// @Tags restore-targets
// @Produce json
// @Success 200 {array} RestoreTarget
// @Failure 401
// @Router /api/v1/restore-targets [get]
func (c *RestoreTargetController) GetRestoreTargets(ctx *gin.Context) {
	user, ok := users_middleware.GetUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	targets, err := c.restoreTargetService.GetRestoreTargets(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, targets)
}

// GetRestoreTarget
// @Summary Get a restore target
// @Description Get a specific restore target by ID
// @Tags restore-targets
// @Produce json
// @Param id path string true "Restore Target ID"
// @Success 200 {object} RestoreTarget
// @Failure 400
// @Failure 401
// @Failure 404
// @Router /api/v1/restore-targets/{id} [get]
func (c *RestoreTargetController) GetRestoreTarget(ctx *gin.Context) {
	user, ok := users_middleware.GetUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid restore target ID"})
		return
	}

	target, err := c.restoreTargetService.GetRestoreTarget(ctx, user, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "restore target not found"})
		return
	}

	ctx.JSON(http.StatusOK, target)
}

// UpdateRestoreTarget
// @Summary Update a restore target
// @Description Update name and description of a restore target
// @Tags restore-targets
// @Accept json
// @Produce json
// @Param id path string true "Restore Target ID"
// @Param request body UpdateRestoreTargetRequest true "Update data"
// @Success 200 {object} RestoreTarget
// @Failure 400
// @Failure 401
// @Failure 404
// @Router /api/v1/restore-targets/{id} [put]
func (c *RestoreTargetController) UpdateRestoreTarget(ctx *gin.Context) {
	user, ok := users_middleware.GetUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid restore target ID"})
		return
	}

	var req UpdateRestoreTargetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	target, err := c.restoreTargetService.UpdateRestoreTarget(ctx, user, id, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, target)
}

// DeleteRestoreTarget
// @Summary Delete a restore target
// @Description Delete a restore target
// @Tags restore-targets
// @Param id path string true "Restore Target ID"
// @Success 200 {object} map[string]string
// @Failure 400
// @Failure 401
// @Failure 404
// @Router /api/v1/restore-targets/{id} [delete]
func (c *RestoreTargetController) DeleteRestoreTarget(ctx *gin.Context) {
	user, ok := users_middleware.GetUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid restore target ID"})
		return
	}

	err = c.restoreTargetService.DeleteRestoreTarget(ctx, user, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "restore target not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "restore target deleted successfully"})
}
