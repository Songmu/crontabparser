package checron

import (
	"time"
)

type Schedule struct {
	Raw       string
	Minute    *ScheduleEntity
	Hour      *ScheduleEntity
	Day       *ScheduleEntity
	Month     *ScheduleEntity
	DayOfWeek *ScheduleEntity
}

func NewSchedule(raw string, min, hour, day, month, dayOfWeek string) (sche *Schedule, err error) {
	sche.Raw = raw
	var ers errors
	sche.Minute, err = NewScheduleEntity(min, ScheduleMinute)
	if err != nil {
		ers = append(ers, err)
	}
	sche.Hour, err = NewScheduleEntity(hour, ScheduleHour)
	if err != nil {
		ers = append(ers, err)
	}
	sche.Day, err = NewScheduleEntity(day, ScheduleDay)
	if err != nil {
		ers = append(ers, err)
	}
	sche.Month, err = NewScheduleEntity(month, ScheduleMonth)
	if err != nil {
		ers = append(ers, err)
	}
	sche.DayOfWeek, err = NewScheduleEntity(dayOfWeek, ScheduleDayOfWeek)
	if err != nil {
		ers = append(ers, err)
	}
	return sche, ers.err()
}

func (sche *Schedule) Match(t time.Time) bool {
	if !sche.Minute.Match(t.Minute()) || !sche.Hour.Match(t.Hour()) || !sche.Month.Match(int(t.Month())) {
		return false
	}
	if sche.DayOfWeek.Raw != "*" {
		if sche.DayOfWeek.Match(int(t.Weekday())) {
			return true
		}
	}
	return sche.Day.Match(t.Day())
}
