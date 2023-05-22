package pkg

type Table interface {
	InsertOne(interface{}) (int, error)
	InsertMany(interface{}) error
	FindOne(interface{}, interface{}) error
	FindAll(interface{}) error
	//FindByOne()
	//FindByMany()
	//FindByAll()
	//UpdateOne()
	//UpdateMany()
	//DeleteOne()
	//DeleteByAll()
	//DeleteAll()
}
