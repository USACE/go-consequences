package consequences

import (
	"strings"
)

type VirtualResultsWriter struct {
	data strings.Builder
}

func InitVirtualResultsWriter() *VirtualResultsWriter {
	var b strings.Builder
	return &VirtualResultsWriter{data: b}
}
func (srw *VirtualResultsWriter) Write(b []byte) {
	srw.data.Write(b)
}
func (srw *VirtualResultsWriter) Close() {
	//do nothing
}
func (srw *VirtualResultsWriter) Bytes() []byte {
	return []byte(srw.data.String())
}
