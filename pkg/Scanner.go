package pkg

import (
	"reflect"
)

func ScanTagsFromKeyInStruct(obj interface{}, key string) map[string]interface{} {
	num := reflect.TypeOf(obj).NumField()
	res := make(map[string]interface{})
	for i := 0; i < num; i++ {
		val := reflect.TypeOf(obj).Field(i).Tag.Get(key)
		if reflect.ValueOf(obj).Field(i).Kind() == reflect.Int {
			res[val] = reflect.ValueOf(obj).Field(i).Int()
		} else {
			res[val] = reflect.ValueOf(obj).Field(i).String()
		}
	}
	return res
}

//func ScanTypeFromKeyInStruct(obj interface{}, key string) map[string]string {
//	num := reflect.TypeOf(obj).NumField()
//	res := make(map[string]string)
//	for i := 0; i < num; i++ {
//		dataType := reflect.TypeOf(obj).Field(i).Type.String()
//		name := reflect.TypeOf(obj).Field(i).Tag.Get(key)
//		if name == "id" {
//			if reflect.TypeOf(obj).Field(i).Tag.Get("serial") == "true" {
//				dataType = "serial"
//			}
//			if reflect.TypeOf(obj).Field(i).Tag.Get("uq") == "true" {
//				dataType += "_uq"
//			}
//		}
//		res[name] = dataType
//	}
//	return res
//}

// >> rel >> fk >> field - foreign key
// >> rel >> pk >> field - primary key

func ScanTypeFromKeyInStruct(obj interface{}) map[string]map[string]string {
	num := reflect.TypeOf(obj).NumField()
	res := make(map[string]map[string]string)
	keys := [6]string{"fk", "pk", "uq", "rel", "field", "serial"}
	for i := 0; i < num; i++ {
		keyValues := make(map[string]string)
		dataType := reflect.TypeOf(obj).Field(i).Type.String()
		name := reflect.TypeOf(obj).Field(i).Tag.Get("db")

		if name == "" {
			continue
		}

		if name == "id" && reflect.TypeOf(obj).Field(i).Tag.Get("serial") == "true" {
			keyValues["dataType"] = "serial"
			res[name] = keyValues
			continue
		}

		for _, key := range keys {
			keyValue := reflect.TypeOf(obj).Field(i).Tag.Get(key)
			keyValues[key] = keyValue
		}
		keyValues["dataType"] = dataType
		res[name] = keyValues
	}
	return res
}
