package crontabparser

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
			if sche.Raw() != tc.Input {
				t.Errorf("result of raw() is strange. out: %s, expect: %s", sche.Raw(), tc.Input)
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

	errTestCases := []struct {
		Name  string
		Input string
	}{
		{
			Name:  "completely invalid",
			Input: "invalid",
		},
		{
			Name:  "invalid definition",
			Input: "@invalid",
		},
		{
			Name:  "invalid",
			Input: "invalid",
		},
		{
			Name:  "invalid entities",
			Input: "61 25 0 0 SAN",
		},
	}
	for _, tc := range errTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := ParseSchedule(tc.Input)
			if err == nil {
				t.Errorf("error should be occurred but nil")
			}
		})
	}
}

func TestWarnings(t *testing.T) {
	testCases := []struct {
		Name  string
		Input string

		WarningNum int
	}{
		{
			Name:       "no error",
			Input:      "*/1 * * * *",
			WarningNum: 0,
		},
		{
			Name:       "asterisk minutes",
			Input:      "* * * * *",
			WarningNum: 1,
		},
		{
			Name:       "both specified day and day-of-week",
			Input:      "*/1 * 13 * Fri",
			WarningNum: 1,
		},
		{
			Name:       "2 warnings",
			Input:      "* * 13 * Fri",
			WarningNum: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			sche, err := ParseSchedule(tc.Input)
			if err != nil {
				t.Errorf("error should be nil but: %s", err)
			}
			warningNum := len(sche.Warnings())
			if tc.WarningNum != warningNum {
				t.Errorf("output num is missmatched. out=%d, expedted=%d", warningNum, tc.WarningNum)
			}
		})
	}
}
