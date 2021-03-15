package structureprovider

import (
	"github.com/USACE/go-consequences/geography"
	"github.com/USACE/go-consequences/structures"
)

type StreamProvider interface {
	ByFips(fipscode string, sp StreamProcessor)
	ByBbox(bbox geography.BBox, sp StreamProcessor)
}
type StreamProcessor func(str structures.StructureStochastic) //consequences.Receptor)
