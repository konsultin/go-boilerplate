package constant

import "github.com/konsultin/project-goes-here/libs/errk"

var (
	CurrentAuthSessionExpired  = errk.NewError("E_AUTH_2", "Current session has expired")
	LoginDetectedAnotherDevice = errk.NewError("E_AUTH_3", "Login detected from another device")
	ResourceNotFound           = errk.NewError("E_NOTFOUND", "Resource not found")
)
