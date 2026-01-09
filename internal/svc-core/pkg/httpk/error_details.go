package httpk

import (
	"errors"
	"fmt"

	"github.com/konsultin/errk"
	fhttp "github.com/valyala/fasthttp"
)

type ErrorDetails struct {
	HttpStatus uint32
	Code       string
	Message    string
	Source     *SourceErrorDetails
}

type SourceErrorDetails struct {
	Message  string
	Traces   []string
	Metadata map[string]string
}

func GetErrorDetails(err error, withSource bool) *ErrorDetails {
	var hErr *errk.Error

	ok := errors.As(err, &hErr)
	if !ok {
		hErr = errk.InternalError().Wrap(err)
	}

	errMeta := hErr.Metadata()
	httpStatus, ok := errMeta[HttpStatusMetadata].(uint32)
	if !ok {
		httpStatus = fhttp.StatusInternalServerError
	}

	details := &ErrorDetails{
		HttpStatus: httpStatus,
		Code:       hErr.Code(),
		Message:    hErr.Message(),
	}

	if withSource {
		sourceMsg := ""

		if sourceErr := errors.Unwrap(hErr); sourceErr != nil {
			sourceMsg = sourceErr.Error()
		} else {
			sourceMsg = hErr.Message()
		}

		rawMetadata, _ := errMeta[ErrorMetadata].(map[string]interface{})
		var metadata map[string]string
		if rawMetadata != nil {
			metadata = make(map[string]string)
			for k, v := range rawMetadata {
				metadata[k] = fmt.Sprintf("%v", v)
			}
		}
		details.Source = &SourceErrorDetails{
			Message:  sourceMsg,
			Traces:   hErr.Traces(),
			Metadata: metadata,
		}
	}

	return details

}
