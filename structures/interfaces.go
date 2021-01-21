package structures

type StructureProvider interface{
	ProvideStructures() StructureStochastic
}