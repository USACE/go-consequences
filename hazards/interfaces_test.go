package hazards

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMarshalMultiParameterJSON(t *testing.T) {
	d := Depth | Salinity
	b, _ := json.Marshal(d)
	fmt.Println(string(b))
	var d2 Parameter
	json.Unmarshal(b, &d2)
	fmt.Print(d2.String())
	if d2 != d {
		t.Error("unmarshal failed.")
	}
}
