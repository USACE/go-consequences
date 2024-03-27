package resultswriters

import (
	"strings"

	"github.com/USACE/go-consequences/consequences"
)

type VirtualResultsWriter struct {
	data                 strings.Builder
	headerHasBeenWritten bool
}

func InitVirtualResultsWriter() *VirtualResultsWriter {
	var b strings.Builder
	return &VirtualResultsWriter{data: b}
}
func (srw *VirtualResultsWriter) Write(r consequences.Result) {
	if !srw.headerHasBeenWritten {
		srw.data.WriteString("{\"consequences\":[")
		srw.headerHasBeenWritten = true
	}
	b, _ := r.MarshalJSON()
	s := string(b) + ","
	srw.data.WriteString(s)
}
func (srw *VirtualResultsWriter) Close() {
	srw.data.WriteString("]}")
}
func (srw *VirtualResultsWriter) Bytes() []byte {
	srw.Close()
	return []byte(srw.data.String())
}
