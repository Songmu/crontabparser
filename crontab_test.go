package checron

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		Name    string
		Input   string
		HasUser bool
		Expects []Type
		Valid   bool
	}{
		{
			Name: "normal",
			Input: `# comment
HOGE=FUGA
* * * * * perl
  
@daily perl
`,
			Expects: []Type{
				TypeComment,
				TypeEnv,
				TypeJob,
				TypeEmpty,
				TypeJob,
			},
			Valid: true,
		},
		{
			Name: "has user",
			Input: ` # comment
 HOGE=FUGA
 * * * * * songmu perl
  
 @daily songmu perl
`,
			HasUser: true,
			Expects: []Type{
				TypeComment,
				TypeEnv,
				TypeJob,
				TypeEmpty,
				TypeJob,
			},
			Valid: true,
		},
		{
			Name: "invlid schedule",
			Input: `HOGE=FUGA
* * * *R * perl
`,
			Expects: []Type{
				TypeEnv,
				TypeJob,
			},
			Valid: false,
		},
		{
			Name: "invlid schedule definition",
			Input: `HOGE=FUGA
@hoge perl
`,
			Expects: []Type{
				TypeEnv,
				TypeJob,
			},
			Valid: false,
		},
		{
			Name: "invalid with having user",
			Input: `# comment
HOGE=FUGA
* * * * * perl
`,
			HasUser: true,
			Expects: []Type{
				TypeComment,
				TypeEnv,
				TypeJob,
			},
			Valid: false,
		},
		{
			Name: "invalid env",
			Input: ` 'HOGE=FUGA
* * * * * perl
`,
			Expects: []Type{
				TypeEnv,
				TypeJob,
			},
			Valid: false,
		},
		{
			Name:    "invalid line",
			Input:   "invalid invalid\n",
			Expects: []Type{TypeInvalid},
			Valid:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ct, err := Parse(strings.NewReader(tc.Input), tc.HasUser)
			if err != nil {
				t.Fatalf("failed to parse: %s", err)
			}
			if ct.Raw() != tc.Input {
				t.Errorf("something went wrong: %s", ct.Raw())
			}
			if (ct.Err() != nil) == tc.Valid {
				t.Errorf("something went wrong: %#v, %s", ct, ct.Err())
			}
			for i, ent := range ct.Entries() {
				expect := tc.Expects[i]
				if ent.Type() != expect {
					t.Errorf("entry %#v has unexpected type: %s, expect: %s", ent, ent.Type(), expect)
				}
			}
		})
	}
}

func TestParse_env(t *testing.T) {
	input := `HOGE=111
* * * * * perl
"HOGE"='222'
@daily perl
`
	ct, err := Parse(strings.NewReader(input), false)
	if err != nil {
		t.Fatalf("failed to parse: %s", err)
	}
	jobs := ct.Jobs()
	if jobs[0].Env()["HOGE"] != "111" {
		t.Errorf("job is wrong: %#v", jobs[0])
	}
	if jobs[1].Env()["HOGE"] != "222" {
		t.Errorf("job is wrong: %#v", jobs[1])
	}
}
