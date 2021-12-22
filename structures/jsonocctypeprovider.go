package structures

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type JsonOccupancyTypeProvider struct {
	path           string
	occupancyTypes map[string]OccupancyTypeStochastic
}

func (jotp *JsonOccupancyTypeProvider) Init(path string) {
	jotp.path = path
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		log.Fatal("structures: unable to read json occupancy type file at path: " + path)
	}
	m := make(map[string]OccupancyTypeStochastic)
	err = json.Unmarshal(b, &m)
	if err != nil {
		log.Fatal("structures: unable to parse json occupancy type file at path: " + path)
	}
	jotp.occupancyTypes = m
}
func (jotp JsonOccupancyTypeProvider) OccupancyTypeMap() map[string]OccupancyTypeStochastic {
	return jotp.occupancyTypes
}
