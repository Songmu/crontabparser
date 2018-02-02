package checron

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
	TypeInvalid Type = iota
	TypeJob
	TypeComment
	TypeEmpty
	TypeEnv
)

type Crontab struct {
	Entries []Entry
	env     map[string]string
}

func Parse(rdr io.Reader, hasUser bool) (*Crontab, error) {
	ct := &Crontab{env: make(map[string]string)}
	lineNum := 0
	scr := bufio.NewScanner(rdr)
	for scr.Scan() {
		lineNum++
		ct.Entries = append(ct.Entries, ct.parseLine(scr.Text(), lineNum, hasUser))
	}
	return ct, scr.Err()
}

var jobReg = regexp.MustCompile(`^\s*(?:@|\*|[0-9])`)

func (ct *Crontab) parseLine(line string, lineNum int, hasUser bool) Entry {
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
		j := &Job{
			raw:     line,
			line:    lineNum,
			hasUser: hasUser,
			env:     cloneMap(ct.env),
		}
		err := j.parse()
		if err != nil {
			j.setError(err)
		}
		return j
	case strings.Contains(line, "="):
		env := &Env{
			raw:  line,
			line: lineNum,
		}
		env.parse() // error handling
		ct.env[env.Key()] = env.Val()
		return env
	default:
		return &Invalid{
			raw:  line,
			line: lineNum,
		}
	}
}

func cloneMap(orig map[string]string) map[string]string {
	newMap := make(map[string]string, len(orig))
	for k, v := range orig {
		newMap[k] = v
	}
	return newMap
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
