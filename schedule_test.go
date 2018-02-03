package checron

import (
	"log"
	"testing"
	"time"
)

func date(str string) time.Time {
	const layout = "2006-01-02 15:04"
	t, err := time.ParseInLocation(layout, str, time.Local)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func TestParseSchedule(t *testing.T) {
	testCases := []struct {
		Name  string
		Input string

		Mathes    []time.Time
		UnMatches []time.Time
	}{
		{
			Name:  "normal",
			Input: "0 1 1 1 *",

			Mathes: []time.Time{
				date("2018-01-01 01:00"),
			},
			UnMatches: []time.Time{
				date("2018-01-01 01:01"),
			},
		},
		{
			Name:  "specify day of week",
			Input: " 11 13-15/2 * 2 Mon",

			Mathes: []time.Time{
				date("2018-02-05 13:11"),
				date("2018-02-12 15:11"),
				date("2018-02-19 13:11"),
				date("2018-02-26 15:11"),
			},
			UnMatches: []time.Time{
				date("2018-01-01 01:01"),
			},
		},
		{
			Name:  "yearly",
			Input: "@yearly",

			Mathes: []time.Time{
				date("2000-01-01 00:00"),
				date("2015-01-01 00:00"),
				date("2016-01-01 00:00"),
				date("2018-01-01 00:00"),
			},
			UnMatches: []time.Time{
				date("2018-01-01 01:01"),
			},
		},
		{
			Name:  "reboot",
			Input: "@reboot",

			UnMatches: []time.Time{
				date("2018-01-01 01:01"),
				date("2000-01-01 00:00"),
				date("2018-02-12 15:11"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			sche, err := ParseSchedule(tc.Input)
			if err != nil {
				t.Errorf("error should be nil but: %s", err)
			}
			for _, ti := range tc.Mathes {
				if !sche.Match(ti) {
					t.Errorf("schedule(%v) does not match %v", sche, ti)
				}
			}
			for _, ti := range tc.UnMatches {
				if sche.Match(ti) {
					t.Errorf("schedule(%v) unintentionally matches %v", sche, ti)
				}
			}
		})
	}
}
