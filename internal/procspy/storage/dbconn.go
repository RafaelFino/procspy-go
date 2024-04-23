package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

type DbConnection struct {
	conn *sql.DB
	path string
	last string
}

func NewDbConnection(path string) *DbConnection {
	return &DbConnection{
		path: path,
	}
}

func (d *DbConnection) makeDBPath() string {
	return fmt.Sprintf("%s/%s.db", d.path, time.Now().Format("2006-01-02"))
}
func (d *DbConnection) GetConn() (*sql.DB, error) {
	path := d.makeDBPath()

	if d.last != path {
		err := d.Close()

		if err != nil {
			log.Printf("[DbConnection] Error closing connection: %s", err)
			return nil, err
		}
	}

	if d.conn == nil {
		log.Printf("[DbConnection] Opening connection to %s", path)
		conn, err := sql.Open("sqlite", path)
		if err != nil {
			log.Printf("[DbConnection] Error connecting to database: %s", err)
			return nil, err
		}
		d.conn = conn
		d.last = path
	}

	return d.conn, nil
}

func (d *DbConnection) Close() error {
	if d.conn == nil {
		log.Printf("[DbConnection] Database is already closed")
		return nil
	}

	err := d.conn.Close()

	if err != nil {
		log.Printf("[DbConnection] Error closing connection: %s", err)
		return err
	}

	log.Printf("[DbConnection] Connection closed for %s", d.last)

	return nil
}

func (d *DbConnection) Exec(query string, args ...interface{}) error {
	conn, err := d.GetConn()

	if err != nil {
		log.Printf("[DbConnection] Error getting connection: %s", err)
		return err
	}

	res, err := conn.Exec(query, args...)

	if err != nil {
		log.Printf("[DbConnection] Error executing query: %s", err)
		return err
	}

	if res != nil {
		affected, err := res.RowsAffected()

		if err != nil {
			log.Printf("[DbConnection] Error getting rows affected: %s", err)
			return err
		}

		lastId, err := res.LastInsertId()

		if err != nil {
			log.Printf("[DbConnection] Error getting last id: %s", err)
			return err
		}

		log.Printf("[DbConnection] Query executed successfully: %d rows affected -> lastId: %d", affected, lastId)
	}

	return nil
}
