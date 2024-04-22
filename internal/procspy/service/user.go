package service

import "procspy/internal/procspy/config"

type Users struct {
	config *config.Server
}

func NewUsers(config *config.Server) *Users {
	return &Users{config: config}
}

func (u *Users) GetUsers() ([]string, error) {
	var ret []string

	for k := range u.config.UserTarges {
		ret = append(ret, k)
	}

	return ret, nil
}

func (u *Users) Exists(user string) bool {
	_, ok := u.config.UserTarges[user]
	return ok
}
