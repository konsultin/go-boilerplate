package middleware

import (
	"strings"
	"sync"
	"time"

	"github.com/konsultin/logk"
	"github.com/valyala/fasthttp"
)

type RateLimiter struct {
	rps     float64
	burst   float64
	ttl     time.Duration
	buckets map[string]*tokenBucket
	mu      sync.Mutex
}

type tokenBucket struct {
	tokens float64
	last   time.Time
}

func NewRateLimiter(rps, burst int) *RateLimiter {
	if rps <= 0 {
		rps = 1
	}
	if burst <= 0 {
		burst = rps
	}

	return &RateLimiter{
		rps:     float64(rps),
		burst:   float64(burst),
		ttl:     5 * time.Minute,
		buckets: make(map[string]*tokenBucket),
	}
}

func (r *RateLimiter) Allow(key string) bool {
	if key == "" {
		key = "global"
	}

	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	bucket, ok := r.buckets[key]
	if !ok || now.Sub(bucket.last) > r.ttl {
		r.buckets[key] = &tokenBucket{tokens: r.burst - 1, last: now}
		return true
	}

	elapsed := now.Sub(bucket.last).Seconds()
	bucket.tokens = minFloat(r.burst, bucket.tokens+(elapsed*r.rps))
	bucket.last = now

	if bucket.tokens < 1 {
		return false
	}

	bucket.tokens--
	return true
}

func RateLimit(rl *RateLimiter, log logk.Logger, onError ErrorResponder) func(fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			if !rl.Allow(clientKey(ctx)) {
				log.Warn("request rejected: rate limit exceeded")
				onError(ctx, fasthttp.StatusTooManyRequests, "TOO_MANY_REQUESTS", "too many requests", nil)
				return
			}
			next(ctx)
		}
	}
}

func clientKey(ctx *fasthttp.RequestCtx) string {
	if forwardedFor := strings.TrimSpace(string(ctx.Request.Header.Peek("X-Forwarded-For"))); forwardedFor != "" {
		parts := strings.Split(forwardedFor, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	if ip := ctx.RemoteIP(); ip != nil {
		return ip.String()
	}

	return "unknown"
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
