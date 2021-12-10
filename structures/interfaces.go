package structures

type RawDFProvider interface {
	Elements(occType string) // Get all substructs under specific occType key
}

type DamageCategory struct {
	Name        string      `json:"Name"`
	Description interface{} `json:"Description"`
	Rebuild     string      `json:"Rebuild"`
	CostFactor  string      `json:"CostFactor"`
}

type FoundationHeightUncertainty struct {
	None struct {
		Value string `json:"_value"`
	} `json:"None"`
}

type StructureUncertainty struct {
	None struct {
		Value string `json:"_value"`
	} `json:"None"`
}

type ContentUncertainty struct {
	None struct {
		Value string `json:"_value"`
	} `json:"None"`
}

type OtherUncertainty struct {
	None struct {
		Value string `json:"_value"`
	} `json:"None"`
}

type VehicleUncertainty struct {
	None struct {
		Value string `json:"_value"`
	} `json:"None"`
}

type StructureDD struct {
	CalculateDamage       string `json:"CalculateDamage"`
	MonotonicCurveUSingle struct {
		UncertaintyType string `json:"UncertaintyType"`
		Ordinate        []struct {
			X     string `json:"X"`
			Value string `json:"_value"`
		} `json:"Ordinate"`
	} `json:"MonotonicCurveUSingle"`
}

type ContentDD struct {
	CalculateDamage       string `json:"CalculateDamage"`
	MonotonicCurveUSingle struct {
		UncertaintyType string `json:"UncertaintyType"`
		Ordinate        []struct {
			X     string `json:"X"`
			Value string `json:"_value"`
		} `json:"Ordinate"`
	} `json:"MonotonicCurveUSingle"`
}

type OtherDD struct {
	CalculateDamage       string `json:"CalculateDamage"`
	MonotonicCurveUSingle struct {
		UncertaintyType string `json:"UncertaintyType"`
		Ordinate        struct {
			X     string `json:"X"`
			Value string `json:"_value"`
		} `json:"Ordinate"`
	} `json:"MonotonicCurveUSingle"`
}

type VehicleDD struct {
	CalculateDamage string `json:"CalculateDamage"`
}

type Prototype struct {
	Name                        string                      `json:"Name"`
	Description                 string                      `json:"Description"`
	DamageCategory              DamageCategory              `json:"DamageCategory"`
	FoundationHeightUncertainty FoundationHeightUncertainty `json:"FoundationHeightUncertainty"`
	StructureUncertainty        StructureUncertainty        `json:"StructureUncertainty"`
	ContentUncertainty          ContentUncertainty          `json:"ContentUncertainty"`
	OtherUncertainty            OtherUncertainty            `json:"OtherUncertainty"`
	VehicleUncertainty          VehicleUncertainty          `json:"VehicleUncertainty"`
	StructureDD                 StructureDD                 `json:"StructureDD"`
	ContentDD                   ContentDD                   `json:"ContentDD"`
	OtherDD                     OtherDD                     `json:"OtherDD"`
	VehicleDD                   VehicleDD                   `json:"VehicleDD"`
}

type RawDFStruct struct {
	OccTypes struct {
		Prototypes []Prototype `json:"OccupancyType"`
	} `json:"OccTypes"`
}

// Map with OccupancyType string as index
type DFStore map[string]Prototype
