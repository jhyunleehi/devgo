package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
)

func Test_T1(t *testing.T) {

	// Open our jsonFile
	jsonFile, err := os.Open("jcollection.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var jdata map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &jdata)
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Printf("[%+v]", jdata)
	err = ListRestApi(jdata)
	if err != nil {
		t.Error(err)
	}
	err = CallRestApi(jdata)
	if err != nil {
		t.Error(err)
	}
}

func ListRestApi(i interface{}) error {
	t := reflect.TypeOf(i)
	if t.String() != "map[string]interface {}" {
		log.Error(t.String())
		return errors.New("failed convert type")
	}

	item, ok := i.(map[string]interface{})
	if !ok {
		log.Errorf("[%v]", i)
		return errors.New("fail convert type")
	}

	val, ok := item["name"]
	if ok {
		log.Printf("Test:[%s]", val)
	}
	vallist, ok := item["item"]
	if ok {
		for _, v := range vallist.([]interface{}) {
			//log.Debugf("[%d][%v]", k, v)
			err := ListRestApi(v)
			if err != nil {
				log.Error(err)
				return err
			}
		}
	}
	return nil
}

func CallRestApi(i interface{}) error {
	item, ok := i.(map[string]interface{})
	if !ok {
		log.Errorf("[%v]", i)
		return errors.New("fail convert type")
	}
	t := reflect.TypeOf(i)
	if t.String() != "map[string]interface {}" {
		log.Error(t.String())
		return errors.New("failed convert type")
	}

	val, ok := item["name"]
	if ok {
		log.Printf("TEST: [%s]", val)
	}
	vallist, ok := item["item"]
	if ok {
		for _, v := range vallist.([]interface{}) {
			//log.Debugf("[%d][%v]", k, v)
			err := CallRestApi(v)
			if err != nil {
				log.Error(err)
				return err
			}
		}
	}
	r := NewResty("")
	err := r.RestApi(item)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
