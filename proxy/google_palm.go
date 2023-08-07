package proxy

import (
	"github.com/mohamedveron/go_app_template/internal/pkg/logger"
	"github.com/mr-destructive/palm"
)

func PrintModels() {
	models, err := palm.ListModels()
	if err != nil {
		logger.Error(err)
	}
	logger.Info("models", models)
}
