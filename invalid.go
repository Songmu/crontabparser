package checron

import "fmt"

// Invalid cron entry
type Invalid struct {
	raw string
}

// Type TypeInvalid
func (er *Invalid) Type() Type {
	return TypeInvalid
}

// Err returns error if error
func (er *Invalid) Err() error {
	return fmt.Errorf("invalid line: %q", er.raw)
}

// Raw content of line
func (er *Invalid) Raw() string {
	return er.raw
}
