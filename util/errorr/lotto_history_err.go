package errorr

type errorr interface {
	Error(er error) error
}

type Error struct{}

func (e *Error) Error(er error) error {
	return er
}
