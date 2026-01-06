package model

import (
	"database/sql"
	"time"

	"github.com/konsultin/project-goes-here/dto"
)

type AuthSession struct {
	BaseField
	Id                    int64                        `db:"id"`
	Xid                   string                       `db:"xid"`
	SubjectId             string                       `db:"subject_id"`
	SubjectTypeId         dto.Role_Enum                `db:"subject_type_id"`
	AuthProviderId        dto.AuthProvider_Enum        `db:"auth_provider_id"`
	DevicePlatformId      dto.DevicePlatform_Enum      `db:"device_platform_id"`
	DeviceId              string                       `db:"device_id"`
	Device                *AuthSessionDevice           `db:"device"`
	NotificationChannelId dto.NotificationChannel_Enum `db:"notification_channel_id"`
	NotificationToken     sql.NullString               `db:"notification_token"`
	ExpiredAt             time.Time                    `db:"expired_at"`
	StatusId              dto.ControlStatus_Enum       `db:"status_id"`
}

type AuthSessionDevice struct {
	DeviceId         string                  `json:"deviceId"`
	DevicePlatformId dto.DevicePlatform_Enum `json:"devicePlatformId"`
	ClientIp         string                  `json:"clientIp"`
}
