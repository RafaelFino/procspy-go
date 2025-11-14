package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"procspy/internal/procspy/config"
	"procspy/internal/procspy/domain"
)

type Target struct {
	urls map[string]string
}

func NewTarget(config *config.Server) *Target {
	return &Target{
		urls: config.UserTarges,
	}
}

func (t *Target) GetTargets(user string) (*domain.TargetList, error) {
	ret := &domain.TargetList{
		Targets: []*domain.Target{},
	}

	for k, v := range t.urls {
		if k == user {
			data, err := t.getFromUrl(v)

			if err != nil {
				log.Printf("[service.Target.GetTargets] Failed to fetch targets from URL '%s' for user '%s': %v", v, user, err)
				return nil, err
			}

			ret, err = domain.TargetListFromJson(data)

			if err != nil {
				log.Printf("[service.Target.GetTargets] Failed to parse target list JSON for user '%s': %v", user, err)
				return nil, err
			}
			break
		}
	}

	for k, v := range ret.Targets {
		v.User = user
		ret.Targets[k] = v
	}

	if ret == nil {
		log.Printf("[service.Target.GetTargets] No targets configured for user '%s'", user)
		return nil, fmt.Errorf("no targets found for user: %s", user)
	}

	return ret, nil
}

func (t *Target) getFromUrl(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Printf("[service.Target.getFromUrl] Failed to fetch URL '%s': %v", url, err)
		return "", err
	}

	if res.StatusCode != 200 {
		log.Printf("[service.Target.getFromUrl] Received non-OK status code %d from URL '%s'", res.StatusCode, url)
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("[service.Target.getFromUrl] Failed to read response body from URL '%s': %v", url, err)
		return "", err
	}

	return string(body), nil
}
