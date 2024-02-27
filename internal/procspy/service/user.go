package service

import "procspy/internal/procspy/domain"

type User struct {
}

func NewUser() *User {
	ret := &User{}

	return ret
}

func (u *User) GetUser(user string) (*domain.User, error) {
	return nil, nil
}
