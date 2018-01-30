package crontab

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

type Entry interface {
	Type() Type
	Err() error
	Raw() string
	Line() int
}

//go:generate stringer -type=Type -trimprefix Type
type Type int

const (
	TypeError Type = iota
	TypeJob
	TypeComment
	TypeEmpty
	TypeEnv
)

type Crontab struct {
	Entries []Entry
}

func Parse(rdr io.Reader, hasUser bool) (*Crontab, error) {
	ct := &Crontab{}
	lineNum := 0
	scr := bufio.NewScanner(rdr)
	for scr.Scan() {
		lineNum++
		ct.Entries = append(ct.Entries, parseLine(scr.Text(), lineNum, hasUser))
	}
	return ct, scr.Err()
}

var jobReg = regexp.MustCompile(`^(?:@|\*|[0-9])`)

func parseLine(line string, lineNum int, hasUser bool) Entry {
	switch {
	case strings.HasPrefix(line, "#"):
		return &Comment{
			raw:  line,
			line: lineNum,
		}
	case strings.TrimSpace(line) == "":
		return &Empty{
			raw:  line,
			line: lineNum,
		}
	case jobReg.MatchString(line):
		return &Job{
			raw:     line,
			line:    lineNum,
			hasUser: hasUser,
		}
	case strings.Contains(line, "="):
		return &Env{
			raw:  line,
			line: lineNum,
		}
	default:
		return &Error{
			raw:  line,
			line: lineNum,
		}
	}
}

func (ct *Crontab) Jobs() (jobs []*Job) {
	for _, ent := range ct.Entries {
		if j, ok := ent.(*Job); ok {
			jobs = append(jobs, j)
		}
	}
	return jobs
}

func (ct *Crontab) Valid() bool {
	for _, ent := range ct.Entries {
		if ent.Err() != nil {
			return false
		}
	}
	return true
}
