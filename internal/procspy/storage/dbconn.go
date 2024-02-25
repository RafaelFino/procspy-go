package procspy_storage

import (
	"fmt"
	"log"

	"database/sql"
	procspy_config "procspy/internal/procspy/config"
)

type DbConnection struct {
	config   *procspy_config.ServerConfig
	conn     *sql.DB
	user     string
	port     int
	host     string
	dbname   string
	password string
}

func NewDbConnection(config *procspy_config.ServerConfig) *DbConnection {
	ret := &DbConnection{
		config: config,
		conn:   nil,
	}

	ret.dbname = config.DBName
	ret.user = config.User
	ret.password = config.Password
	ret.host = config.Host
	ret.port = config.Port

	return ret
}

func (d *DbConnection) connect() error {
	var err error

	d.conn, err = sql.Open("postgres", d.makeConnString())
	if err != nil {
		log.Printf("[DbConnection] Error connecting to database: %s", err)
		return err
	}
	log.Printf("[DbConnection] Connected to database")

	return err
}

func (d *DbConnection) Close() error {
	if d.conn == nil {
		log.Printf("[DbConnection] Database is already closed")
		return nil
	}

	log.Printf("[DbConnection] Disconnecting from database")
	return d.conn.Close()
}

func (d *DbConnection) makeConnString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", d.user, d.password, d.host, d.port, d.dbname)
}

func (d *DbConnection) GetConn() (*sql.DB, error) {
	if d.conn == nil {
		err := d.connect()
		if err != nil {
			return nil, err
		}
	}

	return d.conn, nil
}

func (d *DbConnection) Exec(script string, args ...interface{}) error {
	conn, err := d.GetConn()
	if err != nil {
		log.Printf("[DbConnection] Error getting connection: %s", err)
		return err
	}

	result, err := conn.Exec(script, args...)

	if err != nil {
		log.Printf("[DbConnection] Error executing script: %s \n %s", script, err)
		return err
	}

	if result != nil {
		affected, err := result.RowsAffected()
		if err != nil {
			log.Printf("[DbConnection] Error getting affected rows: %s", err)
			return err
		}

		if affected == 0 {
			log.Printf("[DbConnection] No rows affected by %s", script)
		}
	}

	return nil
}
