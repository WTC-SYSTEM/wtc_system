package db

import (
	"github.com/hawkkiller/wtc_system/recipe_service/pkg/client/postgresql"
	"github.com/hawkkiller/wtc_system/recipe_service/pkg/logging"
)

type db struct {
	client postgresql.Client
	logger logging.Logger
}
