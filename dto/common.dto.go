package dto

type Response[T any] struct {
	Message   string `json:"message"`
	Code      Code   `json:"code"`
	Data      T      `json:"data"`
	Timestamp int64  `json:"timestamp"`
}

type ErrorTrace struct {
	Source string   `json:"source"`
	Trace  []string `json:"trace"`
	Err    error    `json:"err"`
}

type Subject struct {
	Id       string `json:"id"`
	FullName string `json:"fullName"`
	Role     string `json:"role"`
}

type SimpleData struct {
	Value string `json:"value"`
	Name  string `json:"name"`
}

type Empty struct{}