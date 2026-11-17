# Unauthenticated discovery endpoints & user/session lifecycle

GET {{base_url}}/api/delivery-methods
HTTP 200
[Asserts]
jsonpath "$" isCollection
jsonpath "$[?(@.id == 'EV im')]" count == 1
jsonpath "$[?(@.id == 'patch tw')]" count == 1

GET {{base_url}}/api/notifications/methods
HTTP 200
[Asserts]
jsonpath "$" contains "webPush"

GET {{base_url}}/api/push-info
HTTP 200
[Asserts]
jsonpath "$.applicationServerKey" isString
jsonpath "$.applicationServerKey" != ""

# Protected route without Authorization header must fail
GET {{base_url}}/api/me
HTTP 403

# Register a fresh user
POST {{base_url}}/api/register
Content-Type: application/json
{
  "name": "hurl-user-{{newUuid}}"
}
HTTP 200
[Captures]
user_secret: jsonpath "$.secret"
user_name: jsonpath "$.name"
[Asserts]
jsonpath "$.secret" isString
jsonpath "$.locale" isString

# Authenticate with the secret
POST {{base_url}}/api/auth
Content-Type: application/json
User-Agent: hurl-e2e-test/1.0
{
  "secret": "{{user_secret}}"
}
HTTP 200
[Captures]
token: jsonpath "$.token"
[Asserts]
jsonpath "$.token" != ""

# Fetch current user profile
GET {{base_url}}/api/me
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.secret" == "{{user_secret}}"
jsonpath "$.name" == "{{user_name}}"

# List sessions & delete current session
GET {{base_url}}/api/me/sessions
Authorization: Bearer {{token}}
HTTP 200
[Captures]
session_id: jsonpath "$[0].id"
[Asserts]
jsonpath "$" count == 1

DELETE {{base_url}}/api/me/sessions?id={{session_id}}
Authorization: Bearer {{token}}
HTTP 204

# Deleted session token is rejected with 401
GET {{base_url}}/api/me
Authorization: Bearer {{token}}
HTTP 401
