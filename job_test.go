package checron

import (
	"reflect"
	"testing"
)

func TestParseJob(t *testing.T) {
	testCases := []struct {
		Name    string
		Input   string
		HasUser bool

		Valid    bool
		User     string
		Command  string
		Schedule [5]string
	}{
		{
			Name:  "all asterisk",
			Input: " * * * * * perl",

			Valid:    true,
			Command:  "perl",
			Schedule: [5]string{"*", "*", "*", "*", "*"},
		},
		{
			Name:  "normal",
			Input: `*/15 1-11/4 1 1 1 perl -E 'say "Hello"'`,

			Valid:    true,
			Command:  `perl -E 'say "Hello"'`,
			Schedule: [5]string{"*/15", "1-11/4", "1", "1", "1"},
		},
		{
			Name:    "hourly has user",
			Input:   `@hourly songmu perl -E`,
			HasUser: true,

			Valid:    true,
			Command:  `perl -E`,
			User:     "songmu",
			Schedule: definitions["@hourly"],
		},
		{
			Name:  "all asterisk",
			Input: "* * * *  perl",

			Valid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			j := ParseJob(tc.Input, tc.HasUser, nil)
			if tc.Input != j.Raw() {
				t.Errorf("invalid raw input. out=%q, expect=%q", j.Raw(), tc.Input)
			}
			if !tc.Valid {
				if j.Err() == nil {
					t.Errorf("error should be occurred but nil: %#v", j)
				}
				return
			}
			if j.Err() != nil {
				t.Errorf("error should be nil but: %s", j.Err())
			}
			if j.Stdin() != nil {
				t.Errorf("something went wrong")
			}
			if tc.User != j.User() {
				t.Errorf("invalid user. out=%q, expect=%q", j.User(), tc.User)
			}
			if tc.Command != j.Command() {
				t.Errorf("invalid command. out=%q, expect=%q", j.Command(), tc.Command)
			}
			outSche := [5]string{
				j.Schedule().minute.Raw(),
				j.Schedule().hour.Raw(),
				j.Schedule().day.Raw(),
				j.Schedule().month.Raw(),
				j.Schedule().dayOfWeek.Raw(),
			}
			if !reflect.DeepEqual(outSche, tc.Schedule) {
				t.Errorf("invalid schedule.\n   out: %v\nexpect: %v", outSche, tc.Schedule)
			}
		})
	}
}

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

func TestParseCommand(t *testing.T) {
	testCases := []struct {
		Name  string
		Input string

		Command, Stdin string
	}{
		{
			Name:    "normal",
			Input:   `perl -E "say 'Hello'"`,
			Command: `perl -E "say 'Hello'"`,
		},
		{
			Name:    "escaped percent",
			Input:   `/path/to/cmd > /path/to/log.$(date +\%Y\%m\%d) 2>&1`,
			Command: `/path/to/cmd > /path/to/log.$(date +%Y%m%d) 2>&1`,
		},
		{
			Name:    "with backslashes which not for escaping",
			Input:   `/path/to/cmd\ \`,
			Command: `/path/to/cmd\ \`,
		},
		{
			Name:    "It's 10pm",
			Input:   `mail -s "It's 10pm" joe%Joe,%%Where are your kids?%`,
			Command: `mail -s "It's 10pm" joe`,
			Stdin:   "Joe,\n\nWhere are your kids?\n",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			cmd, in := parseCommand(tc.Input)
			if cmd != tc.Command {
				t.Errorf("invalid command: out=%s, expect=%s", cmd, tc.Command)
			}
			if in != tc.Stdin {
				t.Errorf("invalid stdin: out=%s, expect=%s", in, tc.Stdin)
			}
		})
	}
}
