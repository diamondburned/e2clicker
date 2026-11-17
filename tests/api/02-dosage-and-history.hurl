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
[Captures]
dose1_id: jsonpath "$.id"
[Asserts]
jsonpath "$.id" matches "^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
jsonpath "$.deliveryMethod" == "EV im"
jsonpath "$.dose" == 5
jsonpath "$.takenAt" == "2025-01-10T12:00:00Z"

# Record a second dose at the exact same timestamp (allowed with UUIDv7 IDs)
POST {{base_url}}/api/dosage/dose
Authorization: Bearer {{token}}
Content-Type: application/json
{
  "takenAt": "2025-01-10T12:00:00Z",
  "comment": "second dose same second"
}
HTTP 200
[Captures]
dose2_id: jsonpath "$.id"
[Asserts]
jsonpath "$.id" matches "^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
jsonpath "$.id" != "{{dose1_id}}"
jsonpath "$.comment" == "second dose same second"

# Range > 1 year must be rejected
GET {{base_url}}/api/dosage?start=2023-01-01T00:00:00Z&end=2025-02-01T00:00:00Z
Authorization: Bearer {{token}}
HTTP 400

# Query within range returns both recorded doses
GET {{base_url}}/api/dosage?start=2025-01-01T00:00:00Z&end=2025-02-01T00:00:00Z
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.history" count == 2

# Edit first dose by UUID (leaves second dose untouched and persists comment)
PUT {{base_url}}/api/dosage/dose/{{dose1_id}}
Authorization: Bearer {{token}}
Content-Type: application/json
{
  "deliveryMethod": "EV im",
  "dose": 6,
  "takenAt": "2025-01-10T12:00:00Z",
  "comment": "first dose"
}
HTTP 204

# Delete second dose via single-item DELETE /dosage/dose/{doseID}
DELETE {{base_url}}/api/dosage/dose/{{dose2_id}}
Authorization: Bearer {{token}}
HTTP 204

# Reject format=estrannaise with Accept: application/json
GET {{base_url}}/api/dosage/export-doses?format=estrannaise
Authorization: Bearer {{token}}
Accept: application/json
HTTP 400

# Export as CSV (e2clicker v2 format)
GET {{base_url}}/api/dosage/export-doses
Authorization: Bearer {{token}}
Accept: text/csv
HTTP 200
[Asserts]
header "Content-Type" contains "text/csv"
body contains "id,deliveryMethod,dose,takenAt,takenOffAt,comment"
body contains "{{dose1_id}},EV im,6,2025-01-10T12:00:00Z,,first dose"

# Export as CSV (estrannaise format)
GET {{base_url}}/api/dosage/export-doses?format=estrannaise
Authorization: Bearer {{token}}
Accept: text/csv
HTTP 200
[Asserts]
header "Content-Type" contains "text/csv"
body contains "dose,days,model"
body contains "6,0,EV im"

# Export as JSON
GET {{base_url}}/api/dosage/export-doses
Authorization: Bearer {{token}}
Accept: application/json
HTTP 200
[Asserts]
header "Content-Type" contains "application/json"
jsonpath "$" count == 1
jsonpath "$[0].id" == "{{dose1_id}}"
jsonpath "$[0].dose" == 6
jsonpath "$[0].comment" == "first dose"

# Import v2 CSV (1 duplicate id + 1 new explicit id)
POST {{base_url}}/api/dosage/import-doses
Authorization: Bearer {{token}}
Content-Type: text/csv
```
id,deliveryMethod,dose,takenAt,takenOffAt,comment
{{dose1_id}},EV im,6,2025-01-10T12:00:00Z,,first dose
019450a0-0000-7000-8000-000000000010,EV im,4,2025-01-15T12:00:00Z,,imported v2 dose
```
HTTP 200
[Asserts]
jsonpath "$.records" == 2
jsonpath "$.succeeded" == 1

# Import legacy v1 CSV (1 new + 1 duplicate of same v1 row in-file)
POST {{base_url}}/api/dosage/import-doses
Authorization: Bearer {{token}}
Content-Type: text/csv
```
deliveryMethod,dose,takenAt,takenOffAt,comment
EV im,4,2025-01-20T12:00:00Z,,imported v1 dose
EV im,4,2025-01-20T12:00:00Z,,imported v1 dose
```
HTTP 200
[Asserts]
jsonpath "$.records" == 2
jsonpath "$.succeeded" == 1

# Import legacy ID-less v1 JSON (duplicate of v1 CSV row should insert 0)
POST {{base_url}}/api/dosage/import-doses
Authorization: Bearer {{token}}
Content-Type: application/json
[
  {
    "deliveryMethod": "EV im",
    "dose": 4,
    "takenAt": "2025-01-20T12:00:00Z",
    "comment": "imported v1 dose"
  }
]
HTTP 200
[Asserts]
jsonpath "$.records" == 1
jsonpath "$.succeeded" == 0

# Delete doses via query param ids
DELETE {{base_url}}/api/dosage/dose?ids={{dose1_id}}&ids=019450a0-0000-7000-8000-000000000010
Authorization: Bearer {{token}}
HTTP 204

GET {{base_url}}/api/dosage?start=2025-01-01T00:00:00Z&end=2025-02-01T00:00:00Z
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.history" count == 1

# Clear dosage schedule
DELETE {{base_url}}/api/dosage
Authorization: Bearer {{token}}
HTTP 204

# Second user to verify estrannaise CSV import rejection without hitting import rate limit
POST {{base_url}}/api/register
Content-Type: application/json
{
  "name": "hurl-dosage-estrannaise-{{newUuid}}"
}
HTTP 200
[Captures]
user2_secret: jsonpath "$.secret"

POST {{base_url}}/api/auth
Content-Type: application/json
{
  "secret": "{{user2_secret}}"
}
HTTP 200
[Captures]
token2: jsonpath "$.token"

POST {{base_url}}/api/dosage/import-doses
Authorization: Bearer {{token2}}
Content-Type: text/csv
```
dose,days,model
6,0,EV im
```
HTTP 400
