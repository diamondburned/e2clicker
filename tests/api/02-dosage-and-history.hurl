# Dosage schedule, dose history CRUD, and CSV/JSON export & import

POST {{base_url}}/api/register
Content-Type: application/json
{
  "name": "hurl-dosage-{{newUuid}}"
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

# Empty dosage before configuration
GET {{base_url}}/api/dosage
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.dosage" not exists
jsonpath "$.history" not exists

# Invalid delivery method rejected
PUT {{base_url}}/api/dosage
Authorization: Bearer {{token}}
Content-Type: application/json
{
  "deliveryMethod": "not-a-real-method",
  "dose": 5,
  "interval": 5
}
HTTP 400

# Configure EV im regimen
PUT {{base_url}}/api/dosage
Authorization: Bearer {{token}}
Content-Type: application/json
{
  "deliveryMethod": "EV im",
  "dose": 5,
  "interval": 5
}
HTTP 204

GET {{base_url}}/api/dosage
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.dosage.deliveryMethod" == "EV im"
jsonpath "$.dosage.dose" == 5
jsonpath "$.dosage.interval" == 5

# Record a dose with an explicit timestamp
POST {{base_url}}/api/dosage/dose
Authorization: Bearer {{token}}
Content-Type: application/json
{
  "takenAt": "2025-01-10T12:00:00Z"
}
HTTP 200
[Asserts]
jsonpath "$.deliveryMethod" == "EV im"
jsonpath "$.dose" == 5
jsonpath "$.takenAt" == "2025-01-10T12:00:00Z"

# Range > 1 year must be rejected
GET {{base_url}}/api/dosage?start=2023-01-01T00:00:00Z&end=2025-02-01T00:00:00Z
Authorization: Bearer {{token}}
HTTP 400

# Query within range returns the recorded dose
GET {{base_url}}/api/dosage?start=2025-01-01T00:00:00Z&end=2025-02-01T00:00:00Z
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.history" count == 1
jsonpath "$.history[0].dose" == 5
jsonpath "$.history[0].takenAt" == "2025-01-10T12:00:00Z"

# Edit dose
PUT {{base_url}}/api/dosage/dose/2025-01-10T12:00:00Z
Authorization: Bearer {{token}}
Content-Type: application/json
{
  "deliveryMethod": "EV im",
  "dose": 6,
  "takenAt": "2025-01-10T12:00:00Z"
}
HTTP 204

# Export as CSV
GET {{base_url}}/api/dosage/export-doses
Authorization: Bearer {{token}}
Accept: text/csv
HTTP 200
[Asserts]
header "Content-Type" contains "text/csv"
body contains "deliveryMethod,dose,takenAt,takenOffAt,comment"
body contains "EV im,6,2025-01-10T12:00:00Z"

# Export as JSON
GET {{base_url}}/api/dosage/export-doses
Authorization: Bearer {{token}}
Accept: application/json
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$" count == 1
jsonpath "$[0].dose" == 6

# Import CSV (1 duplicate + 1 new)
POST {{base_url}}/api/dosage/import-doses
Authorization: Bearer {{token}}
Content-Type: text/csv
```
deliveryMethod,dose,takenAt,takenOffAt,comment
EV im,6,2025-01-10T12:00:00Z,,
EV im,4,2025-01-15T12:00:00Z,,imported dose
```
HTTP 200
[Asserts]
jsonpath "$.records" == 2
jsonpath "$.succeeded" == 1

# Delete doses via query param
DELETE {{base_url}}/api/dosage/dose?doseTimes=2025-01-10T12:00:00Z&doseTimes=2025-01-15T12:00:00Z
Authorization: Bearer {{token}}
HTTP 204

GET {{base_url}}/api/dosage?start=2025-01-01T00:00:00Z&end=2025-02-01T00:00:00Z
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.history" count == 0

# Clear dosage schedule
DELETE {{base_url}}/api/dosage
Authorization: Bearer {{token}}
HTTP 204
