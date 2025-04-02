package errors

// 401
type UnauthorizedError struct {
	Msg string
}

func (e *UnauthorizedError) Error() string {
	return e.Msg
}

var (
	Err401Unauthorized = UnauthorizedError{Msg: "invalid token"}
	Err401TokenExpired = UnauthorizedError{Msg: "token expired"}
)
