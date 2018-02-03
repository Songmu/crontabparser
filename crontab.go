package checron

import (
	"bufio"
	"fmt"
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
	raw     string
	entries []Entry
	env     map[string]string

	errors errors
}

func Parse(rdr io.Reader, hasUser bool) (*Crontab, error) {
	ct := &Crontab{env: make(map[string]string)}
	scr := bufio.NewScanner(rdr)
	for scr.Scan() {
		ct.entries = append(ct.entries, ct.parseLine(scr.Text(), hasUser))
	}
	return ct, scr.Err()
}

func (ct *Crontab) Raw() string {
	return ct.raw
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

func (ct *Crontab) Err() error {
	if ct.errors == nil {
		ct.errors = []error{}
		for i, ent := range ct.entries {
			if ent.Err() != nil {
				ct.errors = append(ct.errors, fmt.Errorf("line %d: %s", i+1, ent.Err().Error()))
			}
		}
	}
	return ct.errors.err()
}

var jobReg = regexp.MustCompile(`^\s*(?:@|\*|[0-9])`)

func (ct *Crontab) parseLine(line string, hasUser bool) Entry {
	switch {
	case strings.HasPrefix(strings.TrimSpace(line), "#"):
		return &Comment{raw: line}
	case strings.TrimSpace(line) == "":
		return &Empty{raw: line}
	case jobReg.MatchString(line):
		return NewJob(line, hasUser, cloneMap(ct.env))
	case strings.Contains(line, "="):
		env := newEnv(line)
		if env.Err() == nil {
			ct.env[env.Key()] = env.Val()
		}
		return env
	default:
		return &Invalid{raw: line}
	}
}

func cloneMap(orig map[string]string) map[string]string {
	newMap := make(map[string]string, len(orig))
	for k, v := range orig {
		newMap[k] = v
	}
	return newMap
}
