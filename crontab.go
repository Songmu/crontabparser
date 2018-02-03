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
	entries []Entry
	env     map[string]string
}

func Parse(rdr io.Reader, hasUser bool) (*Crontab, error) {
	ct := &Crontab{env: make(map[string]string)}
	scr := bufio.NewScanner(rdr)
	for scr.Scan() {
		ct.entries = append(ct.entries, ct.parseLine(scr.Text(), hasUser))
	}
	return ct, scr.Err()
}

var jobReg = regexp.MustCompile(`^\s*(?:@|\*|[0-9])`)

func (ct *Crontab) parseLine(line string, hasUser bool) Entry {
	switch {
	case strings.HasPrefix(line, "#"):
		return &Comment{
			raw: line,
		}
	case strings.TrimSpace(line) == "":
		return &Empty{
			raw: line,
		}
	case jobReg.MatchString(line):
		j := &Job{
			raw:     line,
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
			raw: line,
		}
		err := env.parse()
		if err == nil {
			ct.env[env.Key()] = env.Val()
		}
		return env
	default:
		return &Invalid{
			raw: line,
		}
	}
}

func (ct *Crontab) Entries() []Entry {
	return ct.entries
}

func (ct *Crontab) Jobs() (jobs []*Job) {
	for _, ent := range ct.entries {
		if j, ok := ent.(*Job); ok {
			jobs = append(jobs, j)
		}
	}
	return jobs
}

func (ct *Crontab) Valid() bool {
	for _, ent := range ct.entries {
		if ent.Err() != nil {
			return false
		}
	}
	return true
}

func cloneMap(orig map[string]string) map[string]string {
	newMap := make(map[string]string, len(orig))
	for k, v := range orig {
		newMap[k] = v
	}
	return newMap
}
