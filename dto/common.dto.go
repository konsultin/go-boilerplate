package dto

import "time"

type Response[T any] struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Data      T         `json:"data"`
	Timestamp time.Time `json:"timestamp"`
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
