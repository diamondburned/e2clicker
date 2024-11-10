<!-- Generator: Widdershins v4.0.1 -->

<h1 id="e2clicker-service">e2clicker service v0</h1>

> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

Base URLs:

* <a href="/api">/api</a>

# Authentication

- HTTP Authentication, scheme: bearer 

<h1 id="e2clicker-service-dosage">dosage</h1>

## List all available delivery methods

<a id="opIddeliveryMethods"></a>

`GET /delivery-methods`

> Example responses

> 200 Response

```json
[
  {
    "id": "string",
    "units": "string",
    "name": "string"
  }
]
```

<h3 id="list-all-available-delivery-methods-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully retrieved delivery methods.|Inline|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<h3 id="list-all-available-delivery-methods-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[[DeliveryMethod](#schemadeliverymethod)]|false|none|none|

<aside class="success">
This operation does not require authentication
</aside>

## Get the user's dosage schedule

<a id="opIddosageSchedule"></a>

`GET /dosage/schedule`

> Example responses

> 200 Response

```json
{
  "schedule": {
    "deliveryMethod": "string",
    "dose": 0.1,
    "interval": 0.1,
    "concurrence": 0
  }
}
```

<h3 id="get-the-user's-dosage-schedule-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully retrieved dosage schedule.|Inline|

<h3 id="get-the-user's-dosage-schedule-responseschema">Response Schema</h3>

Status Code **200**

*The user's dosage schedule. If the user has no schedule set, this will be null.*

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Set the user's dosage schedule

<a id="opIdsetDosageSchedule"></a>

`PUT /dosage/schedule`

> Body parameter

```json
{
  "deliveryMethod": "string",
  "dose": 0.1,
  "interval": 0.1,
  "concurrence": 0
}
```

<h3 id="set-the-user's-dosage-schedule-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[DosageSchedule](#schemadosageschedule)|false|none|

> Example responses

> default Response

```json
{
  "message": "string",
  "details": null,
  "internal": true,
  "internalCode": "string"
}
```

<h3 id="set-the-user's-dosage-schedule-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Successfully set dosage schedule.|None|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Clear the user's dosage schedule

<a id="opIdclearDosageSchedule"></a>

`DELETE /dosage/schedule`

> Example responses

> default Response

```json
{
  "message": "string",
  "details": null,
  "internal": true,
  "internalCode": "string"
}
```

<h3 id="clear-the-user's-dosage-schedule-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Successfully cleared dosage schedule.|None|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Get the user's dosage history within a time range

<a id="opIddoseHistory"></a>

`GET /dosage/history`

<h3 id="get-the-user's-dosage-history-within-a-time-range-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|start|query|string(date-time)|true|none|
|end|query|string(date-time)|true|none|

> Example responses

> 200 Response

```json
{
  "history": [
    {
      "id": 0,
      "deliveryMethod": "string",
      "dose": 0.1,
      "takenAt": "2019-08-24T14:15:22Z",
      "takenOffAt": "2019-08-24T14:15:22Z"
    }
  ]
}
```

<h3 id="get-the-user's-dosage-history-within-a-time-range-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully retrieved dosage history.|[DosageHistory](#schemadosagehistory)|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Record a new dosage to the user's history

<a id="opIdrecordDose"></a>

`POST /dosage/history`

> Body parameter

```json
{
  "takenAt": "2019-08-24T14:15:22Z"
}
```

<h3 id="record-a-new-dosage-to-the-user's-history-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|object|false|none|
|» takenAt|body|string(date-time)|true|The time the dosage was taken.|

> Example responses

> 200 Response

```json
{
  "id": 0,
  "deliveryMethod": "string",
  "dose": 0.1,
  "takenAt": "2019-08-24T14:15:22Z",
  "takenOffAt": "2019-08-24T14:15:22Z"
}
```

<h3 id="record-a-new-dosage-to-the-user's-history-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully recorded dosage.|[DosageObservation](#schemadosageobservation)|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Update a dosage in the user's history

<a id="opIdeditDose"></a>

`PUT /dosage/history`

> Body parameter

```json
{
  "id": 0,
  "deliveryMethod": "string",
  "dose": 0.1,
  "takenAt": "2019-08-24T14:15:22Z",
  "takenOffAt": "2019-08-24T14:15:22Z"
}
```

<h3 id="update-a-dosage-in-the-user's-history-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[DosageObservation](#schemadosageobservation)|false|none|

> Example responses

> default Response

```json
{
  "message": "string",
  "details": null,
  "internal": true,
  "internalCode": "string"
}
```

<h3 id="update-a-dosage-in-the-user's-history-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Successfully updated dosage.|None|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Delete multiple dosages from the user's history

<a id="opIdforgetDoses"></a>

`DELETE /dosage/history`

> Body parameter

```json
{
  "dose_ids": [
    0
  ]
}
```

<h3 id="delete-multiple-dosages-from-the-user's-history-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|object|false|none|
|» dose_ids|body|[integer]|true|none|

> Example responses

> default Response

```json
{
  "message": "string",
  "details": null,
  "internal": true,
  "internalCode": "string"
}
```

<h3 id="delete-multiple-dosages-from-the-user's-history-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Successfully deleted dosages.|None|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

<h1 id="e2clicker-service-user">user</h1>

## Register a new account

<a id="opIdregister"></a>

`POST /register`

> Body parameter

```json
{
  "name": "string"
}
```

<h3 id="register-a-new-account-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|object|false|none|
|» name|body|string|true|The name to register with|

> Example responses

> 200 Response

```json
{
  "name": "string",
  "locale": "string",
  "hasAvatar": true,
  "secret": "string"
}
```

<h3 id="register-a-new-account-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully logged in.|Inline|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<h3 id="register-a-new-account-responseschema">Response Schema</h3>

<aside class="success">
This operation does not require authentication
</aside>

## Authenticate a user and obtain a session

<a id="opIdauth"></a>

`POST /auth`

> Body parameter

```json
{
  "secret": "string"
}
```

<h3 id="authenticate-a-user-and-obtain-a-session-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|User-Agent|header|string|false|The user agent of the client making the request.|
|body|body|object|false|none|

> Example responses

> 200 Response

```json
{
  "token": "string"
}
```

<h3 id="authenticate-a-user-and-obtain-a-session-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully logged in.|Inline|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<h3 id="authenticate-a-user-and-obtain-a-session-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» token|string|true|none|The session token|

<aside class="success">
This operation does not require authentication
</aside>

## Get the current user

<a id="opIdcurrentUser"></a>

`GET /me`

> Example responses

> 200 Response

```json
{
  "name": "string",
  "locale": "string",
  "hasAvatar": true
}
```

<h3 id="get-the-current-user-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully retrieved the current user.|[User](#schemauser)|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Get the current user's avatar

<a id="opIdcurrentUserAvatar"></a>

`GET /me/avatar`

> Example responses

> 200 Response

> default Response

```json
{
  "message": "string",
  "details": null,
  "internal": true,
  "internalCode": "string"
}
```

<h3 id="get-the-current-user's-avatar-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully retrieved the user's avatar.|string|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Set the current user's avatar

<a id="opIdsetCurrentUserAvatar"></a>

`PUT /me/avatar`

> Body parameter

<h3 id="set-the-current-user's-avatar-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|string(binary)|false|none|

> Example responses

> default Response

```json
{
  "message": "string",
  "details": null,
  "internal": true,
  "internalCode": "string"
}
```

<h3 id="set-the-current-user's-avatar-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Successfully set the user's avatar.|None|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Get the current user's secret

<a id="opIdcurrentUserSecret"></a>

`GET /me/secret`

> Example responses

> 200 Response

```json
{
  "secret": "string"
}
```

<h3 id="get-the-current-user's-secret-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully retrieved the user's secret.|Inline|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<h3 id="get-the-current-user's-secret-responseschema">Response Schema</h3>

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## List the current user's sessions

<a id="opIdcurrentUserSessions"></a>

`GET /me/sessions`

> Example responses

> 200 Response

```json
[
  {
    "id": 0,
    "createdAt": "2019-08-24T14:15:22Z",
    "lastUsed": "2019-08-24T14:15:22Z",
    "expiresAt": "2019-08-24T14:15:22Z"
  }
]
```

<h3 id="list-the-current-user's-sessions-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully retrieved the user's sessions.|Inline|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<h3 id="list-the-current-user's-sessions-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[[Session](#schemasession)]|false|none|[A session for a user.]|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Delete one of the current user's sessions

<a id="opIddeleteUserSession"></a>

`DELETE /me/sessions`

> Body parameter

```json
{
  "id": 0
}
```

<h3 id="delete-one-of-the-current-user's-sessions-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|object|false|none|
|» id|body|integer(int64)|true|The session identifier to delete|

> Example responses

> default Response

```json
{
  "message": "string",
  "details": null,
  "internal": true,
  "internalCode": "string"
}
```

<h3 id="delete-one-of-the-current-user's-sessions-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Successfully deleted the user's sessions.|None|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

# Schemas

<h2 id="tocS_Error">Error</h2>
<!-- backwards compatibility -->
<a id="schemaerror"></a>
<a id="schema_Error"></a>
<a id="tocSerror"></a>
<a id="tocserror"></a>

```json
{
  "message": "string",
  "details": null,
  "internal": true,
  "internalCode": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|message|string|true|none|A message describing the error|
|details|any|false|none|Additional details about the error|
|internal|boolean|false|none|Whether the error is internal|
|internalCode|string|false|none|An internal code for the error (useless for clients)|

<h2 id="tocS_DeliveryMethod">DeliveryMethod</h2>
<!-- backwards compatibility -->
<a id="schemadeliverymethod"></a>
<a id="schema_DeliveryMethod"></a>
<a id="tocSdeliverymethod"></a>
<a id="tocsdeliverymethod"></a>

```json
{
  "id": "string",
  "units": "string",
  "name": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|true|none|A short string representing the delivery method. This is what goes into the DeliveryMethod fields.|
|units|string|true|none|The units of the delivery method.|
|name|string|true|none|The full name of the delivery method.|

<h2 id="tocS_DosageSchedule">DosageSchedule</h2>
<!-- backwards compatibility -->
<a id="schemadosageschedule"></a>
<a id="schema_DosageSchedule"></a>
<a id="tocSdosageschedule"></a>
<a id="tocsdosageschedule"></a>

```json
{
  "deliveryMethod": "string",
  "dose": 0.1,
  "interval": 0.1,
  "concurrence": 0
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|deliveryMethod|string|true|none|The delivery method to use.|
|dose|number(float)|true|none|The dosage amount.|
|interval|number(double)|true|none|The interval between doses in days.|
|concurrence|integer|false|none|The number of estrogen patches on the body at once. Only relevant if delivery method is patch.|

<h2 id="tocS_DosageHistory">DosageHistory</h2>
<!-- backwards compatibility -->
<a id="schemadosagehistory"></a>
<a id="schema_DosageHistory"></a>
<a id="tocSdosagehistory"></a>
<a id="tocsdosagehistory"></a>

```json
{
  "history": [
    {
      "id": 0,
      "deliveryMethod": "string",
      "dose": 0.1,
      "takenAt": "2019-08-24T14:15:22Z",
      "takenOffAt": "2019-08-24T14:15:22Z"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|history|[[DosageObservation](#schemadosageobservation)]|true|none|none|

<h2 id="tocS_DosageObservation">DosageObservation</h2>
<!-- backwards compatibility -->
<a id="schemadosageobservation"></a>
<a id="schema_DosageObservation"></a>
<a id="tocSdosageobservation"></a>
<a id="tocsdosageobservation"></a>

```json
{
  "id": 0,
  "deliveryMethod": "string",
  "dose": 0.1,
  "takenAt": "2019-08-24T14:15:22Z",
  "takenOffAt": "2019-08-24T14:15:22Z"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|integer(int64)|true|none|The unique identifier for the observation.|
|deliveryMethod|string|true|none|The delivery method used.|
|dose|number(float)|true|none|The dosage amount.|
|takenAt|string(date-time)|true|none|The time the dosage was taken.|
|takenOffAt|string(date-time)|false|none|The time the dosage was taken off. This is only relevant for patch delivery methods.|

<h2 id="tocS_UserSecret">UserSecret</h2>
<!-- backwards compatibility -->
<a id="schemausersecret"></a>
<a id="schema_UserSecret"></a>
<a id="tocSusersecret"></a>
<a id="tocsusersecret"></a>

```json
"string"

```

A secret and unique user identifier. This secret is generated once and never changes. It is used to both authenticate and identify a user, so it should be kept secret.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|A secret and unique user identifier. This secret is generated once and never changes. It is used to both authenticate and identify a user, so it should be kept secret.|

<h2 id="tocS_Locale">Locale</h2>
<!-- backwards compatibility -->
<a id="schemalocale"></a>
<a id="schema_Locale"></a>
<a id="tocSlocale"></a>
<a id="tocslocale"></a>

```json
"string"

```

A locale identifier.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|A locale identifier.|

<h2 id="tocS_User">User</h2>
<!-- backwards compatibility -->
<a id="schemauser"></a>
<a id="schema_User"></a>
<a id="tocSuser"></a>
<a id="tocsuser"></a>

```json
{
  "name": "string",
  "locale": "string",
  "hasAvatar": true
}

```

A user of the system.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|name|string|true|none|The user's name|

<h2 id="tocS_Session">Session</h2>
<!-- backwards compatibility -->
<a id="schemasession"></a>
<a id="schema_Session"></a>
<a id="tocSsession"></a>
<a id="tocssession"></a>

```json
{
  "id": 0,
  "createdAt": "2019-08-24T14:15:22Z",
  "lastUsed": "2019-08-24T14:15:22Z",
  "expiresAt": "2019-08-24T14:15:22Z"
}

```

A session for a user.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|integer(int64)|true|none|The session identifier|
|createdAt|string(date-time)|true|none|The time the session was created|
|lastUsed|string(date-time)|true|none|The last time the session was used|
|expiresAt|string(date-time)|false|none|The time the session expires, or null if it never expires|

