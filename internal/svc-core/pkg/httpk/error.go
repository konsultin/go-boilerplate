package httpk

import (
	"github.com/konsultin/errk"
	fhttp "github.com/valyala/fasthttp"
)

const (
	HttpStatusMetadata      = "httpStatus"
	OverrideMessageMetadata = "message"
	ErrorMetadata           = "errorMetadata"
)

func withStatus(status uint32) errk.SetOptionFn {
	return errk.AddMetadata(HttpStatusMetadata, status)
}

func OverrideMessage(message string) errk.SetOptionFn {
	return errk.AddMetadata(OverrideMessageMetadata, message)
}

var b = errk.NewBuilder(pkgNamespace)

var InternalError = b.NewError("500", "Internal Error",
	errk.WithHTTPStatus(fhttp.StatusInternalServerError))

var BadRequestError = b.NewError("400", "Bad Request",
	errk.WithHTTPStatus(fhttp.StatusBadRequest),
)

var UnauthorizedError = b.NewError("401", "Unauthorized",
	errk.WithHTTPStatus(fhttp.StatusUnauthorized),
)

var ForbiddenError = b.NewError("403", "Forbidden",
	errk.WithHTTPStatus(fhttp.StatusForbidden),
)

var NotFoundError = b.NewError("404", "Not Found",
	errk.WithHTTPStatus(fhttp.StatusNotFound),
)

var CancelError = b.NewError("408", "Request Canceled",
	errk.WithHTTPStatus(fhttp.StatusRequestTimeout),
)

var InvalidPayloadError = b.NewError("422", "Invalid Payload",
	errk.WithHTTPStatus(fhttp.StatusUnprocessableEntity),
)
