package checron

import "testing"

func TestNewEnv(t *testing.T) {
	testCases := []struct {
		Name  string
		Input string

		Valid bool
		Key   string
		Val   string
	}{
		{
			Name:  "normal",
			Input: "hoge=fuga",

			Valid: true,
			Key:   "hoge",
			Val:   "fuga",
		},
		{
			Name:  "trim space",
			Input: " hoge = fuga ",

			Valid: true,
			Key:   "hoge",
			Val:   "fuga",
		},
		{
			Name:  "quoted and needs trimming",
			Input: ` "hoge" = 'fuga' `,

			Valid: true,
			Key:   "hoge",
			Val:   "fuga",
		},
		{
			Name:  "invalid",
			Input: `'ho'ge'='fuga'`,

			Valid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			env := newEnv(tc.Input)
			if !tc.Valid {
				if env.Err() == nil {
					t.Errorf("error should be occurred but nil: %#v", env)
				}
				return
			}
			if env.Err() != nil {
				t.Errorf("error should be nil but: %s", env.Err())
			}
			if tc.Key != env.Key() {
				t.Errorf("invalid key. out=%q, expect=%q", env.Key(), tc.Key)
			}
			if tc.Val != env.Val() {
				t.Errorf("invalid val. out=%q, expect=%q", env.Val(), tc.Val)
			}
		})
	}
}

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
