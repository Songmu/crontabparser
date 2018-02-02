package checron

import (
	"fmt"
	"strings"
	"time"
)

var definitions = map[string][5]string{
	"@yearly":   [5]string{"0", "0", "1", "1", "*"},
	"@annually": [5]string{"0", "0", "1", "1", "*"},
	"@monthly":  [5]string{"0", "0", "1", "*", "*"},
	"@weekly":   [5]string{"0", "0", "*", "*", "0"},
	"@daily":    [5]string{"0", "0", "*", "*", "*"},
	"@hourly":   [5]string{"0", "*", "*", "*", "*"},
	"@reboot":   [5]string{"0", "0", "0", "0", "0"}, // XXX
}

type Schedule struct {
	Raw       string
	Minute    *ScheduleEntity
	Hour      *ScheduleEntity
	Day       *ScheduleEntity
	Month     *ScheduleEntity
	DayOfWeek *ScheduleEntity
}

func ParseSchedule(raw string) (sche *Schedule, err error) {
	if strings.HasPrefix(raw, "@") {
		def, ok := definitions[raw]
		if !ok {
			return nil, fmt.Errorf("invalid schedule definition: %q", raw)
		}
		if raw == "@reboot" {
			return &Schedule{}, nil
		}
		return newSchedule(raw, def[0], def[1], def[2], def[3], def[4])
	}
	flds := strings.Fields(raw)
	if len(flds) != 5 {
		return nil, fmt.Errorf("invalid schedule: %q", raw)
	}
	return newSchedule(raw, flds[0], flds[1], flds[2], flds[3], flds[4])
}

func newSchedule(raw string, min, hour, day, month, dayOfWeek string) (sche *Schedule, err error) {
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
