package dorm

type Database struct {
	driver Driver
}

type Driver interface {
	ConnTable(string, interface{}) Table
	RegisterSchemas(interface{}) error
}

func (d Database) Table(name string, strct interface{}) Table {
	return d.driver.ConnTable(name, strct)
}
func (d Database) RegisterSchemas(schemas interface{}) error {
	return d.driver.RegisterSchemas(schemas)
}

func NewDatabase(driver Driver) *Database {
	return &Database{
		driver: driver,
	}
}
