package structureprovider

import (
	"fmt"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/structures"
	"github.com/dewberry/gdal"
)

func StructureSchema() []string {
	s := make([]string, 10)
	s[0] = "fd_id"
	s[1] = "cbfips"
	s[2] = "x"
	s[3] = "y"
	s[4] = "st_damcat"
	s[5] = "occtype"
	s[6] = "val_struct"
	s[7] = "val_cont"
	s[8] = "found_ht"
	s[9] = "found_type"
	return s
}

func OptionalSchema() []string {
	s := make([]string, 7)
	s[0] = "num_story"
	s[1] = "pop2amu65"
	s[2] = "pop2amo65"
	s[3] = "pop2pmu65"
	s[4] = "pop2pmo65"
	s[5] = "ground_elv"
	s[6] = "bldgtype"
	s[7] = "firmzone"
	return s
}

func featuretoStructure(
	f *gdal.Feature,
	m map[string]structures.OccupancyTypeStochastic,
	defaultOcctype structures.OccupancyTypeStochastic,
	idxs []int,
	oidxs []int,
) (structures.StructureStochastic, error) {
	defer f.Destroy()
	s := structures.StructureStochastic{}
	s.Name = fmt.Sprintf("%v", f.FieldAsInteger(idxs[0]))
	OccTypeName := f.FieldAsString(idxs[5])
	var occtype = defaultOcctype
	//dont have access to foundation type in the structure schema yet.
	if idxs[9] > 0 {
		if otf, okf := m[OccTypeName+"-"+f.FieldAsString(idxs[9])]; okf {
			occtype = otf
		} else {
			if ot, ok := m[OccTypeName]; ok {
				occtype = ot
			} else {
				occtype = defaultOcctype
				msg := "Using default " + OccTypeName + " not found"
				fmt.Println(msg)
				//return s, errors.New(msg)
			}
		}
	} else {
		if ot, ok := m[OccTypeName]; ok {
			occtype = ot
		} else {
			occtype = defaultOcctype
			msg := "Using default " + OccTypeName + " not found"
			fmt.Println(msg)
			//return s, errors.New(msg)
		}
	}

	s.OccType = occtype
	s.CBFips = f.FieldAsString(idxs[1])
	s.X = f.FieldAsFloat64(idxs[2])
	s.Y = f.FieldAsFloat64(idxs[3])
	s.DamCat = f.FieldAsString(idxs[4])
	s.FoundType = f.FieldAsString(idxs[9])
	s.StructVal = consequences.ParameterValue{Value: f.FieldAsFloat64(idxs[6])}
	s.ContVal = consequences.ParameterValue{Value: f.FieldAsFloat64(idxs[7])}
	s.FoundHt = consequences.ParameterValue{Value: f.FieldAsFloat64(idxs[8])}
	if oidxs[0] != -1 {
		s.NumStories = int32(f.FieldAsInteger(oidxs[0]))
	}
	if oidxs[1] != -1 {
		s.Pop2amu65 = int32(f.FieldAsInteger(oidxs[1]))
	}
	if oidxs[2] != -1 {
		s.Pop2amo65 = int32(f.FieldAsInteger(oidxs[2]))
	}
	if oidxs[3] != -1 {
		s.Pop2pmu65 = int32(f.FieldAsInteger(oidxs[3]))
	}
	if oidxs[4] != -1 {
		s.Pop2pmo65 = int32(f.FieldAsInteger(oidxs[4]))
	}
	if oidxs[5] != -1 {
		s.GroundElevation = f.FieldAsFloat64(oidxs[5])
	}
	if oidxs[6] != -1 {
		s.ConstructionType = f.FieldAsString(oidxs[6])
	}
	if oidxs[7] != -1 {
		s.FirmZone = f.FieldAsString(oidxs[7])
	}
	return s, nil
}

func swapOcctypeMap(
	m map[string]structures.OccupancyTypeStochastic,
) map[string]structures.OccupancyTypeDeterministic {
	m2 := make(map[string]structures.OccupancyTypeDeterministic)
	for name, ot := range m {
		m2[name] = ot.CentralTendency()
	}
	return m2
}

func featuretoDeterministicStructure(
	f *gdal.Feature,
	m map[string]structures.OccupancyTypeDeterministic,
	defaultOcctype structures.OccupancyTypeDeterministic,
	idxs []int,
	oidxs []int,
) (structures.StructureDeterministic, error) {
	defer f.Destroy()
	s := structures.StructureDeterministic{}
	s.Name = fmt.Sprintf("%v", f.FieldAsInteger(idxs[0]))
	OccTypeName := f.FieldAsString(idxs[5])
	var occtype = defaultOcctype
	//dont have access to foundation type in the structure schema yet.
	if idxs[9] > 0 {
		if otf, okf := m[OccTypeName+"-"+f.FieldAsString(idxs[9])]; okf {
			occtype = otf
		} else {
			if ot, ok := m[OccTypeName]; ok {
				occtype = ot
			} else {
				occtype = defaultOcctype
				msg := "Using default " + OccTypeName + " not found"
				fmt.Println(msg)
				//return s, errors.New(msg)
			}
		}
	} else {
		if ot, ok := m[OccTypeName]; ok {
			occtype = ot
		} else {
			occtype = defaultOcctype
			msg := "Using default " + OccTypeName + " not found"
			fmt.Println(msg)
			//return s, errors.New(msg)
		}
	}

	s.OccType = occtype
	s.CBFips = f.FieldAsString(idxs[1])
	s.X = f.FieldAsFloat64(idxs[2])
	s.Y = f.FieldAsFloat64(idxs[3])
	s.DamCat = f.FieldAsString(idxs[4])
	s.StructVal = f.FieldAsFloat64(idxs[6])
	s.ContVal = f.FieldAsFloat64(idxs[7])
	s.FoundHt = f.FieldAsFloat64(idxs[8])
	s.FoundType = f.FieldAsString(idxs[9])
	if oidxs[0] != -1 {
		s.NumStories = int32(f.FieldAsInteger(oidxs[0]))
	}
	if oidxs[1] != -1 {
		s.Pop2amu65 = int32(f.FieldAsInteger(oidxs[1]))
	}
	if oidxs[2] != -1 {
		s.Pop2amo65 = int32(f.FieldAsInteger(oidxs[2]))
	}
	if oidxs[3] != -1 {
		s.Pop2pmu65 = int32(f.FieldAsInteger(oidxs[3]))
	}
	if oidxs[4] != -1 {
		s.Pop2pmo65 = int32(f.FieldAsInteger(oidxs[4]))
	}
	if oidxs[5] != -1 {
		s.GroundElevation = f.FieldAsFloat64(oidxs[5])
	}
	if oidxs[6] != -1 {
		s.ConstructionType = f.FieldAsString(oidxs[6])
	}
	if oidxs[7] != -1 {
		s.FirmZone = f.FieldAsString(oidxs[7])
	}
	return s, nil
}
