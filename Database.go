package dorm

import (
	"dorm/pkg"
)

type Database struct {
	driver Driver
}

type Driver interface {
	ConnTable(string, interface{}) pkg.Table
}

func (d Database) Table(name string, strct interface{}) pkg.Table {
	return d.driver.ConnTable(name, strct)
}

func NewDatabase(driver Driver) *Database {
	return &Database{
		driver: driver,
	}
}
