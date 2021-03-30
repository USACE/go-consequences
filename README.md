# go-consequences

An Economic Consequences Library written in Golang

This library provides tools to support evaluating natural hazards interacting with concequence receptors. An example would be a flood represented by depth interacting with a residential structure to produce an estimate of economic losses at the residential structure.

## Packages
The following packages are in this library:
- census
- compute
- consequences
- crops
- geography
- hazardproviders
- hazards
- paireddata
- structureprovider
- structure

### census
The census package contains a map of state FIP codes to county FIP codes to support consequences computations and iteration across the entire United States.

### compute
The compute package combines hazard providers with the streaming consequence receptor provider and a results writer to produce a compute iterating over all consequence receptors in the streaming consequence receptor provider within the geographic extent of the hazard provider and writes the results to the results writer selected.

### consequences
The consequences package contains the interfaces behind a consequences receptor and the consequences results. It contains the interface for the streaming consequences receptor. This facilitates consequence estimation for any implementation of consequence receptor. It contains the results writer interfaces (much like closable writer from go) to enable the writing of the atomic unit of the consequence receptor result. There are default implementations for a summary writer, streaming writer to any io.writer implementation, a geojson writer, and a json writer.

### crops
The crops package contains the logic for agricultural consequences leveraging the NASS CDL data. It implements the consequence receptor interface for crops. This package is a work in progress.

### hazardproviders
The hazard providers package defines the interface for HazardProvider. A hazard provider can provide a hazard for a point location and produces a hazard.HazardEvent. This package includes hazardprovider implementations for geotif files, specifically for depth events and duration and arrival time events.

### hazards
This package contains the inteface for HazardEvent which is an abstraction of any hazard. Various hazards are stored in the hazards package, the primary hazard under review is flood. A HazardEvent contains a parameter bitflag which describes what damage driving parameters are present in the hazardevent to quickly ascertain which types of consequence receptors might be vunerable (and to what severity).

### paireddata
The paireddata object provides a linear interpolation of x and y data. This is used in the representation of depth damage relationships for the occupancy types described by the NSI structures. 

### structureprovider
The structure provider package implements the streaming consequence receptor provider interfaces for the structure.StochasticStructure type. It contains implementations for the NSI api, geopackage, and shapefiles. This means a user can supply structure inventories using either geopackage or shapefile, or through the streaming services of the NSI.

### structure
The structure package contains the types for DeterministicStructure and StochasticStructure. The primary path of execution starts with a stochastic structure. A stochastic structure can be sampled to produce a deterministic structure. The deterministic structure (and the stochastic structure) implements the consequences receptor interface to produce a consequences result for a hazard event. The package also includes occupancy types for the standard USACE damage functions for residential structures (based on the EGMs) and additional damage functions for commercial industrial and public structures (mostly sourced from Galveston). Work is underway to add the NACCS coastal curves as well as some recent coastal curves produced by FEMA. OccupancyTypes are by default produced commensurate with their hazard and thier ability to operate stochastically. A deterministic occupancy type when asked to sample produces itself, a stochastic occupancy type curve samples its damage relationships, and produces a deterministic image of the occupancy type. 
Damage for hazardEvents with a depth parameter are implemented. If the hazard event is coastal in nature, a coastal damage function (if supplied for the occupancy type) is provided, if no damage function is specified for the hazardevent in question, the default (inland) depth damage relationship for the occupancy type is produced. 


## Testing
Tests have been developed for most of the code related to flood damage estimation. The tests can be compiled using the general calls listed below on a package level. 

```
C:\Examples\Go_Consequences>go test ./paireddata -c
C:\Examples\Go_Consequences>.\paireddata.test -test.v