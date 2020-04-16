package sync

import (
	"context"

	"github.com/redhat-developer/tekton-hub/backend/api/pkg/app"
	"go.uber.org/zap"
)

type SyncService struct {
	app       app.Config
	log       *zap.SugaredLogger
	clonePath string
}

func New(app app.Config) *SyncService {
	return &SyncService{
		app: app,
		log: app.Logger().With("service", "sync"),
	}
}

func (s *SyncService) Sync(context context.Context) {

}
