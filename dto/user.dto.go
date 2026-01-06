package dto

type User struct {
	Id         int64    `json:"id"`
	Xid        string   `json:"xid"`
	Phone      string   `json:"phone"`
	FullName   string   `json:"fullName"`
	Email      string   `json:"email"`
	Age        string   `json:"age,omitempty"`
	Avatar     *File    `json:"avatar,omitempty"`
	Status     *Status  `json:"status"`
	ModifiedBy *Subject `json:"modifiedBy"`
	CreatedAt  int64    `json:"createdAt"`
	UpdatedAt  int64    `json:"updatedAt"`
	Version    int64    `json:"version"`
}
