# go-consequences

An Economic Consequences Library written in Golang

This library provides tools to support evaluating natural hazards interacting with concequence receptors. An example would be a flood represented by depth interacting with a residential structure to produce an estimate of economic losses at the residential structure.

## Packages
The following packages are in this library:
- census
- compute
- consequences
- crops
- hazard_providers
- hazards
- nsi
- paireddata

### census
The census package contains a map of state FIP codes to county FIP codes to support consequences computations and iteration across the entire United States.

### compute
The compute package is intended to support the Computable interface which takes a set of request args describing the compute and produces a simulation summary.

### consequences
The consequences package contains the interfaces behind a consequences receptor and the consequences results. It also contains a map of the current occupancy types supported by the National Structure Inventory (NSI) and an implementation of a structure type for the consequences receptor. This facilitates flood consequence estimation for any structure in the NSI.

### crops
The crops package contains the logic for agricultural consequences leveraging the NASS CDL data.

### hazard_providers
The hazard providers package is designed to fulfil the interface HazardProvider which suports multiple hazard types from multiple hazard sources.

### hazards
Various hazards are stored in the hazards package, the primary hazard under review is flood, but fire is also functional for the structure consequence receptor

### nsi
The NSI package provides access to the NSI api bounding box endpoint so that structures can be retrieved for the extent of a grid representing the area of interest ofor a hazard.

### paireddata
The paireddata object provides a linear interpolation of x and y data. This is used in the representation of depth damage relationships for the occupancy types described by the NSI structures. 

## Testing
Tests have been developed for most of the code related to flood damage estimation. The tests can be compiled using the general calls listed below on a package level. 

```
C:\Examples\Go_Consequences>go test ./paireddata -c
C:\Examples\Go_Consequences>.\paireddata.test -test.v