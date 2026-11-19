# Notification preferences, Web Push device endpoints, SSRF checks, and test notifications

POST {{base_url}}/api/register
Content-Type: application/json
{
  "name": "hurl-notif-{{newUuid}}"
}
HTTP 200
[Captures]
user_secret: jsonpath "$.secret"

POST {{base_url}}/api/auth
Content-Type: application/json
{
  "secret": "{{user_secret}}"
}
HTTP 200
[Captures]
token: jsonpath "$.token"

GET {{base_url}}/api/notifications/preferences
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.notificationConfigs" exists

GET {{base_url}}/api/notifications/status
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.lastAttempt" not exists

# Test notification with no configured channels must fail with 400
POST {{base_url}}/api/notifications/test
Authorization: Bearer {{token}}
HTTP 400
[Asserts]
jsonpath "$.message" contains "no notification channels configured"

# SSRF / invalid Web Push endpoints must be rejected
PUT {{base_url}}/api/notifications/web-push/dev1
Authorization: Bearer {{token}}
Content-Type: application/json
{
  "deviceID": "dev1",
  "endpoint": "http://fcm.googleapis.com/fcm/send/insecure",
  "keys": {
    "p256dh": "BH1uVr_j2jFqnZFmOaqvBgwLZhMIRgGf23u-DfrN8yyr_0WX4UsgVIhWil6z__ddipWL4E8LqGgAc9yxIuTaJz8",
    "auth": "MDEyMzQ1Njc4OWFiY2RlZg"
  }
}
HTTP 400

PUT {{base_url}}/api/notifications/web-push/dev1
Authorization: Bearer {{token}}
Content-Type: application/json
{
  "deviceID": "dev1",
  "endpoint": "https://169.254.169.254/latest/meta-data/",
  "keys": {
    "p256dh": "BH1uVr_j2jFqnZFmOaqvBgwLZhMIRgGf23u-DfrN8yyr_0WX4UsgVIhWil6z__ddipWL4E8LqGgAc9yxIuTaJz8",
    "auth": "MDEyMzQ1Njc4OWFiY2RlZg"
  }
}
HTTP 400

PUT {{base_url}}/api/notifications/web-push/dev1
Authorization: Bearer {{token}}
Content-Type: application/json
{
  "deviceID": "dev1",
  "endpoint": "https://localhost/push",
  "keys": {
    "p256dh": "BH1uVr_j2jFqnZFmOaqvBgwLZhMIRgGf23u-DfrN8yyr_0WX4UsgVIhWil6z__ddipWL4E8LqGgAc9yxIuTaJz8",
    "auth": "MDEyMzQ1Njc4OWFiY2RlZg"
  }
}
HTTP 400

# Mismatched path vs body deviceID must be rejected
PUT {{base_url}}/api/notifications/web-push/dev1
Authorization: Bearer {{token}}
Content-Type: application/json
{
  "deviceID": "other-dev",
  "endpoint": "https://fcm.googleapis.com/fcm/send/test-sub",
  "keys": {
    "p256dh": "BH1uVr_j2jFqnZFmOaqvBgwLZhMIRgGf23u-DfrN8yyr_0WX4UsgVIhWil6z__ddipWL4E8LqGgAc9yxIuTaJz8",
    "auth": "MDEyMzQ1Njc4OWFiY2RlZg"
  }
}
HTTP 400

# Subscribe Web Push device
PUT {{base_url}}/api/notifications/web-push/dev1
Authorization: Bearer {{token}}
Content-Type: application/json
{
  "deviceID": "dev1",
  "endpoint": "https://fcm.googleapis.com/fcm/send/test-sub",
  "keys": {
    "p256dh": "BH1uVr_j2jFqnZFmOaqvBgwLZhMIRgGf23u-DfrN8yyr_0WX4UsgVIhWil6z__ddipWL4E8LqGgAc9yxIuTaJz8",
    "auth": "MDEyMzQ1Njc4OWFiY2RlZg"
  }
}
HTTP 204

# GET /notifications/preferences must return device metadata without endpoint or keys
GET {{base_url}}/api/notifications/preferences
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.notificationConfigs.webPush" count == 1
jsonpath "$.notificationConfigs.webPush[0].deviceID" == "dev1"
jsonpath "$.notificationConfigs.webPush[0].endpoint" not exists
jsonpath "$.notificationConfigs.webPush[0].keys" not exists

# Updating email preferences must preserve existing Web Push subscriptions
PUT {{base_url}}/api/notifications/preferences
Authorization: Bearer {{token}}
Content-Type: application/json
{
  "notificationConfigs": {
    "email": [
      {
        "address": "test@example.com",
        "name": "Test User"
      }
    ],
    "webPush": []
  }
}
HTTP 204

GET {{base_url}}/api/notifications/preferences
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.notificationConfigs.email" count == 1
jsonpath "$.notificationConfigs.email[0].address" == "test@example.com"
jsonpath "$.notificationConfigs.webPush" count == 1
jsonpath "$.notificationConfigs.webPush[0].deviceID" == "dev1"

# Unsubscribe Web Push device while preserving email preferences
DELETE {{base_url}}/api/notifications/web-push/dev1
Authorization: Bearer {{token}}
HTTP 204

GET {{base_url}}/api/notifications/preferences
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.notificationConfigs.email" count == 1
jsonpath "$.notificationConfigs.webPush" not exists
