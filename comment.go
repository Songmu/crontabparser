package crontabparser

// Comment cron line
type Comment struct {
	raw string
}

// Type TypeCommend
func (co *Comment) Type() Type {
	return TypeComment
}

// Err always return nil
func (co *Comment) Err() error {
	return nil
}

// Raw content of comment line
func (co *Comment) Raw() string {
	return co.raw
}
