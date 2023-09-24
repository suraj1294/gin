package utils

import (
	"encoding/json"
	"reflect"
	"strings"
)

func MarshalRequest(data interface{}) (*map[string]interface{}, *error) {

	var inInterface map[string]interface{}
	inrec, err := json.Marshal(data)

	if err != nil {
		return nil, &err
	}

	json.Unmarshal(inrec, &inInterface)

	return &inInterface, nil

}

/*
This function will help you to convert your object from struct to map[string]interface{} based on your JSON tag in your structs.
Example how to use posted in sample_test.go file.
*/
func StructToMap(item interface{}) map[string]interface{} {

	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")

		// remove omitEmpty
		omitEmpty := false
		if strings.HasSuffix(tag, "omitempty") {
			omitEmpty = true
			idx := strings.Index(tag, ",")
			if idx > 0 {
				tag = tag[:idx]
			} else {
				tag = ""
			}
		}

		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = StructToMap(field)
			} else {
				if !(omitEmpty && reflectValue.Field(i).IsZero()) {
					res[tag] = field
				}
			}
		}
	}
	return res
}
