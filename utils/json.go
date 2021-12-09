package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type flex interface{}

func ReadJson(path string) (flex, error) {

	jsonFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Unable to read json file from:" + path)
	}

	var data flex
	err2 := json.Unmarshal([]byte(jsonFile), &data)
	if err2 != nil {
		log.Fatal("Unable to marshal json to interface")
	}

	return data, err
}
