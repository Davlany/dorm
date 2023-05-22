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
