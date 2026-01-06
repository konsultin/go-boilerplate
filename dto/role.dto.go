package dto

type Role struct {
	Xid         string                `json:"xid,omitempty"`
	Version     int64                 `json:"version,omitempty"`
	Name        string                `json:"name,omitempty"`
	RoleType    *RoleType_Result      `json:"roleType,omitempty"`
	CreatedAt   int64                 `json:"createdAt,omitempty"`
	UpdatedAt   int64                 `json:"updatedAt,omitempty"`
	ModifiedBy  *Subject              `json:"modifiedBy,omitempty"`
	Description string                `json:"description,omitempty"`
	Privileges  []string              `json:"privileges,omitempty"`
	Status      *ControlStatus_Result `json:"status,omitempty"`
}

type RoleType_Result struct {
	Id   RoleType_Enum `json:"id,omitempty"`
	Name string        `json:"name,omitempty"`
}
