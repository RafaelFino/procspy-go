package storage

import (
	"errors"
	"log"
	domain "procspy/internal/procspy/domain"
)

type User struct {
	conn *DbConnection
}

func NewUser(dbConn *DbConnection) *User {
	ret := &User{
		conn: dbConn,
	}

	err := ret.Init()

	if err != nil {
		log.Printf("[Storage.User] Error initializing storage: %s", err)
	}

	return ret
}

func (u *User) Init() error {
	create := `
CREATE TABLE IF NOT EXISTS users (
	name varchar(128) PRIMARY KEY NOT NULL,
	key TEXT DEFAULT NULL,
	approved BOOLEAN DEFAULT FALSE,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP()
);
`
	if u.conn == nil {
		log.Printf("[Storage.User] Error creating tables: db is nil")
		return errors.New("db is nil")
	}

	err := u.conn.Exec(create)

	if err != nil {
		log.Printf("[Storage.User] Error creating tables: %s", err)
	}

	return err
}

func (u *User) Close() error {
	if u.conn == nil {
		log.Printf("[Storage.User] Database is already closed")
		return nil
	}

	return u.conn.Close()
}

func (u *User) CreateUser(name string, key string) error {
	insert := `
INSERT INTO users (name, key) VALUES (?, ?)
`
	if u.conn == nil {
		log.Printf("[Storage.User] Error creating user: db is nil")
		return errors.New("db is nil")
	}

	err := u.conn.Exec(insert, name, key)

	if err != nil {
		log.Printf("[Storage.User] Error creating user: %s", err)
	}

	return err
}

func (u *User) GetUser(name string) (*domain.User, error) {
	query := `
SELECT
	key,
	approved,
	created_at
FROM
	users
WHERE
	name = ?
ORDER BY created_at DESC
LIMIT 1;
`
	if u.conn == nil {
		log.Printf("[Storage.User] Error getting user: db is nil")
		return nil, errors.New("db is nil")
	}

	conn, err := u.conn.GetConn()
	if err != nil {
		log.Printf("[Storage.User] Error getting connection: %s", err)
		return nil, err
	}

	rows, err := conn.Query(query, name)
	if err != nil {
		log.Printf("[Storage.User] Error getting user: %s", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var key string
		var approved bool
		var createdAt string

		err = rows.Scan(&key, &approved, &createdAt)

		user := *domain.NewUser(name)
		user.SetKey(key)
		user.SetApproved(approved)
		user.SetCreatedAt(createdAt)

		if err != nil {
			log.Printf("[Storage.User] Error scanning user: %s", err)
			return nil, err
		}

		return &user, err
	}

	return nil, err
}
