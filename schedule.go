package crontab

type Schedule struct {
	Raw       string
	Minute    *ScheduleEntity
	Hour      *ScheduleEntity
	Day       *ScheduleEntity
	Month     *ScheduleEntity
	DayOfWeek *ScheduleEntity
}

type ScheduleEntity struct {
	Raw  string
	Type ScheduleType
}

//go:generate stringer -type=ScheduleType -trimprefix Schedule
type ScheduleType int

const (
	ScheduleMinute ScheduleType = iota
	ScheduleHour
	ScheduleDay
	ScheduleMonth
	ScheduleDayOfWeek
)

func NewSchedule(raw string, min, hour, day, month, dayOfWeek string) (sche *Schedule, err error) {
	sche.Raw = raw
	sche.Minute = &ScheduleEntity{
		Raw:  min,
		Type: ScheduleMinute,
	}
	sche.Hour = &ScheduleEntity{
		Raw:  hour,
		Type: ScheduleHour,
	}
	sche.Day = &ScheduleEntity{
		Raw:  day,
		Type: ScheduleDay,
	}
	sche.Month = &ScheduleEntity{
		Raw:  month,
		Type: ScheduleMonth,
	}
	sche.DayOfWeek = &ScheduleEntity{
		Raw:  dayOfWeek,
		Type: ScheduleDayOfWeek,
	}
	return sche, nil
}
