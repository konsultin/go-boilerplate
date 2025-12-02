package svcCore

import (
	"encoding/json"

	f "github.com/valyala/fasthttp"
)

// s.response writes a JSON response with a consistent header and handles marshal errors gracefully.
func (s *Server) response(ctx *f.RequestCtx, statusCode int, payload any) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.SetStatusCode(statusCode)

	body, err := json.Marshal(payload)
	if err != nil {
		if s.log != nil {
			s.log.Errorf("failed to marshal response: %v", err)
		}
		ctx.SetStatusCode(f.StatusInternalServerError)
		ctx.SetBodyString(`{"message":"internal server error","code":"INTERNAL_ERROR","data":null,"timestamp":0}`)
		return
	}

	ctx.SetBody(body)
}
