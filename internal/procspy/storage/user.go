package procspy_storage

import (
	"errors"
	"log"
	procspy_domains "procspy/internal/procspy/domain"
)

type UserStorage struct {
	conn *DbConnection
}

func NewUserStorage(dbConn *DbConnection) *UserStorage {
	ret := &UserStorage{
		conn: dbConn,
	}

	err := ret.Init()

	if err != nil {
		log.Printf("[UserStorage] Error initializing storage: %s", err)
	}

	return ret
}

func (u *UserStorage) Init() error {
	create := `
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	name varchar(128) NOT NULL,
	key TEXT DEFAULT NULL,
	approved BOOLEAN DEFAULT FALSE,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP()
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP()
);
`
	if u.conn == nil {
		log.Printf("[UserStorage] Error creating tables: db is nil")
		return errors.New("db is nil")
	}

	err := u.conn.Exec(create)

	if err != nil {
		log.Printf("[UserStorage] Error creating tables: %s", err)
	}

	return err
}

func (u *UserStorage) Close() error {
	if u.conn == nil {
		log.Printf("[UserStorage] Database is already closed")
		return nil
	}

	return u.conn.Close()
}

func (u *UserStorage) CreateUser(name string, key string) error {
	insert := `
INSERT INTO users (name, key) VALUES (?, ?)
`
	if u.conn == nil {
		log.Printf("[UserStorage] Error creating user: db is nil")
		return errors.New("db is nil")
	}

	err := u.conn.Exec(insert, name, key)

	if err != nil {
		log.Printf("[UserStorage] Error creating user: %s", err)
	}

	return err
}

func (u *UserStorage) GetUserById(id int) (*procspy_domains.User, error) {
	query := `
SELECT
	id,
	name,
	key,
	approved,
	created_at
FROM
	users
WHERE
	id = ?
ORDER BY created_at DESC
LIMIT 1;
`
	return u.loadUser(query, id)
}

func (u *UserStorage) GetUser(name string) (*procspy_domains.User, error) {
	query := `
SELECT
	id,
	name,
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
	return u.loadUser(query, name)
}

func (u *UserStorage) loadUser(query string, args ...interface{}) (*procspy_domains.User, error) {
	if u.conn == nil {
		log.Printf("[UserStorage] Error getting user: db is nil")
		return nil, errors.New("db is nil")
	}

	conn, err := u.conn.GetConn()
	if err != nil {
		log.Printf("[UserStorage] Error getting connection: %s", err)
		return nil, err
	}

	rows, err := conn.Query(query, args...)
	if err != nil {
		log.Printf("[UserStorage] Error getting user: %s", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var key string
		var approved bool
		var createdAt string

		err = rows.Scan(&id, &name, &key, &approved, &createdAt)

		user := *procspy_domains.NewUser(name)
		user.SetId(id)
		user.SetKey(key)
		user.SetApproved(approved)
		user.SetCreatedAt(createdAt)

		if err != nil {
			log.Printf("[UserStorage] Error scanning user: %s", err)
			return nil, err
		}

		return &user, err
	}

	return nil, err
}
