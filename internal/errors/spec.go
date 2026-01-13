package errors

import (
	"github.com/go-konsultin/errk"
	"github.com/konsultin/project-goes-here/internal/svc-core/constant"
	fhttp "github.com/valyala/fasthttp"
)

var b = errk.NewBuilder(constant.ServiceName)

var ResourceNotFound = b.NewError("E_COMM_1", "Resource not found",
	errk.WithHTTPStatus(fhttp.StatusNotFound),
)

// Session Errors
var SessionGenerationFailed = b.NewError("E_SESS_1", "Failed to generate session",
	errk.WithHTTPStatus(fhttp.StatusInternalServerError),
)

var SessionTokenInvalid = b.NewError("E_SESS_2", "Invalid session token",
	errk.WithHTTPStatus(fhttp.StatusUnauthorized),
)

var SessionExpired = b.NewError("E_SESS_3", "Session has expired",
	errk.WithHTTPStatus(fhttp.StatusUnauthorized),
)

// Client Auth Errors
var InvalidCredentials = b.NewError("E_AUTH_1", "Invalid credentials",
	errk.WithHTTPStatus(fhttp.StatusUnauthorized),
)

var InvalidClientType = b.NewError("E_AUTH_2", "Invalid client type",
	errk.WithHTTPStatus(fhttp.StatusUnauthorized),
)
