package consequences

import (
	"fmt"
	"io"
	"os"
)

type streamingResultsWriter struct {
	filepath string
	w        io.Writer
}

func InitStreamingResultsWriterFromFile(filepath string) streamingResultsWriter {
	w, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	return streamingResultsWriter{filepath: filepath, w: w}
}
func InitStreamingResultsWriter(w io.Writer) streamingResultsWriter {
	return streamingResultsWriter{filepath: "not applicapble", w: w}
}
func (srw streamingResultsWriter) Write(r Result) {
	b, _ := r.MarshalJSON()
	s := string(b) + "\n"
	fmt.Fprintf(srw.w, s)
}
func (srw streamingResultsWriter) Close() {
	w2, ok := srw.w.(io.WriteCloser)
	if ok {
		w2.Close()
	}
}
