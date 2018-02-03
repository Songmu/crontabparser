package checron

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"
)

// Job entry in crontab
type Job struct {
	raw string
	env map[string]string

	user     string
	command  string
	stdin    string
	schedule *Schedule

	err error
}

// ParseJob parses job line and returns the *Job
func ParseJob(raw string, hasUser bool, env map[string]string) *Job {
	jo := &Job{
		raw: raw,
		env: env,
	}
	jo.err = jo.parse(hasUser)
	return jo
}

// User of the job
func (jo *Job) User() string {
	return jo.user
}

// Command of the job
func (jo *Job) Command() string {
	return jo.command
}

// Stdin returns input
func (jo *Job) Stdin() io.Reader {
	if jo.stdin == "" {
		return nil
	}
	return strings.NewReader(jo.stdin)
}

// Schedule of the job
func (jo *Job) Schedule() *Schedule {
	return jo.schedule
}

// Type of the job
func (jo *Job) Type() Type {
	return TypeJob
}

// Err of the job
func (jo *Job) Err() error {
	return jo.err
}

// Raw content of the job
func (jo *Job) Raw() string {
	return jo.raw
}

// Env of the job
func (jo *Job) Env() map[string]string {
	return jo.env
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

func (jo *Job) parse(hasUser bool) (err error) {
	if m := scheduleReg.FindStringSubmatch(strings.TrimSpace(jo.raw)); len(m) == 3 {
		jo.schedule, err = ParseSchedule(strings.TrimSpace(m[1]))
		if err != nil {
			return err
		}
		if hasUser {
			flds := fieldsN(m[2], 2)
			if len(flds) != 2 {
				return fmt.Errorf("field: %q is invalid", jo.raw)
			}
			jo.user = flds[0]
			jo.command, jo.stdin = parseCommand(flds[1])
			return nil
		}
		jo.command, jo.stdin = parseCommand(m[2])
		return nil
	}
	return fmt.Errorf("field: %q is invalid", jo.raw)
}

func parseCommand(rawCmd string) (cmd, stdin string) {
	cmdBuf := &bytes.Buffer{}
	stdinBuf := &bytes.Buffer{}
	var endedCmd, escaped bool
	out := func() *bytes.Buffer {
		if endedCmd {
			return stdinBuf
		}
		return cmdBuf
	}
	for _, chr := range rawCmd {
		if escaped {
			if chr != '%' {
				out().WriteRune('\\')
			}
			out().WriteRune(chr)
			escaped = false
			continue
		}

		if chr == '\\' {
			escaped = true
			continue
		}

		if chr == '%' {
			if !endedCmd {
				endedCmd = true
				continue
			}
			out().WriteString("\n")
			continue
		}
		out().WriteRune(chr)
	}
	if escaped {
		out().WriteRune('\\')
	}
	return cmdBuf.String(), stdinBuf.String()
}
