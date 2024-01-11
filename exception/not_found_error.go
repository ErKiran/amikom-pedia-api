package exception

type NotFoundError struct {
	Error string
}
type BadRequestError struct {
	Error string
}

func NewNotFoundError(error string) NotFoundError {
	return NotFoundError{Error: error}
}

func NewBadRequestError(error string) BadRequestError {
	return BadRequestError{Error: error}
}
