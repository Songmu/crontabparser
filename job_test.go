package crontab

import (
	"reflect"
	"testing"
)

func TestFieldsN(t *testing.T) {
	testCases := []struct {
		Name   string
		Input  string
		N      int
		Expect []string
	}{
		{"Simple", "  aaa  bbb	 ccc   DDDD  ", 3, []string{"aaa", "bbb", "ccc   DDDD"}},
		{"Just", "  aaa  bbb	 ccc   DDDD  ", 4, []string{"aaa", "bbb", "ccc", "DDDD"}},
		{"Big N", "  aaa  bbb	 ccc   DDDD  ", 5, []string{"aaa", "bbb", "ccc", "DDDD"}},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			out := fieldsN(tc.Input, tc.N)
			if !reflect.DeepEqual(out, tc.Expect) {
				t.Errorf("%s:\n   out: %#v\nexpect: %#v", tc.Name, out, tc.Expect)
			}
		})
	}
}
