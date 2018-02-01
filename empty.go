package checron

type Empty struct {
	raw  string
	line int
}

func (em *Empty) Type() Type {
	return TypeEmpty
}

func (em *Empty) Err() error {
	return nil
}

func (em *Empty) Raw() string {
	return em.raw
}

func (em *Empty) Line() int {
	return em.line
}
