package structures

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/HydrologicEngineeringCenter/go-statistics/paireddata"
	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
	"github.com/USACE/go-consequences/hazards"
)

type OccupancyTypesContainer struct {
	OccupancyTypes map[string]OccupancyTypeStochastic `json:"occupancytypes"`
}

//DamageFunctionFamily is to support a family of damage functions stored by hazard parameter types
type DamageFunctionFamily struct {
	DamageFunctions map[hazards.Parameter]DamageFunction `json:"damagefunctions"` //parameter is a bitflag
}

func (dff DamageFunctionFamily) MarshalJSON() ([]byte, error) {
	s := "{\"damagefunctions\":{"
	for key, val := range dff.DamageFunctions {
		pstring, err := json.Marshal(key)
		if err != nil {
			return nil, errors.New("structures: could not marshal damage function family parameter key")
		}
		s += fmt.Sprintf("%v:", string(pstring))
		vstring, err := json.Marshal(val)
		if err != nil {
			return nil, errors.New("structures: could not marshal damage function family parameter value")
		}
		s += fmt.Sprintf("%v,", string(vstring))
	}
	s = strings.TrimRight(s, ",")
	s += "}"
	return []byte(s), nil
}
func (dff *DamageFunctionFamily) UnmarshalJSON(b []byte) error {
	m := map[string]interface{}{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	valueBytes, err := json.Marshal(m["damagefunctions"])
	if err != nil {
		return err
	}
	var functions map[string]DamageFunction
	if err = json.Unmarshal(valueBytes, &functions); err != nil {
		return err
	}
	damgfunctions := make(map[hazards.Parameter]DamageFunction)
	for key, value := range functions {
		var p hazards.Parameter
		key = "\"" + key + "\""
		b := []byte(key)
		err := json.Unmarshal(b, &p)
		if err != nil {
			return errors.New("structures: could not unmarshal parameter key " + key)
		}
		damgfunctions[p] = value
	}
	dff.DamageFunctions = damgfunctions
	return nil
}

//DamageFunctionFamilyStochastic is to support a family of damage functions stored by hazard parameter types that can represent uncertain paired data
type DamageFunctionFamilyStochastic struct {
	DamageFunctions map[hazards.Parameter]DamageFunctionStochastic `json:"damagefunctions"` //parameter is a bitflag
}

func (dffs DamageFunctionFamilyStochastic) MarshalJSON() ([]byte, error) {
	s := "{\"damagefunctions\":{"
	for key, val := range dffs.DamageFunctions {
		pstring, err := json.Marshal(key)
		if err != nil {
			return nil, errors.New("structures: could not marshal damage function family stochastic parameter key")
		}
		s += fmt.Sprintf("%v:", string(pstring))
		vstring, err := json.Marshal(val)
		if err != nil {
			return nil, errors.New("structures: could not marshal damage function family stochastic parameter value")
		}
		s += fmt.Sprintf("%v,", string(vstring))
	}
	s = strings.TrimRight(s, ",")
	s += "}}"
	return []byte(s), nil
}
func (dffs *DamageFunctionFamilyStochastic) UnmarshalJSON(b []byte) error {
	m := map[string]interface{}{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	valueBytes, err := json.Marshal(m["damagefunctions"])
	if err != nil {
		return err
	}
	var functions map[string]DamageFunctionStochastic
	if err = json.Unmarshal(valueBytes, &functions); err != nil {
		return err
	}
	damgfunctions := make(map[hazards.Parameter]DamageFunctionStochastic)
	for key, value := range functions {
		var p hazards.Parameter
		key = "\"" + key + "\""
		b := []byte(key)
		err := json.Unmarshal(b, &p)
		if err != nil {
			return errors.New("structures: could not unmarshal parameter key " + key)
		}
		damgfunctions[p] = value
	}
	dffs.DamageFunctions = damgfunctions
	return nil
}

type DamageFunction struct {
	Source         string                `json:"source"`
	DamageDriver   hazards.Parameter     `json:"damagedriver"`
	DamageFunction paireddata.PairedData `json:"damagefunction"`
}
type DamageFunctionStochastic struct {
	Source         string                           `json:"source"`
	DamageDriver   hazards.Parameter                `json:"damagedriver"`
	DamageFunction paireddata.UncertaintyPairedData `json:"damagefunction"`
}

//OccupancyTypeStochastic is used to describe an occupancy type with uncertainty in the damage relationships it produces an OccupancyTypeDeterministic through the UncertaintyOccupancyTypeSampler interface
type OccupancyTypeStochastic struct { //this is mutable
	Name                     string                                    `json:"name"`
	ComponentDamageFunctions map[string]DamageFunctionFamilyStochastic `json:"componentdamagefunctions"`
}

//OccupancyTypeDeterministic is used to describe an occupancy type without uncertainty in the damage relationships
type OccupancyTypeDeterministic struct {
	Name                     string                          `json:"name"`
	ComponentDamageFunctions map[string]DamageFunctionFamily `json:"componentdamagefunctions"`
}

func (otc *OccupancyTypesContainer) ExtendMap(extension map[string]OccupancyTypeStochastic) error {
	for key, value := range extension {
		_, exists := otc.OccupancyTypes[key]
		if exists {
			return errors.New("structures: occupancy type " + key + " already exists")
		} else {
			otc.OccupancyTypes[key] = value
		}
	}
	return nil
}
func (otc *OccupancyTypesContainer) MergeMap(additionalDFs map[string]OccupancyTypeStochastic) error {
	for key, value := range additionalDFs {
		curval, exists := otc.OccupancyTypes[key] //occupancy type
		if exists {
			for componentkey, cdff := range value.ComponentDamageFunctions {
				curcdff, componentExists := curval.ComponentDamageFunctions[componentkey]
				if componentExists {
					for parameterkey, sdf := range cdff.DamageFunctions {
						_, sdfexists := curcdff.DamageFunctions[parameterkey]
						if sdfexists {
							return errors.New("structures: occupancy type " + key + " already exists with parameter " + parameterkey.String() + " on component " + componentkey)
						} else {
							curcdff.DamageFunctions[parameterkey] = sdf
						}
					}
					curval.ComponentDamageFunctions[componentkey] = curcdff
				} else {
					curval.ComponentDamageFunctions[componentkey] = cdff
				}
			}
			otc.OccupancyTypes[key] = curval
		} else {
			otc.OccupancyTypes[key] = value
		}
	}
	return nil
}
func (otc *OccupancyTypesContainer) OverrideMap(overrides map[string]OccupancyTypeStochastic) error {
	for key, value := range overrides {
		curval, exists := otc.OccupancyTypes[key]
		if exists {
			for componentkey, cdff := range value.ComponentDamageFunctions {
				curcdff, componentExists := curval.ComponentDamageFunctions[componentkey]
				if componentExists {
					for parameterkey, sdf := range cdff.DamageFunctions {
						_, sdfexists := curcdff.DamageFunctions[parameterkey]
						if sdfexists {
							curcdff.DamageFunctions[parameterkey] = sdf
						} else {
							return errors.New("structures: occupancy type " + key + " doesn't currently exist with parameter " + parameterkey.String() + " on component " + componentkey)
						}
					}
					curval.ComponentDamageFunctions[componentkey] = curcdff
				} else {
					curval.ComponentDamageFunctions[componentkey] = cdff
				}
			}
			otc.OccupancyTypes[key] = curval
		} else {
			return errors.New("structures: occupancy type " + key + " doesn't currently exist.")
		}
	}
	return nil
}
func (otc OccupancyTypesContainer) OcctypeReport() ([]byte, error) {
	var ret string
	ret += fmt.Sprintf("|%v| %v| %v| %v| %v|\n", "occtype", "componenttype", "compoundhazard", "damageDriver", "source")
	ret += fmt.Sprintf("|%v| %v| %v| %v| %v|\n", "-----", "-----", "-----", "-----", "-----")
	for occtype, occvalue := range otc.OccupancyTypes {
		for componenttype, componentvalue := range occvalue.ComponentDamageFunctions {
			for functiontype, functionvalue := range componentvalue.DamageFunctions {
				row := fmt.Sprintf("|%v| %v| %v| %v| %v|\n", occtype, componenttype, functiontype.String(), functionvalue.DamageDriver.String(), functionvalue.Source)
				ret += row
			}
		}
	}
	return []byte(ret), nil
}

//GetComponentDamageFunctionForHazard provides a hazard specific damage function for a component (e.g. structure, content, car, or other)
func (o OccupancyTypeDeterministic) GetComponentDamageFunctionForHazard(component string, h hazards.HazardEvent) (DamageFunction, error) {
	c, cok := o.ComponentDamageFunctions[component]
	if cok {
		r, rok := c.DamageFunctions[h.Parameters()]
		if rok {
			return r, nil
		} else {
			return c.DamageFunctions[hazards.Default], nil //errors.New("using default damage function")
		}
	}
	return DamageFunction{}, errors.New("component does not exist for this occupancy type")
}

//UncertaintyOccupancyTypeSampler provides the pattern for an OccupancyTypeStochastic to produce an OccupancyTypeDeterministic
type UncertaintyOccupancyTypeSampler interface {
	SampleOccupancyType(rand int64) OccupancyTypeDeterministic
	CentralTendencyOccupancyType() OccupancyTypeDeterministic
}

//SampleOccupancyType implements the UncertaintyOccupancyTypeSampler on the OccupancyTypeStochastic interface.
func (o OccupancyTypeStochastic) SampleOccupancyType(seed int64) OccupancyTypeDeterministic {
	r := rand.New(rand.NewSource(seed))
	//iterate through damage function family
	cm := make(map[string]DamageFunctionFamily)
	for ck, cv := range o.ComponentDamageFunctions { //components
		hm := make(map[hazards.Parameter]DamageFunction)
		var cdf = DamageFunctionFamily{DamageFunctions: hm}
		for k, v := range cv.DamageFunctions { //hazards
			df := DamageFunction{}
			df.DamageDriver = v.DamageDriver
			df.Source = v.Source
			df.DamageFunction = samplePairedDataValueSampler(r, v.DamageFunction)
			cdf.DamageFunctions[k] = df
		}
		cm[ck] = cdf
	}
	return OccupancyTypeDeterministic{Name: o.Name, ComponentDamageFunctions: cm}
}
func samplePairedDataValueSampler(r *rand.Rand, df interface{}) paireddata.PairedData {
	retval, ok := df.(paireddata.PairedData)
	if ok {
		return retval
	}
	//must be uncertain
	retval2, ok2 := df.(paireddata.UncertaintyValueSamplerSampler)
	if ok2 {
		pd := retval2.SampleValueSampler(r.Float64())
		pd2, ok3 := pd.(paireddata.PairedData)
		if ok3 {
			pd2.ForceMonotonicInRange(0.0, 100.0)
			return pd2
		}
		return pd2 //this is actually not ok, but i dont have many options atm.
	}
	return retval
}
func centralTendencyPairedDataValueSampler(df interface{}) paireddata.PairedData {
	retval, ok := df.(paireddata.PairedData)
	if ok {
		return retval
	}
	//must be uncertain
	retval2, ok2 := df.(paireddata.UncertaintyValueSamplerSampler)
	if ok2 {
		ret := retval2.CentralTendency()
		retpd, ok3 := ret.(paireddata.PairedData)
		if ok3 {
			return retpd
		}
		return paireddata.PairedData{}
	}
	return retval
}

//CentralTendency implements the UncertaintyOccupancyTypeSampler on the OccupancyTypeStochastic interface.
func (o OccupancyTypeStochastic) CentralTendency() OccupancyTypeDeterministic {
	//iterate through damage function family
	cm := make(map[string]DamageFunctionFamily)
	for ck, cv := range o.ComponentDamageFunctions { //components
		hm := make(map[hazards.Parameter]DamageFunction)
		var cdf = DamageFunctionFamily{DamageFunctions: hm}
		for k, v := range cv.DamageFunctions { //hazards
			df := DamageFunction{}
			df.DamageDriver = v.DamageDriver
			df.Source = v.Source
			df.DamageFunction = centralTendencyPairedDataValueSampler(v.DamageFunction)
			cdf.DamageFunctions[k] = df
		}
		cm[ck] = cdf
	}
	return OccupancyTypeDeterministic{Name: o.Name, ComponentDamageFunctions: cm}
}
func arrayToDetermnisticDistributions(vals []float64) []statistics.ContinuousDistribution {
	dists := make([]statistics.ContinuousDistribution, len(vals))
	for i, v := range vals {
		dists[i] = statistics.DeterministicDistribution{Value: v}
	}
	return dists
}
