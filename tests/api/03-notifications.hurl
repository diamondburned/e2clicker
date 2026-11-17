# Notification preferences baseline

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
    ]
  }
}
HTTP 204

GET {{base_url}}/api/notifications/preferences
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.notificationConfigs.email" count == 1
jsonpath "$.notificationConfigs.email[0].address" == "test@example.com"
