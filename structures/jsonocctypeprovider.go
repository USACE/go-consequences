package structures

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type JsonOccupancyTypeProvider struct {
	path                    string
	occupancyTypesContainer OccupancyTypesContainer
}

func (jotp *JsonOccupancyTypeProvider) Init(path string) {
	jotp.path = path
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		log.Fatal("structures: unable to read json occupancy type file at path: " + path)
	}
	//m := make(map[string]OccupancyTypeStochastic)
	c := OccupancyTypesContainer{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		log.Fatal("structures: unable to parse json occupancy type file at path: " + path)
	}
	jotp.occupancyTypesContainer = c
}
func (jotp JsonOccupancyTypeProvider) OccupancyTypeMap() map[string]OccupancyTypeStochastic {
	return jotp.occupancyTypesContainer.OccupancyTypes
}
func (jotp JsonOccupancyTypeProvider) Write(path string) error {
	w, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer w.Close()
	b, err := json.Marshal(jotp.occupancyTypesContainer)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}
