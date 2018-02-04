package crontabparser

import (
	"reflect"
	"testing"
)

func TestNewScheduleEntity(t *testing.T) {
	testCases := []struct {
		Name  string
		Input string
		Type  scheduleType

		Expanded []int
	}{
		{
			Name:  "minutes: asterisk and increments",
			Input: "*/15",
			Type:  scheduleMinute,

			Expanded: []int{0, 15, 30, 45},
		},
		{
			Name:  "hour: multi range",
			Input: "2,3-7/2,5",
			Type:  scheduleHour,

			Expanded: []int{2, 3, 5, 7},
		},
		{
			Name:  "day: first day of month",
			Input: "1",
			Type:  scheduleDay,

			Expanded: []int{1},
		},
		{
			Name:  "month: asterisk",
			Input: "*",
			Type:  scheduleMonth,

			Expanded: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		},
		{
			Name:  "month: aliase range with increments",
			Input: "FEB-Sep/3,10-12",
			Type:  scheduleMonth,

			Expanded: []int{2, 5, 8, 10, 11, 12},
		},
		{
			Name:  "day_of_week: multi",
			Input: "3,*/2",
			Type:  scheduleDayOfWeek,

			Expanded: []int{0, 2, 3, 4, 6},
		},
		{
			Name:  "day_of_week: aliases",
			Input: "mon-Tue,4",
			Type:  scheduleDayOfWeek,

			Expanded: []int{1, 2, 4},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			se, err := newScheduleEntity(tc.Input, tc.Type)
			if err != nil {
				t.Errorf("error should be nil but: %s", err)
			}
			if !reflect.DeepEqual(se.Expanded(), tc.Expanded) {
				t.Errorf("invalid expanded values.\n   out: %v\nexpect: %v", se.Expanded(), tc.Expanded)
			}
		})
	}

	invalidTestCases := []struct {
		Name  string
		Input string
		Type  scheduleType
	}{
		{
			Name:  "invalid month",
			Input: "0",
			Type:  scheduleMonth,
		},
		{
			Name:  "invalid range of days",
			Input: "10-32/3",
			Type:  scheduleDay,
		},
		{
			Name:  "day_of_week: invalid aliase",
			Input: "mon-Tuo",
			Type:  scheduleDayOfWeek,
		},
		{
			Name:  "hour: invalid step",
			Input: "1-10/N",
			Type:  scheduleHour,
		},
	}
	for _, tc := range invalidTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := newScheduleEntity(tc.Input, tc.Type)
			if err == nil {
				t.Errorf("error should be occurred but nil")
			}
		})
	}
}
