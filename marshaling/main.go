package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Node struct {
	Id    string `json:"-"`
	Name  string `json:"name,string"`
	Group string `json:"group,string"`
	Size  int64  `json:"size"`
}

type Node2 struct {
	Id    string
	Name  string
	Group string
	Size  int
}

func main() {
	node := Node{
		Id:    "ID-123456",
		Name:  "jhyunlee",
		Group: "myGroup",
		Size:  1234,
	}
	nodestr1, _ := json.Marshal(node)
	fmt.Printf("[%s]\n", nodestr1)

	nodestr2, _ := MarshalJsonIgnoreTags(node)
	fmt.Printf("[%s]\n", nodestr2)

	n1 := Node{}
	json.Unmarshal(nodestr1, &n1)
	fmt.Printf("[%+v]\n", n1)

	n2 := Node2{}
	json.Unmarshal(nodestr2, &n2)

	fmt.Printf("[%+v]\n", n2)

	n3 := Node{}
	UnmarshalJasonIgnoreTags(nodestr2, &n3)

	fmt.Printf("===>>>[%+v]\n", n3)

}

func MarshalJsonIgnoreTags(obj interface{}) ([]byte, error) {
	m := make(map[string]interface{})
	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	for i := 0; i < val.NumField(); i++ {
		fieldName := typ.Field(i).Name
		m[fieldName] = val.Field(i).Interface()
	}
	return json.Marshal(m)
}

func UnmarshalJasonIgnoreTags(data []byte, v interface{}) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	val := reflect.ValueOf(v).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		// Convert struct field name to a format that matches the JSON key
		//jsonKey := strings.ToLower(typ.Field(i).Name)
		jsonKey := typ.Field(i).Name

		if value, exists := m[jsonKey]; exists {
			field := val.Field(i)
			// Check if the destination field is an integer type and
			// the source value is a float64
			if field.Kind() >= reflect.Int && field.Kind() <= reflect.Int64 {
				if floatValue, ok := value.(float64); ok && floatValue == float64(int64(floatValue)) {
					field.SetInt(int64(floatValue))
					continue
				}
			}

			if field.IsValid() && field.CanSet() {
				field.Set(reflect.ValueOf(value))
			}
		}
	}

	return nil
}
