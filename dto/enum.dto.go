package dto

// ===== Role Enums =====

type Role_Enum int32

const (
	Role_UNKNOWN_ROLE Role_Enum = 0

	// Admin
	Role_ANONYMOUS_ADMIN Role_Enum = 1
	Role_ADMIN           Role_Enum = 3

	// User
	Role_ANONYMOUS_USER Role_Enum = 2
	Role_USER           Role_Enum = 4
)

var (
	Role_Enum_name = map[int32]string{
		0: "UNKNOWN_ROLE",
		1: "ANONYMOUS_ADMIN",
		2: "ANONYMOUS_USER",
		3: "ADMIN",
		4: "USER",
	}

	Role_Enum_value = map[string]int32{
		"UNKNOWN_ROLE":    0,
		"ANONYMOUS_ADMIN": 1,
		"ANONYMOUS_USER":  2,
		"ADMIN":           3,
		"USER":            4,
	}
)

// ===== RoleType Enums =====

type RoleType_Enum int32

const (
	RoleType_UNKNOWN RoleType_Enum = 0
	RoleType_ADMIN   RoleType_Enum = 1
	RoleType_USER    RoleType_Enum = 2
	RoleType_SYSTEM  RoleType_Enum = 9
)

var (
	RoleType_Enum_name = map[int32]string{
		0: "UNKNOWN",
		1: "ADMIN",
		2: "USER",
		9: "SYSTEM",
	}
	RoleType_Enum_value = map[string]int32{
		"UNKNOWN": 0,
		"ADMIN":   1,
		"USER":    2,
		"SYSTEM":  9,
	}
)

// ===== ControlStatus Enums =====

type ControlStatus_Enum int32

const (
	ControlStatus_UNKNOWN_CONTROL_STATUS ControlStatus_Enum = 0
	ControlStatus_ACTIVE                 ControlStatus_Enum = 1
	ControlStatus_INACTIVE               ControlStatus_Enum = 2
	ControlStatus_PENDING                ControlStatus_Enum = 3
	ControlStatus_LOCKED                 ControlStatus_Enum = 4
	ControlStatus_DRAFT                  ControlStatus_Enum = 5
)

var (
	ControlStatus_Enum_name = map[int32]string{
		0: "UNKNOWN_CONTROL_STATUS",
		1: "ACTIVE",
		2: "INACTIVE",
		3: "PENDING",
		4: "LOCKED",
		5: "DRAFT",
	}
	ControlStatus_Enum_value = map[string]int32{
		"UNKNOWN_CONTROL_STATUS": 0,
		"ACTIVE":                 1,
		"INACTIVE":               2,
		"PENDING":                3,
		"LOCKED":                 4,
		"DRAFT":                  5,
	}
)

// ===== DevicePlatform Enums =====

type DevicePlatform_Enum int32

const (
	DevicePlatform_UNKNOWN     DevicePlatform_Enum = 0
	DevicePlatform_ANDROID     DevicePlatform_Enum = 1
	DevicePlatform_IOS         DevicePlatform_Enum = 2
	DevicePlatform_WEB_BROWSER DevicePlatform_Enum = 3
)

var (
	DevicePlatform_Enum_name = map[int32]string{
		0: "UNKNOWN",
		1: "ANDROID",
		2: "IOS",
		3: "WEB_BROWSER",
	}
	DevicePlatform_Enum_value = map[string]int32{
		"UNKNOWN":     0,
		"ANDROID":     1,
		"IOS":         2,
		"WEB_BROWSER": 3,
	}
)

// ===== NotificationChannel Enums =====

type NotificationChannel_Enum int32

const (
	NotificationChannel_UNKNOWN NotificationChannel_Enum = 0
	NotificationChannel_FCM     NotificationChannel_Enum = 1
	NotificationChannel_APNS    NotificationChannel_Enum = 2
)

var (
	NotificationChannel_Enum_name = map[int32]string{
		0: "UNKNOWN",
		1: "FCM",
		2: "APNS",
	}
	NotificationChannel_Enum_value = map[string]int32{
		"UNKNOWN": 0,
		"FCM":     1,
		"APNS":    2,
	}
)

// ===== AssetType Enums =====

type AssetType_Enum int32

const (
	AssetType_USER_AVATAR AssetType_Enum = 1
)

// ===== AuthProvider Enums =====

type AuthProvider_Enum int32

const (
	AuthProvider_BASIC AuthProvider_Enum = 1
	AuthProvider_OAUTH AuthProvider_Enum = 2
)
