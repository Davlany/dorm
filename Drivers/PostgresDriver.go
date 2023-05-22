package Drivers

import (
	"dorm/pkg"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type PostgresDriver struct {
	conn *sqlx.DB
}

type PgTable struct {
	name string
	pd   *PostgresDriver
}

func (pd PostgresDriver) ConnTable(name string) pkg.Table {
	return PgTable{
		name: name,
		pd:   &pd,
	}
}

func NewPostgresDriver(user, password, dbName, port, sslMode string) (*PostgresDriver, error) {
	conn, err := sqlx.Connect("postgres", fmt.Sprintf("user = %s password = %s dbname = %s sslmode = %s port = %s", user, password, dbName, sslMode, port))
	if err != nil {
		return nil, err
	}
	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	return &PostgresDriver{
		conn: conn,
	}, nil
}
