package crontab

import "testing"

func TestDequote(t *testing.T) {
	testCases := []struct {
		Name    string
		Input   string
		Expect  string
		IsError bool
	}{
		{"No Quote", " abc d ", "abc d", false},
		{"Double Quote", `"abc 'd"`, "abc 'd", false},
		{"Single Quote", `'abc "d'`, `abc "d`, false},
		{"Inline Quote Chars", `abc ="d'`, `abc ="d'`, false},
		{"Invalid Double Quote", `"abc 'd"'`, "", true},
		{"Invalid Single Quote", `'abc 'd"`, "", true},
		{"Invalid Double Quote2", `"aa""`, "", true},
		{"Invalid Single Quote2", `'aa'a'`, "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			out, err := dequote(tc.Input)
			if tc.IsError {
				if err == nil {
					t.Errorf("%s: error should be occurred but nil", tc.Name)
				}
			} else {
				if out != tc.Expect {
					t.Errorf("%s: out=%q, expact=%q", tc.Name, out, tc.Expect)
				}
			}
		})
	}
}
