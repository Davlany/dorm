package dorm

type Table interface {
	InsertOne(interface{}) (int, error)
	InsertMany(interface{}) error
	FindOne(interface{}, interface{}) error
	FindAll(interface{}) error
	UpdateOne(interface{}) error
	UpdateMany(interface{}) error
	DeleteOne(interface{}) error
	//FindByOne()
	//FindByMany()
	//FindByAll()
	//DeleteByAll()
	//DeleteAll()
}
