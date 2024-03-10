package service

import (
	"log"
	"procspy/internal/procspy/domain"
	"procspy/internal/procspy/storage"
)

type User struct {
	storage *storage.User
	dbConn  *storage.DbConnection
}

func NewUser(dbConn *storage.DbConnection) *User {
	ret := &User{
		dbConn:  dbConn,
		storage: storage.NewUser(dbConn),
	}

	return ret
}

func (u *User) CreateUser(name string, key string) error {
	log.Printf("[service.User] Creating user: %s", name)
	return u.storage.CreateUser(name, key)
}

func (u *User) ApproveUser(name string) error {
	log.Printf("[service.User] Approving user: %s", name)
	return u.storage.ApproveUser(name)
}

func (u *User) GetUser(user string) (*domain.User, error) {
	log.Printf("[service.User] Getting user: %s", user)
	return u.storage.GetUser(user)
}
