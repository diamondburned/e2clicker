<!-- Generator: Widdershins v4.0.1 -->

<h1 id="e2clicker-service">e2clicker service v0</h1>

> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

Base URLs:

* <a href="/api">/api</a>

# Authentication

- HTTP Authentication, scheme: bearer 

<h1 id="e2clicker-service-default">Default</h1>

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
  "has_avatar": true,
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
  "has_avatar": true
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

`POST /me/avatar`

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
  "has_avatar": true
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
  "created_at": "2019-08-24T14:15:22Z",
  "last_used": "2019-08-24T14:15:22Z",
  "expires_at": "2019-08-24T14:15:22Z"
}

```

A session for a user.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|integer|true|none|The session identifier|
|created_at|string(date-time)|true|none|The time the session was created|
|last_used|string(date-time)|false|none|The last time the session was used|
|expires_at|string(date-time)|true|none|The time the session expires|

