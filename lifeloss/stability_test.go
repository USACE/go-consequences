package lifeloss_test

import (
	"fmt"
	"testing"

	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/lifeloss"
)

func ExampleStabilityEvaluation() {

	sc := lifeloss.RescDamWoodUnanchored

	h := hazards.DepthandDVEvent{}
	h.SetDV(30)
	h.SetDepth(100)
	result := sc.Evaluate(h)
	fmt.Println(result)

	h.SetDV(35)
	h.SetDepth(100)
	result = sc.Evaluate(h)
	fmt.Println(result)
}

func Test_RescDamUnanchoredWoodStability(t *testing.T) {

	sc := lifeloss.RescDamWoodUnanchored

	h := hazards.DepthandDVEvent{}
	h.SetDV(30)
	h.SetDepth(100)
	result := sc.Evaluate(h)
	if result != lifeloss.Stable {
		t.Fail()
	}

	h.SetDV(35)
	h.SetDepth(100)
	result = sc.Evaluate(h)
	if result != lifeloss.Collapsed {
		t.Fail()
	}
}
func Test_RescDamAnchoredWoodStability(t *testing.T) {

	sc := lifeloss.RescDamWoodAnchored

	h := hazards.DepthandDVEvent{}
	h.SetDV(35)
	h.SetDepth(100)
	result := sc.Evaluate(h)
	if result != lifeloss.Stable {
		t.Fail()
	}

	h.SetDV(75.4)
	h.SetDepth(100)
	result = sc.Evaluate(h)
	if result != lifeloss.Collapsed {
		t.Fail()
	}
}
func Test_RescDamConcreteMasonarySteelStability(t *testing.T) {

	sc := lifeloss.RescDamMasonryConcreteBrick

	h := hazards.DepthandDVEvent{}
	h.SetDV(35)
	h.SetDepth(100)
	result := sc.Evaluate(h)
	if result != lifeloss.Stable {
		t.Fail()
	}

	h.SetDV(75.4)
	h.SetDepth(10)
	result = sc.Evaluate(h)
	if result != lifeloss.Collapsed {
		t.Fail()
	}

	h.SetDV(75.4)
	h.SetDepth(20)
	result = sc.Evaluate(h)
	if result != lifeloss.Stable { //velocity is computed as not being high enough because depth was so high.
		t.Fail()
	}
}
