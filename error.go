package checron

import "fmt"

type Error struct {
	raw  string
	line int
}

func (er *Error) Type() Type {
	return TypeError
}

func (er *Error) Err() error {
	return fmt.Errorf("invalid line: %q", er.raw)
}

func (er *Error) Raw() string {
	return er.raw
}

func (er *Error) Line() int {
	return er.line
}
