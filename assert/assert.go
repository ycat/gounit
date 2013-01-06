package assert

import (
	"testing"
	"runtime"
	"os"
	"path/filepath"
	"errors"
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
type Predicate func (...interface{}) (bool, error)

func AllNil() Predicate {
	return func(ops ...interface{}) (result bool, err error) {
		result = true
		for index := 0; index < len(ops); index ++ {
			if ops[index] != nil {
				result = false
				return
			}
		}
		return 
	}
}

func Not(pred Predicate) Predicate {
	return func(ops ...interface{}) (bool, error) {
		result, err := pred(ops...)
		return !result, err
	}
}


func equalSlice(op1 interface{}, op2 interface{}) bool {
	s1 := reflect.ValueOf(op1)
	s2 := reflect.ValueOf(op2)
	if s1.Len() == s2.Len() {
		for index := 0; index < s1.Len(); index ++ {
			if !equal(s1.Index(index).Interface(), s2.Index(index).Interface()) {
				return false
			}
		}
		return true
	} 
	return false
}

func equalMap(op1 interface{}, op2 interface{}) bool {
	m1 := reflect.ValueOf(op1)
	m2 := reflect.ValueOf(op2)
	if m1.Len() == m2.Len() {
		keys := m1.MapKeys()
		for _, key := range keys {
			value1 := m1.MapIndex(key)
			value2 := m2.MapIndex(key)
			if !value2.IsValid() {
				return false
			}
			if !equal(value1.Interface(), value2.Interface()) {
				return false
			}
		}
		return true
	}
	return false
}

func alwaysNotEqual(op1 interface{}, op2 interface{}) bool {
	return false
}

func selectEqualFunction(op1 interface{}, op2 interface{}) (func (interface{}, interface{}) bool) {
	kind1 := reflect.TypeOf(op1).Kind()
	kind2 := reflect.TypeOf(op2).Kind()
	
	isSlice1 := kind1 == reflect.Slice || kind1 == reflect.Array
	isSlice2 := kind2 == reflect.Slice || kind2 == reflect.Array

	if isSlice1 != isSlice2 {
		return alwaysNotEqual
	}
	
	if isSlice1 && isSlice2 {
		return equalSlice
	}

	isMap1 := kind1 == reflect.Map
	isMap2 := kind2 == reflect.Map
	if isMap1 != isMap2 {
		return alwaysNotEqual
	}
	
	if isMap1 && isMap2 {
		return equalMap
	}
	
	return equalDefault
}

func equalDefault(op1 interface{}, op2 interface{}) bool {
	return op1 == op2
}

func equal(op1 interface{}, op2 interface{}) bool {
	return selectEqualFunction(op1, op2)(op1, op2)
}

func Eq() Predicate {
	return func(ops ...interface{}) (result bool, err error) {
		if len(ops) <= 1 {
			err = errors.New(fmt.Sprintf("too few arguments for Equivalance: %v", ops))
			return
		}
		op0 := ops[0]
		result = true
		for index := 1; index < len(ops); index ++ {
			if !equal(op0, ops[index]) {
				result = false
				return
			}
		}
		return
	}
}

func (as *Assert) Assert(pred Predicate, format string, a ...interface{}) varArgsFunc {
	return func(ops ...interface{}) {
		result, err := pred(ops...)
		if err != nil {
			as.reportError(2, "%v", err)
		} else if !result {
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

func (as *Assert) Equal(format string, a ...interface{}) varArgsFunc {
	return as.Assert(Eq(), format, a...)
}

func (as *Assert) NotEqual(format string, a ...interface{}) varArgsFunc {
	return as.Assert(Not(Eq()), format, a...)
}