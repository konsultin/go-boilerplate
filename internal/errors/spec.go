package errors

import (
	"github.com/Konsultin/project-goes-here/internal/svc-core/constant"
	"github.com/Konsultin/project-goes-here/libs/errk"
	fhttp "github.com/valyala/fasthttp"
)

var b = errk.NewBuilder(constant.ServiceName)

var ResourceNotFound = b.NewError("E_COMM_1", "Resource not found",
	errk.WithHTTPStatus(fhttp.StatusNotFound),
)
