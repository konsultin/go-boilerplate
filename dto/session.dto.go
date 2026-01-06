package dto

type CreateUserSession_Result struct {
	Session      *Session `json:"session,omitempty"`
	RefreshToken *Session `json:"refreshToken,omitempty"`
	Scopes       []string `json:"scopes,omitempty"`
}

type CreateUserSession_Result_Data struct {
	User           *User    `json:"user"`
	AccessSession  *Session `json:"accessSession"`
	RefreshSession *Session `json:"refreshSession"`
	AccessScopes   []string `json:"accessScopes"`
}

type Session struct {
	Token     string `json:"token,omitempty"`
	ExpiredAt int64  `json:"expiredAt,omitempty"`
}

type RefreshSession_Payload struct {
	RefreshToken string         `json:"refreshToken" validate:"required"`
	Device       *DeviceSession `json:"device,omitempty" validate:"omitempty"`
}

type DeviceSession struct {
	DeviceId              string                   `json:"deviceId,omitempty" validate:"omitempty,max=128"`
	DevicePlatformId      DevicePlatform_Enum      `json:"devicePlatformId,omitempty"`
	ClientIP              string                   `json:"clientIP,omitempty" validate:"omitempty,ip"`
	NotificationChannelId NotificationChannel_Enum `json:"notificationChannelId,omitempty"`
	NotificationToken     string                   `json:"notificationChannelToken,omitempty" validate:"omitempty,max=512"`
}
