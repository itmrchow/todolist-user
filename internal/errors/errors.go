package errors

// 500
type InternalServerError struct {
}

type MsgInternalServerError string

const (
	MsgInternalError MsgInternalServerError = "internal server error"
)

func (e *InternalServerError) Error() string {
	return string(MsgInternalError)
}

// 400
type BadRequestError struct {
	Msg MsgBadRequestError
}

type MsgBadRequestError string

const (
	MsgEmailAlreadyExists MsgBadRequestError = "email already exists"
)

func (e *BadRequestError) Error() string {
	return string(e.Msg)
}
