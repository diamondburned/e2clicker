// Package notification provides a way to send notifications to users.
package notification

import (
	"github.com/samber/do/v2"
)

var Package = do.Package(
	do.Lazy(do.InvokeStruct[NotificationService]),
	do.Lazy(do.InvokeStruct[UserNotificationService]),
	do.Lazy(do.InvokeStruct[GotifyService]),
	do.Lazy(do.InvokeStruct[PushoverService]),
)
