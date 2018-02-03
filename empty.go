package checron

type Empty struct {
	raw string
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
