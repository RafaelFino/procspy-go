package service

import (
	"log"
	"procspy/internal/procspy/domain"
	"procspy/internal/procspy/storage"
)

type Target struct {
	storage *storage.Target
	dbConn  *storage.DbConnection
}

func NewTarget(dbConn *storage.DbConnection) *Target {
	ret := &Target{
		dbConn:  dbConn,
		storage: storage.NewTarget(dbConn),
	}

	return ret
}

func (t *Target) InsertTarget(target *domain.Target) error {
	log.Printf("[service.Target] Inserting target: %s", target.Name)
	return t.storage.InsertTarget(target)
}

func (t *Target) DeleteTargets(user string) error {
	log.Printf("[service.Target] Deleting targets: %s", user)
	return t.storage.DeleteTargets(user)
}

func (t *Target) GetTargets(user string) (map[string]*domain.Target, error) {
	log.Printf("[service.Target] Getting targets: %s", user)
	return t.storage.GetTargets(user)
}
