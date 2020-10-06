package compute

type StructureSimulation struct {
	//some sort of input
	//some sort of progress
	//

}

type Computeable interface {
	Compute() // what arguments?
}

type ProgressReportable interface {
	ReportProgress() float64
}
