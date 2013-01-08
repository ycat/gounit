package assert

import (
	"testing"
	"runtime"
	"os"
	"path/filepath"
	"fmt"
	"reflect"
)

type Assert struct {
	t *testing.T
}

func New(t *testing.T) Assert {
	return Assert {t}
}

func relativePath(file string) (relpath string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	relpath, err = filepath.Rel(wd, file)
	return
}

func (as *Assert) reportError(cd int, format string, a ...interface{}) {
	_, file, line, ok := runtime.Caller(cd)
	if ok {
		file, _ = relativePath(file)
	} else {
		file = "?"
	}
	msg := fmt.Sprintf(format, a...)
	as.t.Errorf("%v:%v: %v", file, line, msg)
}


type varArgsFunc func (...interface{})
type Predicate func (...interface{}) bool

func AllNil() Predicate {
	return func(ops ...interface{}) bool {
		for index := 0; index < len(ops); index ++ {
			if ops[index] != nil {
				return false
			}
		}
		return true
	}
}

func Not(pred Predicate) Predicate {
	return func(ops ...interface{}) bool {
		return !pred(ops...)
	}
}


func Eq(equalFunc (func (interface{}, interface{}) bool)) Predicate {
	return func(ops ...interface{}) bool {
		if len(ops) <= 1 {
			panic(fmt.Sprintf("too few arguments for Equivalance: %v", ops))
		}
		op0 := ops[0]
		for index := 1; index < len(ops); index ++ {
			if !equalFunc(op0, ops[index]) {
				return false
			}
		}
		return true
	}
}


func (as *Assert) Assert(pred Predicate, format string, a ...interface{}) varArgsFunc {
	return func(ops ...interface{}) {
		defer func() {
			if r := recover(); r != nil {
				as.reportError(10, "panic with following parameters %v. Panic: %v", a, r)
			}
		}()
		if !pred(ops...) {
			as.reportError(2, format, a...)
		}
	}
}

func (as *Assert) IsAllNil(format string, a ...interface{}) varArgsFunc  {
	return as.Assert(AllNil(), format, a...)
}

func (as *Assert) ExistNotNil(format string, a ...interface{}) varArgsFunc  {
	return as.Assert(Not(AllNil()), format, a...)
}

func shallowEqual(a interface{}, b interface{}) bool {
	return a == b
}

func (as *Assert) Equal(format string, a ...interface{}) varArgsFunc {
	return as.Assert(Eq(shallowEqual), format, a...)
}

func (as *Assert) NotEqual(format string, a ...interface{}) varArgsFunc {
	return as.Assert(Not(Eq(shallowEqual)), format, a...)
}

func (as *Assert) DeepEqual(format string, a ...interface{}) varArgsFunc {
	return as.Assert(Eq(reflect.DeepEqual), format, a...)
}

func (as *Assert) NotDeepEqual(format string, a ...interface{}) varArgsFunc {
	return as.Assert(Not(Eq(reflect.DeepEqual)), format, a...)
}