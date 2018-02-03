package checron

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// Entry in crontab
type Entry interface {
	Type() Type
	Err() error
	Raw() string
}

// Type of cron entry
//go:generate stringer -type=Type -trimprefix Type
type Type int

// Types constraints
const (
	TypeInvalid Type = iota
	TypeJob
	TypeComment
	TypeEmpty
	TypeEnv
)

// Crontab represents crontab
type Crontab struct {
	raw     string
	entries []Entry
	env     map[string]string

	errors errors
}

// Parse crontab
func Parse(rdr io.Reader, hasUser bool) (*Crontab, error) {
	buf := &bytes.Buffer{}
	ct := &Crontab{env: make(map[string]string)}
	scr := bufio.NewScanner(rdr)
	for scr.Scan() {
		line := scr.Text()
		buf.WriteString(line + "\n")
		ct.entries = append(ct.entries, ct.parseLine(line, hasUser))
	}
	ct.raw = buf.String()
	return ct, scr.Err()
}

// Raw content of crontab
func (ct *Crontab) Raw() string {
	return ct.raw
}

// Entries returns cron entries
func (ct *Crontab) Entries() []Entry {
	return ct.entries
}

// Jobs returns jobs in crontab
func (ct *Crontab) Jobs() (jobs []*Job) {
	for _, ent := range ct.entries {
		if j, ok := ent.(*Job); ok {
			jobs = append(jobs, j)
		}
	}
	return jobs
}

// Err returns error in crontab if error exists
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
		return ParseJob(line, hasUser, cloneMap(ct.env))
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
