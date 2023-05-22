package dorm

import (
	"dorm/pkg"
)

type Database struct {
	driver Driver
}

type Driver interface {
	ConnTable(string) pkg.Table
}

func (d Database) Table(name string) pkg.Table {
	return d.driver.ConnTable(name)
}

func NewDatabase(driver Driver) *Database {
	return &Database{
		driver: driver,
	}
}
