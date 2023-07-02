package structs

var InternalErrorMessage = "An internal error occured"
var BadRequestErrorMessage = "Bad request, check your parameters"

type ServerError struct {
	Message string
}

func GetInternalServerError() ServerError {
	return ServerError{Message: InternalErrorMessage}
}

type BadRequestError struct {
	Message string
}

func GetBadRequestError() BadRequestError {
	return BadRequestError{Message: BadRequestErrorMessage}
}
