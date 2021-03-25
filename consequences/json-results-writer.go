package consequences

import (
	"fmt"
	"io"
	"os"
)

type jsonResultsWriter struct {
	filepath             string
	w                    io.Writer
	headerHasBeenWritten bool
}

func InitJsonResultsWriterFromFile(filepath string) jsonResultsWriter {
	w, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	return jsonResultsWriter{filepath: filepath, w: w}
}
func InitJsonResultsWriter(w io.Writer) jsonResultsWriter {
	return jsonResultsWriter{filepath: "not applicapble", w: w}
}
func (srw jsonResultsWriter) Write(r Result) {
	if !srw.headerHasBeenWritten {
		fmt.Fprintf(srw.w, "{\"consequences\":[")
		srw.headerHasBeenWritten = true
	}
	b, _ := r.MarshalJSON()
	s := string(b) + ","
	fmt.Fprintf(srw.w, s)
}
func (srw jsonResultsWriter) Close() {
	fmt.Fprintf(srw.w, "]}")
	w2, ok := srw.w.(io.WriteCloser)
	if ok {
		w2.Close()
	}
}
