package crontabparser

import "strings"

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
