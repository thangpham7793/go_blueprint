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

type nilTracer struct{}

func (n *nilTracer) Trace(a ...interface{}) {}

//Off creates a nil tracer that ignores calls to Trace()
func Off() Tracer {
	return &nilTracer{}
}
