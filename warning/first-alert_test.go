package warning

import (
	"fmt"
	"testing"
)

func TestComputeWarningDiffusion(t *testing.T) {
	result := ComputeCurve(WarningSystemParameters{B: 95.0, C: .06})
	for i, xval := range result.Xvals {
		fmt.Printf("%v,%v\n", xval, result.Yvals[i])
	}
}
