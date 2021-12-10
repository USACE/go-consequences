package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

//
// Read json file into a flexible interface
// i should be a pointer to a struct - https://mholt.github.io/json-to-go/
func ReadJson(path string, i interface{}) error {

	jsonFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Unable to read json file from:" + path)
	}

	err2 := json.Unmarshal([]byte(jsonFile), i)
	if err2 != nil {
		log.Fatal("Unable to marshal json to interface")
	}

	return err
}

// func removeOneLevel(f *flex) (flex, error) {

// 	for idx, v := range f {
// 		return v
// 	}
// }

// // Remove flex interface overhead, output first level should be an array
// func TraverseToArray(f *flex) (flex, error) {

// 	while

// 	return

// }
