package errors

// 500
type InternalServerError struct {
	Msg string
}

func (e *InternalServerError) Error() string {
	return string(e.Msg)
}

var (
	Err500InternalServer = InternalServerError{Msg: "internal server error"}
)
