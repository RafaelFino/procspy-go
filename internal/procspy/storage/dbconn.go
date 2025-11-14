package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type DbConnection struct {
	conn *sql.DB
	path string
}

func NewDbConnection(path string) *DbConnection {
	return &DbConnection{
		path: path,
	}
}

func (d *DbConnection) makeDBPath() string {
	if d.path == ":memory:" {
		return ":memory:"
	}
	return fmt.Sprintf("%s/procspy.db", d.path)
}
func (d *DbConnection) GetConn() (*sql.DB, error) {
	path := d.makeDBPath()

	if d.conn == nil {
		log.Printf("[storage.DbConnection.GetConn] Opening database connection to '%s'", path)
		conn, err := sql.Open("sqlite", path)
		if err != nil {
			log.Printf("[storage.DbConnection.GetConn] Failed to connect to database '%s': %v", path, err)
			return nil, err
		}
		d.conn = conn
	}

	return d.conn, nil
}

func (d *DbConnection) Close() error {
	if d.conn == nil {
		log.Printf("[storage.DbConnection.Close] Database connection is already closed")
		return nil
	}

	err := d.conn.Close()

	if err != nil {
		log.Printf("[storage.DbConnection.Close] Failed to close database connection: %v", err)
		return err
	}

	d.conn = nil

	log.Printf("[storage.DbConnection.Close] Database connection closed successfully for '%s'", d.makeDBPath())

	return nil
}

func (d *DbConnection) Exec(query string, args ...any) error {
	conn, err := d.GetConn()

	if err != nil {
		log.Printf("[storage.DbConnection.Exec] Failed to get database connection: %v", err)
		return err
	}

	res, err := conn.Exec(query, args...)

	if err != nil {
		log.Printf("[storage.DbConnection.Exec] Failed to execute query: %v", err)
		return err
	}

	if res != nil {
		affected, err := res.RowsAffected()

		if err != nil {
			log.Printf("[storage.DbConnection.Exec] Failed to get rows affected count: %v", err)
			return err
		}

		lastId, err := res.LastInsertId()

		if err != nil {
			log.Printf("[storage.DbConnection.Exec] Failed to get last insert ID: %v", err)
			return err
		}

		log.Printf("[storage.DbConnection.Exec] Query executed successfully (%d rows affected, last insert ID: %d)", affected, lastId)
	}

	return nil
}
