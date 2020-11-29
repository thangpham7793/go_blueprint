package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Return from New should not be nil")
	} else {
		msg := "Hello from package!\n"
		tracer.Trace(msg)
		if buf.String() != msg {
			t.Errorf("Trace should not write '%s'.", buf.String())
		}
	}
}
