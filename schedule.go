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
	raw       string
	minute    *ScheduleEntity
	hour      *ScheduleEntity
	day       *ScheduleEntity
	month     *ScheduleEntity
	dayOfWeek *ScheduleEntity
}

func (sche *Schedule) Raw() string {
	return sche.raw
}

func (sche *Schedule) Minute() *ScheduleEntity {
	return sche.minute
}

func (sche *Schedule) Hour() *ScheduleEntity {
	return sche.hour
}

func (sche *Schedule) Day() *ScheduleEntity {
	return sche.day
}

func (sche *Schedule) Month() *ScheduleEntity {
	return sche.month
}

func (sche *Schedule) DayOfWeek() *ScheduleEntity {
	return sche.dayOfWeek
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
	sche.raw = raw
	var ers errors
	sche.minute, err = NewScheduleEntity(min, ScheduleMinute)
	if err != nil {
		ers = append(ers, err)
	}
	sche.hour, err = NewScheduleEntity(hour, ScheduleHour)
	if err != nil {
		ers = append(ers, err)
	}
	sche.day, err = NewScheduleEntity(day, ScheduleDay)
	if err != nil {
		ers = append(ers, err)
	}
	sche.month, err = NewScheduleEntity(month, ScheduleMonth)
	if err != nil {
		ers = append(ers, err)
	}
	sche.dayOfWeek, err = NewScheduleEntity(dayOfWeek, ScheduleDayOfWeek)
	if err != nil {
		ers = append(ers, err)
	}
	return sche, ers.err()
}

func (sche *Schedule) Match(t time.Time) bool {
	if !sche.minute.Match(t.Minute()) || !sche.hour.Match(t.Hour()) || !sche.month.Match(int(t.Month())) {
		return false
	}
	if sche.dayOfWeek.raw != "*" {
		if sche.dayOfWeek.Match(int(t.Weekday())) {
			return true
		}
	}
	return sche.day.Match(t.Day())
}
