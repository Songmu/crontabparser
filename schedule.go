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

type entityParam struct {
	Range   [2]int
	Aliases []string
}

var entityParams = map[ScheduleType]entityParam{
	ScheduleMinute: {
		Range: [2]int{0, 59},
	},
	ScheduleHour: {
		Range: [2]int{0, 23},
	},
	ScheduleDay: {
		Range: [2]int{1, 31},
	},
	ScheduleMonth: {
		Range:   [2]int{1, 12},
		Aliases: []string{"jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec"},
	},
	ScheduleDayOfWeek: {
		Range:   [2]int{0, 7},
		Aliases: []string{"sun", "mon", "tue", "wed", "thu", "fri", "sat"},
	},
}

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
