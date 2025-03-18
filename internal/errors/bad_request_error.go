package errors

// 400
type BadRequestError struct {
	Msg string
}

func (e *BadRequestError) Error() string {
	return e.Msg
}

var (
	Err400EmailAlreadyExists = BadRequestError{Msg: "email already exists"}
	Err400InvalidLoginInfo   = BadRequestError{Msg: "invalid login info"}
)
