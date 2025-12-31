package routek

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/konsultin/project-goes-here/dto"
	"github.com/konsultin/project-goes-here/libs/errk"
	"github.com/valyala/fasthttp"
)

type Responder struct {
	debug bool
}

// NewResponder creates a responder; debug=true will include error details in responses.
func NewResponder(debug bool) *Responder {
	return &Responder{debug: debug}
}

// Success sends a successful dto.Response with the given status, code, message, and payload data.
func (r *Responder) Success(ctx *fasthttp.RequestCtx, status int, code dto.Code, message string, data any) {
	resp := dto.Response[any]{
		Message:   message,
		Code:      code,
		Data:      data,
		Timestamp: time.Now().UTC().UnixMilli(),
	}
	r.write(ctx, status, resp)
}

// Error standardizes error responses, mapping errk metadata to HTTP status and code.
func (r *Responder) Error(ctx *fasthttp.RequestCtx, status int, code dto.Code, message string, err error) {
	var data any

	if err != nil {
		if e := statusFromErrk(err); e != nil {
			code = dto.Code(e.Code())
			if e.Message() != "" {
				message = e.Message()
			}
			if st, ok := httpStatusFromMetadata(e.Metadata()); ok {
				status = st
			}
			if r.debug {
				data = map[string]any{
					"error":  e.Error(),
					"traces": e.Traces(),
				}
			}
		} else if r.debug {
			data = map[string]any{"error": err.Error()}
		}
	}

	if status == 0 {
		status = fasthttp.StatusInternalServerError
	}
	if code == "" {
		code = dto.CodeInternalError
	}
	if message == "" {
		message = "internal server error"
	}

	resp := dto.Response[any]{
		Message:   message,
		Code:      code,
		Data:      data,
		Timestamp: time.Now().UTC().UnixMilli(),
	}
	r.write(ctx, status, resp)
}

// write marshals the payload and writes it to the response, with a resilient fallback when marshaling fails.
func (r *Responder) write(ctx *fasthttp.RequestCtx, status int, payload any) {
	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal response: %v", err)
		fallback := dto.Response[any]{
			Message:   "internal server error",
			Code:      dto.CodeInternalError,
			Data:      nil,
			Timestamp: time.Now().UTC().UnixMilli(),
		}
		fallbackBody, fallbackErr := json.Marshal(fallback)
		if fallbackErr != nil {
			log.Printf("failed to marshal fallback response: %v", fallbackErr)
			ctx.Response.Header.Set("Content-Type", "application/json")
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			ctx.SetBodyString(
				fmt.Sprintf(
					`{"message":"internal server error","code":"INTERNAL_ERROR","data":null,"timestamp":%d}`,
					time.Now().UTC().UnixMilli(),
				),
			)
			return
		}
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody(fallbackBody)
		return
	}

	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.SetStatusCode(status)
	ctx.SetBody(body)
}

// statusFromErrk extracts *errk.Error for additional metadata handling.
func statusFromErrk(err error) *errk.Error {
	var e *errk.Error
	if errors.As(err, &e) {
		return e
	}
	return nil
}

// httpStatusFromMetadata reads http_status from errk metadata when present.
func httpStatusFromMetadata(md map[string]interface{}) (int, bool) {
	if md == nil {
		return 0, false
	}
	if v, ok := md["http_status"]; ok {
		switch t := v.(type) {
		case int:
			return t, true
		case int64:
			return int(t), true
		case float64:
			return int(t), true
		}
	}
	return 0, false
}
