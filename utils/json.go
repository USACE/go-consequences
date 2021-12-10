package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

// Read json file into a flexible interface
// i should be a pointer to a struct - https://mholt.github.io/json-to-go/
func ReadJson(path string, i interface{}) error {

	jsonFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Unable to read json file from:" + path)
	}

	err2 := json.Unmarshal([]byte(jsonFile), i)

	// Wrong type definitions during unmarshaling probably isn't a big deal
	if err2 != nil && !strings.Contains(
		err2.Error(),
		"json: cannot unmarshal string into Go struct field",
	) {
		log.Fatal("Unable to marshal json to interface")
	}

	return err
}
