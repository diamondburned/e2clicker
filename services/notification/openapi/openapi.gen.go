// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version (devel) DO NOT EDIT.
package openapi

import (
	"time"
)

// Defines values for NotificationType.
const (
	AccountNotice   NotificationType = "account_notice"
	Reminder        NotificationType = "reminder"
	WebPushExpiring NotificationType = "web_push_expiring"
)

// PushDeviceID A short ID associated with the device that the push subscription is for This is used to identify the device when updating its push subscription later on.
// Realistically, this will be handled as an opaque random string generated on the device side, so the server has no way to correlate  it with any fingerprinting.
// The recommended way to generate this string in JavaScript is:
// ```js crypto.randomUUID().slice(0, 8) ```
type PushDeviceID = string

// Notification defines model for Notification.
type Notification struct {
	// Type The type of notification.
	Type NotificationType `json:"type"`

	// Message The message of the notification.
	Message NotificationMessage `json:"message"`

	// Username The username of the user to send the notification to.
	Username string `json:"username"`
}

// NotificationType The type of notification.
type NotificationType string

// NotificationMessage The message of the notification. This is derived from the notification type but can be overridden by the user.
type NotificationMessage struct {
	// Title The title of the notification.
	Title string `json:"title"`

	// Message The message of the notification.
	Message string `json:"message"`
}

// PushInfo This is returned by the server and contains information that the client would need to subscribe to push notifications.
type PushInfo struct {
	// ApplicationServerKey A Base64-encoded string or ArrayBuffer containing an ECDSA P-256 public key that the push server will use to authenticate your application server. If specified, all messages from your application server must use the VAPID authentication scheme, and include a JWT signed with the corresponding private key. This key IS NOT the same ECDH key that you use to encrypt the data. For more information, see "Using VAPID with WebPush".
	ApplicationServerKey string `json:"applicationServerKey"`
}

// PushSubscription The configuration for a push notification subscription.
// This is the object that is returned by calling PushSubscription.toJSON(). More information can be found at: https://developer.mozilla.org/en-US/docs/Web/API/PushSubscription/toJSON
type PushSubscription struct {
	// DeviceID A short ID associated with the device that the push subscription is for This is used to identify the device when updating its push subscription later on.
	// Realistically, this will be handled as an opaque random string generated on the device side, so the server has no way to correlate  it with any fingerprinting.
	// The recommended way to generate this string in JavaScript is:
	// ```js crypto.randomUUID().slice(0, 8) ```
	DeviceID PushDeviceID `json:"deviceID"`

	// Endpoint The endpoint to send the notification to.
	Endpoint string `json:"endpoint"`

	// ExpirationTime The time at which the subscription expires. This is the time when the subscription will be automatically deleted by the browser.
	ExpirationTime *time.Time `json:"expirationTime,omitempty"`

	// Keys The VAPID keys to encrypt the push notification.
	Keys struct {
		// P256Dh An Elliptic curve Diffie–Hellman public key on the P-256 curve (that is, the NIST secp256r1 elliptic curve). The resulting key is an uncompressed point in ANSI X9.62 format.
		P256Dh string `json:"p256dh"`

		// Auth An authentication secret, as described in Message Encryption for Web Push.
		Auth string `json:"auth"`
	} `json:"keys"`
}

// ReturnedPushSubscription Similar to a [PushSubscription], but specifically for returning to the user. This type contains no secrets.
type ReturnedPushSubscription struct {
	DeviceID PushDeviceID `json:"deviceID"`
	Keys     struct {
		// P256Dh An Elliptic curve Diffie–Hellman public key on the P-256 curve (that is, the NIST secp256r1 elliptic curve). The resulting key is an uncompressed point in ANSI X9.62 format.
		P256Dh string `json:"p256dh"`
	} `json:"keys"`

	// ExpirationTime The time at which the subscription expires. This is the time when the subscription will be automatically deleted by the browser.
	ExpirationTime *time.Time `json:"expirationTime,omitempty"`
}

// UserUnsubscribePushParams defines parameters for UserUnsubscribePush.
type UserUnsubscribePushParams struct {
	// DeviceID The device ID of the push subscription to unsubscribe from.
	DeviceID PushDeviceID `form:"deviceID" json:"deviceID"`
}

// UserSubscribePushJSONRequestBody defines body for UserSubscribePush for application/json ContentType.
type UserSubscribePushJSONRequestBody = PushSubscription
