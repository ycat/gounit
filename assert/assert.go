package assert

import (
	"testing"
	"runtime"
	"fmt"
)

type Assert struct {
	t *testing.T
}

func New(t *testing.T) Assert {
	return Assert {t}
}

func (as *Assert) reportError(errorMsg string, cd int) {
	_, file, line, ok := runtime.Caller(cd); 
	if !ok {
		file = "?"
		line = 0
	} 
	as.t.Errorf("%v:%v: %v", file, line, errorMsg)
}

func (as *Assert) NoError(err error) {
	if err != nil {
		errorMsg := fmt.Sprintf("%v", err)
		as.reportError(errorMsg, 2)
	}		
}

func (as *Assert) Equals(expected interface{}, actual interface{}) {
	if actual != expected {
		errorMsg := fmt.Sprintf("%v != %v", expected, actual)
		as.reportError(errorMsg, 2)
	}
}