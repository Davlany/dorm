package pkg

type Table interface {
	InsertOne(interface{}) (int, error)
	InsertMany(interface{}) error
	//FindOne()
	//FindByOne()
	//FindByMany()
	//FindByAll()
	//FindAll()
	//UpdateOne()
	//UpdateMany()
	//DeleteOne()
	//DeleteByAll()
	//DeleteAll()
}
