package checron

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type errors []error

func (ers errors) Error() string {
	strs := make([]string, len(ers))
	for i, err := range ers {
		strs[i] = err.Error()
	}
	return strings.Join(strs, "\n")
}

func (ers errors) err() error {
	if len(ers) == 0 {
		return nil
	}
	return ers
}

type Job struct {
	raw     string
	hasUser bool
	env     map[string]string
	errors  errors

	user     string
	command  string
	schedule *Schedule
}

func (jo *Job) User() string {
	return jo.user
}

func (jo *Job) Command() string {
	return jo.command
}

func (jo *Job) Schedule() *Schedule {
	return jo.schedule
}

func (jo *Job) Type() Type {
	return TypeJob
}

func (jo *Job) Err() error {
	return jo.errors.err()
}

func (jo *Job) Raw() string {
	return jo.raw
}

func (jo *Job) Env() map[string]string {
	return jo.env
}

func (jo *Job) setError(err error) {
	if err == nil {
		return
	}
	jo.errors = append(jo.errors, err)
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
	if buf.Len() > 0 {
		flds = append(flds, buf.String())
	}
	return flds
}

var scheduleReg = regexp.MustCompile(`^(@\w+|(?:\S+\s+){5})(.*)$`)

func (jo *Job) parse() (err error) {
	if m := scheduleReg.FindStringSubmatch(strings.TrimSpace(jo.raw)); len(m) == 3 {
		jo.schedule, err = ParseSchedule(strings.TrimSpace(m[1]))
		if err != nil {
			return err
		}
		if jo.hasUser {
			flds := fieldsN(m[2], 2)
			if len(flds) != 2 {
				return fmt.Errorf("field: %q is invalid", jo.raw)
			}
			jo.user = flds[1]
			jo.command = flds[2]
			return nil
		}
		jo.command = m[2]
		return nil
	}
	return fmt.Errorf("field: %q is invalid", jo.raw)
}
