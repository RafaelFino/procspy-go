package service

import "procspy/internal/procspy/domain"

type Target struct {
}

func NewTarget() *Target {
	ret := &Target{}

	return ret
}

func (t *Target) InsertTarget(target *domain.Target) error {
	return nil
}

func (t *Target) DeleteTargets(user string) error {
	return nil
}

func (t *Target) GetTargets(user string) (map[string]*domain.Target, error) {
	return nil, nil
}
