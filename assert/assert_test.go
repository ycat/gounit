package assert

import "testing"

func TestEq(t* testing.T) {
	as := New(t)
	if result, _ := Eq()(1, 1); !result {
		as.reportError(1, "Eq()(1, 1) expects true but got false")
	}
	if result, _ := Eq()(1, 1, 1); !result {
		as.reportError(1, "Eq()(1, 1, 1) expects true but got false")
	}
	if result, _ := Eq()(1, 2); result {
		as.reportError(1, "Eq()(1, 2) expects false but got true")
	}
	if result, _ := Eq()(1, 2, 1); result {
		as.reportError(1, "Eq()(1, 2, 1) expects false but got true")
	}
	as.Equal("(%v, %v) should be equivalent", 1, 1)(1, 1)
	as.Equal("(%v, %v, %v) should be equivalent", 1, 1, 1)(1, 1, 1)
	as.NotEqual("(%v, %v) should not be equivalent", 1, 2)(1, 2)
	as.NotEqual("(%v, %v, %v) should not be equivalent", 1, 2, 1)(1, 2, 1)

	a1 := [3]byte{1, 2, 3}
	a2 := [3]byte{1, 2, 3}
	as.Equal("(%v, %v) should be equivalent", a1, a2)(a1, a2)

	a3 := [3]byte{1, 3, 2}
	a4 := [3]byte{1, 2, 3}
	as.NotEqual("(%v, %v) should not be equivalent", a3, a4)(a3, a4)

	s1 := a1[1:]
	s2 := a2[1:]
	as.Equal("(%v, %v) should be equivalent", s1, s2)(s1, s2)

	s3 := a3[:]
	s4 := a4[:]
	as.NotEqual("(%v, %v) should not be equivalent", s3, s4)(s3, s4)

	a5 := [3][3]byte{ {1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	a6 := [3][3]byte{ {1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	as.Equal("(%v, %v) should be equivalent", a5, a6)(a5, a6)

	s5 := a5[:]
	s6 := a6[:]
	as.Equal("(%v, %v) should be equivalent", s5, s6)(s5, s6)

	s7 := a5[0:2]
	s8 := a6[1:3]
	as.NotEqual("(%v, %v) should not be equivalent", s7, s8)(s7, s8)
	
	a7 := [2][]byte{s1, s2}
	a8 := [2][]byte{s1, s2}
	as.Equal("(%v, %v) should be equivalent", a7, a8)(a7, a8)
	
	s9 := a7[:]
	s10 := a8[:]
	as.Equal("(%v, %v) should be equivalent", s9, s10)(s9, s10)

	m1 := make(map[string] int)
	m2 := make(map[string] int)
	as.Equal("(%v, %v) should be equivalent", m1, m2)(m1, m2)
	
	m3 := map[string] string {"k1":"v1", "k2":"v2"}
	m4 := map[string] string {"k2":"v2", "k1":"v1"}
	as.Equal("(%v, %v) should be equivalent", m3, m4)(m3, m4)

	m5 := map[string] string {"k1":"v1", "k2":"v3"}
	m6 := map[string] string {"k2":"v2", "k1":"v1"}
	as.NotEqual("(%v, %v) should be equivalent", m5, m6)(m5, m6)

	m7 := map[string] string {"k1":"v1", "k3":"v2"}
	m8 := map[string] string {"k2":"v2", "k1":"v1"}
	as.NotEqual("(%v, %v) should be equivalent", m7, m8)(m7, m8)

	m9 := map[string] []byte {"k1":s1, "k2":s2}
	m10 := map[string] []byte {"k1":s1, "k2":s2}
	as.Equal("(%v, %v) should be equivalent", m9, m10)(m9, m10)

	m11 := map[string] []byte {"k1":s1, "k2":s3}
	m12 := map[string] []byte {"k1":s1, "k2":s2}
	as.NotEqual("(%v, %v) should not be equivalent", m11, m12)(m11, m12)
	
}

func TestNil(t* testing.T) {
	as := New(t)
	if result, _ := AllNil()(nil); !result {
		as.reportError(1, "AllNil(nil) expects true but got false")
	}
	if result, _ := AllNil()(nil, nil); !result {
		as.reportError(1, "AllNil(nil, nil) expects true but got false")
	}
	if result, _ := AllNil()(1); result {
		as.reportError(1, "AllNil(1) expects false but got true")
	}
	if result, _ := AllNil()(1, nil); result {
		as.reportError(1, "AllNil(1, nil) expects false but got true")
	}
	as.IsAllNil("IsNil(nil) should be true")(nil)
	as.IsAllNil("IsNil(nil, nil) should be true")(nil, nil)
}