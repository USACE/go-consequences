package indirecteconomics

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/USACE/go-consequences/census"
)

func Test_ComputeECAM(t *testing.T) {
	r, err := ComputeEcam("36", "049", 0.10163, 0.52977)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", r)
}
func Test_National_ECAM_Compute(t *testing.T) {
	fipsmap := census.StateToCountyFipsMap()
	//output := strings.Builder{}
	sucsesses := 0
	fails := 0
	total := 0
	w, err := os.OpenFile("/workspaces/Go_Consequences/data/ecamOutputreport.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	defer w.Close()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	for state, counties := range fipsmap {
		for _, county := range counties {
			total += 1
			_, err := ComputeEcam(state, county[2:5], 0.10163, 0.52977)
			if err != nil {
				fmt.Println(err, county)
				fmt.Fprintln(w, err, county)
				fails += 1
			} else {
				sucsesses += 1
				fmt.Println("success", county)
				fmt.Fprintln(w, "success", county)
			}
			time.Sleep(1 * time.Second)
		}
	}
	fmt.Println("successes:", sucsesses, "fails:", fails, "total:", total)
	fmt.Fprintln(w, "successes:", sucsesses, "fails:", fails, "total:", total)
}
