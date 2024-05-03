package errorr

type errorr interface {
	Error() error
}

type Error struct{}

func (e *Error) Error(er error) error {
	return er
}
