package dorm

type Database struct {
	driver Driver
}

type Driver interface {
	ConnTable(string, interface{}) Table
}

func (d Database) Table(name string, strct interface{}) Table {
	return d.driver.ConnTable(name, strct)
}

func NewDatabase(driver Driver) *Database {
	return &Database{
		driver: driver,
	}
}
