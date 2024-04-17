package lifeloss

import (
	"testing"

	"github.com/USACE/go-consequences/hazards"
)

func Test_RescDamUnanchoredWoodStability(t *testing.T) {

	sc := RescDamWoodUnanchored

	h := hazards.DepthandDVEvent{}
	h.SetDV(30)
	h.SetDepth(100)
	result := sc.Evaluate(h)
	if result != Stable {
		t.Fail()
	}

	h.SetDV(35)
	h.SetDepth(100)
	result = sc.Evaluate(h)
	if result != Collapsed {
		t.Fail()
	}
}
func Test_RescDamAnchoredWoodStability(t *testing.T) {

	sc := RescDamWoodAnchored

	h := hazards.DepthandDVEvent{}
	h.SetDV(35)
	h.SetDepth(100)
	result := sc.Evaluate(h)
	if result != Stable {
		t.Fail()
	}

	h.SetDV(75.4)
	h.SetDepth(100)
	result = sc.Evaluate(h)
	if result != Collapsed {
		t.Fail()
	}
}
func Test_RescDamConcreteMasonarySteelStability(t *testing.T) {

	sc := RescDamMasonryConcreteBrick

	h := hazards.DepthandDVEvent{}
	h.SetDV(35)
	h.SetDepth(100)
	result := sc.Evaluate(h)
	if result != Stable {
		t.Fail()
	}

	h.SetDV(75.4)
	h.SetDepth(10)
	result = sc.Evaluate(h)
	if result != Collapsed {
		t.Fail()
	}

	h.SetDV(75.4)
	h.SetDepth(20)
	result = sc.Evaluate(h)
	if result != Stable { //velocity is computed as not being high enough because depth was so high.
		t.Fail()
	}
}
