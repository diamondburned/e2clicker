<!-- Generator: Widdershins v4.0.1 -->

<h1 id="e2clicker-service">e2clicker service v0</h1>

> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

Base URLs:

* <a href="/api/user/v1">/api/user/v1</a>

# Authentication

- HTTP Authentication, scheme: bearer 

<h1 id="e2clicker-service-default">Default</h1>

## Log into an existing account

<a id="opIdlogin"></a>

`POST /login`

> Body parameter

```json
{
  "email": "string",
  "password": "string"
}
```

<h3 id="log-into-an-existing-account-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|User-Agent|header|string|false|The user agent string of the client making the request.|
|body|body|object|false|none|
|» email|body|string|true|The username to log in with|
|» password|body|string|true|The password to log in with|

> Example responses

> 200 Response

```json
{
  "userID": "string",
  "token": "string"
}
```

<h3 id="log-into-an-existing-account-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully logged in.|Inline|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<h3 id="log-into-an-existing-account-responseschema">Response Schema</h3>

<aside class="success">
This operation does not require authentication
</aside>

## Register a new account

<a id="opIdregister"></a>

`POST /register`

> Body parameter

```json
{
  "name": "string",
  "email": "string",
  "password": "string"
}
```

<h3 id="register-a-new-account-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|User-Agent|header|string|false|The user agent string of the client making the request.|
|body|body|object|false|none|
|» name|body|string|true|The name to register with|
|» email|body|string|true|The username to register with|
|» password|body|string|true|The password to register with|

> Example responses

> 200 Response

```json
{
  "user": {
    "id": "string",
    "email": "string",
    "name": "string",
    "locale": "string"
  },
  "token": "string"
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

## Get a user by ID

<a id="opIduser"></a>

`GET /user/{userID}`

<h3 id="get-a-user-by-id-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|userID|path|string|true|The ID of the user to get the avatar for.|

> Example responses

> 200 Response

```json
{
  "id": "string",
  "email": "string",
  "name": "string",
  "locale": "string"
}
```

<h3 id="get-a-user-by-id-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully retrieved the current user.|[User](#schemauser)|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Get a user's avatar by ID

<a id="opIduserAvatar"></a>

`GET /user/{userID}/avatar`

<h3 id="get-a-user's-avatar-by-id-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|userID|path|string|true|The ID of the user to get the avatar for.|

> Example responses

> 200 Response

<h3 id="get-a-user's-avatar-by-id-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully retrieved the user's avatar.|string|

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

<h2 id="tocS_UserID">UserID</h2>
<!-- backwards compatibility -->
<a id="schemauserid"></a>
<a id="schema_UserID"></a>
<a id="tocSuserid"></a>
<a id="tocsuserid"></a>

```json
"string"

```

A unique user identifier.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|A unique user identifier.|

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
  "id": "string",
  "email": "string",
  "name": "string",
  "locale": "string"
}

```

A user of the system.

### Properties

*None*

<h2 id="tocS_SessionToken">SessionToken</h2>
<!-- backwards compatibility -->
<a id="schemasessiontoken"></a>
<a id="schema_SessionToken"></a>
<a id="tocSsessiontoken"></a>
<a id="tocssessiontoken"></a>

```json
"string"

```

A session token string. This is used in the Authorization header to authenticate requests.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|A session token string. This is used in the Authorization header to authenticate requests.|

