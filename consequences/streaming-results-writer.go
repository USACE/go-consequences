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

func InitStreamingResultsWriterFromFile(filepath string) (streamingResultsWriter, error) {
	w, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return streamingResultsWriter{}, err
	}
	return streamingResultsWriter{filepath: filepath, w: w}, nil
}
func InitStreamingResultsWriter(w io.Writer) streamingResultsWriter {
	return streamingResultsWriter{filepath: "not applicapble", w: w}
}
func (srw streamingResultsWriter) Write(r Result) {
	b, _ := r.MarshalJSON()
	s := string(b) + "\n"
	fmt.Fprint(srw.w, s)
}
func (srw streamingResultsWriter) Close() {
	w2, ok := srw.w.(io.WriteCloser)
	if ok {
		w2.Close()
	}
}
