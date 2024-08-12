<!-- Generator: Widdershins v4.0.1 -->

<h1 id="rest-api-definition-for-the-e2clicker-service-">REST API definition for the e2clicker service. v0</h1>

> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

Base URLs:

* <a href="/api/user/v1">/api/user/v1</a>

# Authentication

- HTTP Authentication, scheme: bearer 

<h1 id="rest-api-definition-for-the-e2clicker-service--default">Default</h1>

## login

<a id="opIdlogin"></a>

`POST /login`

> Body parameter

```json
{
  "email": "string",
  "password": "string"
}
```

<h3 id="login-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|object|false|none|
|» email|body|string|true|The username to log in with|
|» password|body|string|true|The password to log in with|

> Example responses

> 200 Response

```json
"string"
```

<h3 id="login-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully logged in.|[SessionToken](#schemasessiontoken)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|The username or password is incorrect.|None|

<aside class="success">
This operation does not require authentication
</aside>

## register

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

<h3 id="register-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|object|false|none|
|» name|body|string|true|The name to register with|
|» email|body|string|true|The username to register with|
|» password|body|string|true|The password to register with|

> Example responses

> 201 Response

```json
"string"
```

<h3 id="register-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|Successfully registered.|[SessionToken](#schemasessiontoken)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|The username is already taken.|None|

<aside class="success">
This operation does not require authentication
</aside>

## user

<a id="opIduser"></a>

`GET /user/{userID}`

<h3 id="user-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|userID|path|string|true|The ID of the user to get, or "me" to get the current user.|

#### Detailed descriptions

**userID**: The ID of the user to get, or "me" to get the current user.

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

<h3 id="user-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully retrieved the current user.|[User](#schemauser)|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## userAvatar

<a id="opIduserAvatar"></a>

`GET /user/{userID}/avatar`

<h3 id="useravatar-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|userID|path|string|true|The ID of the user to get the avatar for.|

#### Detailed descriptions

**userID**: The ID of the user to get the avatar for.

> Example responses

> 200 Response

<h3 id="useravatar-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully retrieved the user's avatar.|string|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

# Schemas

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

A session token string.
This is used in the Authorization header to authenticate requests.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|A session token string.<br>This is used in the Authorization header to authenticate requests.|

<h2 id="tocS_SessionTokenObject">SessionTokenObject</h2>
<!-- backwards compatibility -->
<a id="schemasessiontokenobject"></a>
<a id="schema_SessionTokenObject"></a>
<a id="tocSsessiontokenobject"></a>
<a id="tocssessiontokenobject"></a>

```json
{
  "userID": "string",
  "token": "string"
}

```

A session token object that is returned when creating a new session.

### Properties

*None*

