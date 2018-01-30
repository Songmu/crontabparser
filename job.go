package crontab

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type Job struct {
	raw     string
	line    int
	hasUser bool

	User     string
	Command  string
	Schedule *Schedule
}

func (jo *Job) Type() Type {
	return TypeJob
}

func (jo *Job) Err() error {
	return nil
}

func (jo *Job) Raw() string {
	return jo.raw
}

func (jo *Job) Line() int {
	return jo.line
}

var definitions = map[string][5]string{
	"@yearly":   [5]string{"0", "0", "1", "1", "*"},
	"@annually": [5]string{"0", "0", "1", "1", "*"},
	"@monthly":  [5]string{"0", "0", "1", "*", "*"},
	"@weekly":   [5]string{"0", "0", "*", "*", "0"},
	"@daily":    [5]string{"0", "0", "*", "*", "*"},
	"@hourly":   [5]string{"0", "*", "*", "*", "*"},
}

func fieldsN(str string, n int) (flds []string) {
	str = strings.TrimSpace(str)
	offset := 0
	buf := &bytes.Buffer{}
	for _, r := range str {
		if n < 2 {
			flds = append(flds, strings.TrimSpace(str[offset:]))
			break
		}
		offset += len(string(r))
		if unicode.IsSpace(r) {
			if buf.Len() > 0 {
				flds = append(flds, buf.String())
				n--
				buf.Reset()
			}
		} else {
			buf.WriteRune(r)
		}
	}
	return flds
}

func (jo *Job) parse() error {
	if strings.HasPrefix(jo.raw, "@") {
		var flds []string
		if jo.hasUser {
			flds = fieldsN(jo.raw, 3)
			if len(flds) != 3 {
				return fmt.Errorf("field: %q is invalid", jo.raw)
			}
			jo.User = flds[1]
			jo.Command = flds[2]
		} else {
			flds = fieldsN(jo.raw, 2)
			if len(flds) != 2 {
				return fmt.Errorf("field: %q is invalid", jo.raw)
			}
			jo.Command = flds[1]
		}
		def, ok := definitions[flds[0]]
		if !ok {
			return fmt.Errorf("invalid definition: %q", flds[0])
		}
		jo.Schedule = NewSchedule(def[0], def[1], def[2], def[3], def[4])
	} else {
		var flds []string
		if jo.hasUser {
			flds = fieldsN(jo.raw, 7)
			if len(flds) != 7 {
				return fmt.Errorf("field: %q is invalid", jo.raw)
			}
			jo.User = flds[5]
			jo.Command = flds[6]
		} else {
			flds = fieldsN(jo.raw, 6)
			if len(flds) != 6 {
				return fmt.Errorf("field: %q is invalid", jo.raw)
			}
			jo.Command = flds[5]
		}
		jo.Schedule = NewSchedule(flds[0], flds[1], flds[2], flds[3], flds[4])
	}
	return nil
}
