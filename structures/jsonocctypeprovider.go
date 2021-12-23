package structures

import (
	"encoding/json"
	"errors"
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
func (jotp *JsonOccupancyTypeProvider) ExtendMap(extension map[string]OccupancyTypeStochastic) error {
	for key, value := range extension {
		_, exists := jotp.occupancyTypes[key]
		if exists {
			return errors.New("structures: occupancy type " + key + " already exists")
		} else {
			jotp.occupancyTypes[key] = value
		}
	}
	return nil
}
func (jotp *JsonOccupancyTypeProvider) MergeMap(additionalDFs map[string]OccupancyTypeStochastic) error {
	for key, value := range additionalDFs {
		curval, exists := jotp.occupancyTypes[key]
		if exists {
			for parameterkey, sdf := range value.StructureDFF.DamageFunctions {
				_, sdfexists := curval.StructureDFF.DamageFunctions[parameterkey]
				if sdfexists {
					return errors.New("structures: occupancy type " + key + " already exists with parameter " + parameterkey.String())
				} else {
					curval.StructureDFF.DamageFunctions[parameterkey] = sdf
				}
			}
			for parameterkey, cdf := range value.ContentDFF.DamageFunctions {
				_, cdfexists := curval.ContentDFF.DamageFunctions[parameterkey]
				if cdfexists {
					return errors.New("structures: occupancy type " + key + " already exists with parameter " + parameterkey.String())
				} else {
					curval.ContentDFF.DamageFunctions[parameterkey] = cdf
				}
			}
			jotp.occupancyTypes[key] = curval
		} else {
			jotp.occupancyTypes[key] = value
		}
	}
	return nil
}
func (jotp *JsonOccupancyTypeProvider) OverrideMap(overrides map[string]OccupancyTypeStochastic) error {
	for key, value := range overrides {
		curval, exists := jotp.occupancyTypes[key]
		if exists {
			for parameterkey, sdf := range value.StructureDFF.DamageFunctions {
				_, sdfexists := curval.StructureDFF.DamageFunctions[parameterkey]
				if sdfexists {
					curval.StructureDFF.DamageFunctions[parameterkey] = sdf
				} else {
					return errors.New("structures: occupancy type " + key + " doesn't currently exist with parameter " + parameterkey.String())
				}
			}
			for parameterkey, cdf := range value.ContentDFF.DamageFunctions {
				_, cdfexists := curval.ContentDFF.DamageFunctions[parameterkey]
				if cdfexists {
					curval.ContentDFF.DamageFunctions[parameterkey] = cdf
				} else {
					return errors.New("structures: occupancy type " + key + " doesn't currently exist with parameter " + parameterkey.String())
				}
			}
			jotp.occupancyTypes[key] = curval
		} else {
			return errors.New("structures: occupancy type " + key + " doesn't currently exist.")
		}
	}
	return nil
}
