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

type File struct {
	FileName  string `json:"fileName"`
	Url       string `json:"url"`
	Signature string `json:"signature,omitempty"`
}

type Status struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
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

type ControlStatus_Result struct {
	Id   ControlStatus_Enum `json:"id,omitempty"`
	Name string             `json:"name,omitempty"`
}
