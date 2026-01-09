package errors

import (
	"github.com/konsultin/project-goes-here/internal/svc-core/constant"
	"github.com/konsultin/errk"
	fhttp "github.com/valyala/fasthttp"
)

var b = errk.NewBuilder(constant.ServiceName)

var ResourceNotFound = b.NewError("E_COMM_1", "Resource not found",
	errk.WithHTTPStatus(fhttp.StatusNotFound),
)
