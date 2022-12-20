package warning

import (
	"fmt"
	"testing"
)

func TestWarningIssuance(t *testing.T) {
	ld := LindellDist{a: 1.7, b: .6}
	for i := 0.0; i < 25.0; i++ {
		t := i * 15.0
		probIssued := ld.CDF(t)
		fmt.Printf("%v,%v\n", t, probIssued)
	}
}
