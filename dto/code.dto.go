package dto

type Code string

const (
	CodeOK                 Code = "OK"
	CodeBadRequest         Code = "BAD_REQUEST"
	CodeUnauthorized       Code = "UNAUTHORIZED"
	CodeForbidden          Code = "FORBIDDEN"
	CodeNotFound           Code = "NOT_FOUND"
	CodeTooManyRequests    Code = "TOO_MANY_REQUESTS"
	CodeInternalError      Code = "INTERNAL_ERROR"
	CodeServiceUnavailable Code = "SERVICE_UNAVAILABLE"
)
