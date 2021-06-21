package command

type NotFound struct {
	Msg   string
	Group []string
}

type BadArguments struct {
	Msg string
}

func (err NotFound) Error() string {
	return err.Msg
}

func (err BadArguments) Error() string {
	return err.Msg
}

var ErrBadArguments error = &BadArguments{}
var ErrNotFound error = &NotFound{}
