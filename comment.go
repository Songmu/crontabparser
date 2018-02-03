package checron

type Comment struct {
	raw string
}

func (co *Comment) Type() Type {
	return TypeComment
}

func (co *Comment) Err() error {
	return nil
}

func (co *Comment) Raw() string {
	return co.raw
}
