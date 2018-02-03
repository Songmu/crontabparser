package checron

// Empty cron line
type Empty struct {
	raw string
}

// Type returns the type
func (em *Empty) Type() Type {
	return TypeEmpty
}

// Err returns the error or nil
func (em *Empty) Err() error {
	return nil
}

// Raw return raw contents of the line
func (em *Empty) Raw() string {
	return em.raw
}
