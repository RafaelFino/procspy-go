package domain

import (
	"encoding/json"
	"log"
)

type User struct {
	Name      string `json:"name"`
	Key       string `json:"key"`
	Approved  bool   `json:"approved"`
	CreatedAt string `json:"created_at"`
	Token     string `json:"token,omitempty"`
}

func NewUser(name string) *User {
	return &User{
		Name: name,
	}
}

func (u *User) SetKey(key string) {
	u.Key = key
}

func (u *User) SetApproved(approved bool) {
	u.Approved = approved
}

func (u *User) SetCreatedAt(created_at string) {
	u.CreatedAt = created_at
}

func (u *User) SetToken(token string) {
	u.Token = token
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetKey() string {
	return u.Key
}

func (u *User) GetApproved() bool {
	return u.Approved
}

func (u *User) GetCreatedAt() string {
	return u.CreatedAt
}

func (u *User) GetToken() string {
	return u.Token
}

func (u *User) ToJson() string {
	ret, err := json.MarshalIndent(u, "", "\t")
	if err != nil {
		log.Printf("[User] Error parsing json: %s", err)
	}

	return string(ret)
}

func UserFromJson(jsonString string) (*User, error) {
	ret := NewUser("")
	err := json.Unmarshal([]byte(jsonString), &ret)
	if err != nil {
		log.Printf("[User] Error unmarshalling user: %s", err)
		return nil, err
	}

	return ret, nil
}
