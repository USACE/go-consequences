package crops

import (
	"fmt"
	"testing"

	"github.com/USACE/go-consequences/consequences"
)

func Test_StreamProcessor(t *testing.T) {
	nassSp := NassCropProvider{Year: "2018"}
	nassSp.ByFips("19017", func(r consequences.Receptor) {
		c, ok := r.(Crop)
		if ok {
			fmt.Println(c.name)
		}
	})
}
