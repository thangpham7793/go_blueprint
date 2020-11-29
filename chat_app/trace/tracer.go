package trace

import (
	"fmt"
	"io"
)

//Tracer describes an object capable of tracing events throughout code
type Tracer interface {
	Trace(...interface{})
}

//New inits a tracer
func New(w io.Writer) Tracer {
	return &tracer{w}
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Println(t.out)
}
