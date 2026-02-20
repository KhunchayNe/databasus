package restore_targets

import (
	"databasus-backend/internal/storage"
	workspaces_services "databasus-backend/internal/features/workspaces/services"
	"databasus-backend/internal/util/logger"
)

var (
	restoreTargetRepository *RestoreTargetRepository
	restoreTargetService    *RestoreTargetService
	restoreTargetController *RestoreTargetController
)

func GetRestoreTargetController() *RestoreTargetController {
	return restoreTargetController
}

func SetupDependencies() {
	db := storage.GetDb()
	workspaceService := workspaces_services.GetWorkspaceService()

	restoreTargetRepository = NewRestoreTargetRepository(db)
	restoreTargetService = NewRestoreTargetService(restoreTargetRepository, workspaceService)
	restoreTargetController = NewRestoreTargetController(restoreTargetService)

	logger.GetLogger().Info("Restore targets dependencies initialized")
}
