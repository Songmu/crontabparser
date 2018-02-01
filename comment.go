package checron

type Comment struct {
	raw  string
	line int
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

func (co *Comment) Line() int {
	return co.line
}
