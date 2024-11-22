<!-- Generator: Widdershins v4.0.1 -->

<h1 id="e2clicker-service">e2clicker service v0</h1>

> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

Base URLs:

* <a href="https://e2clicker.app/api">https://e2clicker.app/api</a>

* <a href="/api">/api</a>

# Authentication

- HTTP Authentication, scheme: bearer 

<h1 id="e2clicker-service-dosage">dosage</h1>

## List all available delivery methods

<a id="opIddeliveryMethods"></a>

`GET /deliverymethods`

> Example responses

> 200 Response

```json
[
  {
    "id": "string",
    "units": "string",
    "name": "string",
    "description": "string"
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

## Get the user's dosage and optionally their history

<a id="opIddosage"></a>

`GET /dosage`

<h3 id="get-the-user's-dosage-and-optionally-their-history-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|start|query|string(date-time)|false|none|
|end|query|string(date-time)|false|none|

> Example responses

> 200 Response

```json
{
  "dosage": {
    "deliveryMethod": "string",
    "dose": 0.1,
    "interval": 0.1,
    "concurrence": 0
  },
  "history": [
    {
      "id": 0,
      "deliveryMethod": "string",
      "dose": 0.1,
      "takenAt": "2019-08-24T14:15:22Z",
      "takenOffAt": "2019-08-24T14:15:22Z",
      "comment": "string"
    }
  ]
}
```

<h3 id="get-the-user's-dosage-and-optionally-their-history-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully retrieved the dosage, if set.|Inline|

<h3 id="get-the-user's-dosage-and-optionally-their-history-responseschema">Response Schema</h3>

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Set the user's dosage

<a id="opIdsetDosage"></a>

`PUT /dosage`

> Body parameter

```json
{
  "deliveryMethod": "string",
  "dose": 0.1,
  "interval": 0.1,
  "concurrence": 0
}
```

<h3 id="set-the-user's-dosage-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[Dosage](#schemadosage)|true|none|

> Example responses

> default Response

```json
{
  "message": "string",
  "errors": [
    {
      "message": "string",
      "errors": [],
      "details": null,
      "internal": true,
      "internalCode": "string"
    }
  ],
  "details": null,
  "internal": true,
  "internalCode": "string"
}
```

<h3 id="set-the-user's-dosage-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Successfully set the dosage.|None|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Clear the user's dosage schedule

<a id="opIdclearDosage"></a>

`DELETE /dosage`

> Example responses

> default Response

```json
{
  "message": "string",
  "errors": [
    {
      "message": "string",
      "errors": [],
      "details": null,
      "internal": true,
      "internalCode": "string"
    }
  ],
  "details": null,
  "internal": true,
  "internalCode": "string"
}
```

<h3 id="clear-the-user's-dosage-schedule-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Successfully cleared the dosage.|None|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Record a new dosage to the user's history

<a id="opIdrecordDose"></a>

`POST /dosage/dose`

This endpoint is used to record a new dosage observation to the user's history. The current time is automatically used.

> Example responses

> 200 Response

```json
{
  "id": 0,
  "deliveryMethod": "string",
  "dose": 0.1,
  "takenAt": "2019-08-24T14:15:22Z",
  "takenOffAt": "2019-08-24T14:15:22Z",
  "comment": "string"
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

`PUT /dosage/dose`

> Body parameter

```json
{
  "id": 0,
  "deliveryMethod": "string",
  "dose": 0.1,
  "takenAt": "2019-08-24T14:15:22Z",
  "takenOffAt": "2019-08-24T14:15:22Z",
  "comment": "string"
}
```

<h3 id="update-a-dosage-in-the-user's-history-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|any|true|none|

> Example responses

> default Response

```json
{
  "message": "string",
  "errors": [
    {
      "message": "string",
      "errors": [],
      "details": null,
      "internal": true,
      "internalCode": "string"
    }
  ],
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

`DELETE /dosage/dose`

<h3 id="delete-multiple-dosages-from-the-user's-history-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|dose_ids|query|array[integer]|true|none|

> Example responses

> default Response

```json
{
  "message": "string",
  "errors": [
    {
      "message": "string",
      "errors": [],
      "details": null,
      "internal": true,
      "internalCode": "string"
    }
  ],
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

## Export the user's dosage history

<a id="opIdexportDosageHistory"></a>

`GET /dosage/export`

<h3 id="export-the-user's-dosage-history-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Accept|header|string|true|The format to export the dosage history in.|
|start|query|string(date-time)|false|none|
|end|query|string(date-time)|false|none|

#### Enumerated Values

|Parameter|Value|
|---|---|
|Accept|text/csv|
|Accept|application/json|

> Example responses

> 200 Response

```json
[
  {
    "id": 0,
    "deliveryMethod": "string",
    "dose": 0.1,
    "takenAt": "2019-08-24T14:15:22Z",
    "takenOffAt": "2019-08-24T14:15:22Z",
    "comment": "string"
  }
]
```

<h3 id="export-the-user's-dosage-history-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully retrieved the dosage history.|[DosageHistory](#schemadosagehistory)|
|429|[Too Many Requests](https://tools.ietf.org/html/rfc6585#section-4)|The request has been rate limited.|[Error](#schemaerror)|
|default|Default|The request is invalid.|[Error](#schemaerror)|

### Response Headers

|Status|Header|Type|Format|Description|
|---|---|---|---|---|
|429|Retry-After|integer|int32|If the client should retry the request after a certain amount of time (in seconds), this header will be set. Often times, this will be set if the request is being rate limmmited.|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Import a CSV file of dosage history

<a id="opIdimportDosageHistory"></a>

`POST /dosage/import`

> Body parameter

```json
"dose,method,takenAt,takenOffAt,comment\n100,patch tw,2020-01-01T00:00:00Z,,\n100,patch tw,2020-01-02T12:00:00Z,,\n100,patch tw,2020-01-04T00:00:00Z,,\n100,patch tw,2020-01-05T12:00:00Z,,"
```

<h3 id="import-a-csv-file-of-dosage-history-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|Content-Type|header|string|true|The format to import the dosage history as.|
|body|body|[DosageHistoryCSV](#schemadosagehistorycsv)|true|none|

#### Enumerated Values

|Parameter|Value|
|---|---|
|Content-Type|text/csv|
|Content-Type|application/json|

> Example responses

> 200 Response

```json
{
  "records": 0,
  "succeeded": 0,
  "error": {
    "message": "string",
    "errors": [
      {}
    ],
    "details": null,
    "internal": true,
    "internalCode": "string"
  }
}
```

<h3 id="import-a-csv-file-of-dosage-history-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully imported the dosage history.|Inline|

<h3 id="import-a-csv-file-of-dosage-history-responseschema">Response Schema</h3>

Status Code **200**

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|» records|integer|true|none|The number of records in the file.|
|» succeeded|integer|true|none|The number of records actually imported successfully. This is not equal to #records if there were errors or duplicate entries.|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

<h1 id="e2clicker-service-notification">notification</h1>

## Get the server's push notification information

<a id="opIdwebPushInfo"></a>

`GET /pushinfo`

> Example responses

> 200 Response

```json
{
  "applicationServerKey": "string"
}
```

<h3 id="get-the-server's-push-notification-information-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully retrieved the server's push notification information.|[PushInfo](#schemapushinfo)|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<aside class="success">
This operation does not require authentication
</aside>

## Get the user's notification methods

<a id="opIduserNotificationMethods"></a>

`GET /notifications/methods`

> Example responses

> 200 Response

```json
{
  "webPush": [
    {
      "deviceID": "7996e974",
      "expirationTime": "2019-08-24T14:15:22Z",
      "keys": {
        "p256dh": "string"
      }
    }
  ]
}
```

<h3 id="get-the-user's-notification-methods-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully retrieved the user's push notification subscription.|[ReturnedNotificationMethods](#schemareturnednotificationmethods)|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Get the user's push notification subscription

<a id="opIduserPushSubscription"></a>

`GET /notifications/methods/push/{deviceID}`

<h3 id="get-the-user's-push-notification-subscription-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|deviceID|path|[PushDeviceID](#schemapushdeviceid)|true|The device ID of the push subscription to retrieve.|

> Example responses

> 200 Response

```json
{
  "deviceID": "7996e974",
  "endpoint": "string",
  "expirationTime": "2019-08-24T14:15:22Z",
  "keys": {
    "p256dh": "string",
    "auth": "string"
  }
}
```

<h3 id="get-the-user's-push-notification-subscription-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully retrieved the user's push notification subscription. The returned object will contain secrets.|[PushSubscription](#schemapushsubscription)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|The user does not have a push notification subscription for the specified device ID.|[Error](#schemaerror)|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Unsubscribe from push notifications

<a id="opIduserUnsubscribePush"></a>

`DELETE /notifications/methods/push/{deviceID}`

<h3 id="unsubscribe-from-push-notifications-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|deviceID|path|[PushDeviceID](#schemapushdeviceid)|true|The device ID of the push subscription to unsubscribe from.|

> Example responses

> default Response

```json
{
  "message": "string",
  "errors": [
    {
      "message": "string",
      "errors": [],
      "details": null,
      "internal": true,
      "internalCode": "string"
    }
  ],
  "details": null,
  "internal": true,
  "internalCode": "string"
}
```

<h3 id="unsubscribe-from-push-notifications-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Successfully unsubscribed from push notifications.|None|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Create or update a push subscription

<a id="opIduserSubscribePush"></a>

`PUT /notifications/methods/push`

> Body parameter

```json
{
  "deviceID": "7996e974",
  "endpoint": "string",
  "expirationTime": "2019-08-24T14:15:22Z",
  "keys": {
    "p256dh": "string",
    "auth": "string"
  }
}
```

<h3 id="create-or-update-a-push-subscription-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[PushSubscription](#schemapushsubscription)|true|none|

> Example responses

> default Response

```json
{
  "message": "string",
  "errors": [
    {
      "message": "string",
      "errors": [],
      "details": null,
      "internal": true,
      "internalCode": "string"
    }
  ],
  "details": null,
  "internal": true,
  "internalCode": "string"
}
```

<h3 id="create-or-update-a-push-subscription-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Successfully subscribed to push notifications.|None|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

<h1 id="e2clicker-service-ignore">ignore</h1>

## get___ignore_notification__haha_anything_can_go_here_lol

`GET /_ignore/notification/_haha_anything_can_go_here_lol`

> Example responses

> 500 Response

```json
{
  "type": "welcome_message",
  "message": null,
  "username": "string"
}
```

<h3 id="get___ignore_notification__haha_anything_can_go_here_lol-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|This endpoint is not used and should not be called.|Inline|

<h3 id="get___ignore_notification__haha_anything_can_go_here_lol-responseschema">Response Schema</h3>

Status Code **500**

*These fields aren't used for any API routes. They're just here to make oazapfts happy :) See https://github.com/oazapfts/oazapfts/issues/325.*

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|

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
|body|body|object|true|none|
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
|body|body|object|true|none|

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
  "hasAvatar": true,
  "secret": "string"
}
```

<h3 id="get-the-current-user-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully retrieved the current user.|Inline|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<h3 id="get-the-current-user-responseschema">Response Schema</h3>

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
  "errors": [
    {
      "message": "string",
      "errors": [],
      "details": null,
      "internal": true,
      "internalCode": "string"
    }
  ],
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
|body|body|string(binary)|true|none|

> Example responses

> default Response

```json
{
  "message": "string",
  "errors": [
    {
      "message": "string",
      "errors": [],
      "details": null,
      "internal": true,
      "internalCode": "string"
    }
  ],
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

<h3 id="delete-one-of-the-current-user's-sessions-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|query|integer(int64)|true|none|

> Example responses

> default Response

```json
{
  "message": "string",
  "errors": [
    {
      "message": "string",
      "errors": [],
      "details": null,
      "internal": true,
      "internalCode": "string"
    }
  ],
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
  "errors": [
    {
      "message": "string",
      "errors": [],
      "details": null,
      "internal": true,
      "internalCode": "string"
    }
  ],
  "details": null,
  "internal": true,
  "internalCode": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|message|string|true|none|A message describing the error|
|errors|[[Error](#schemaerror)]|false|none|An array of errors that caused this error. If this is populated, then [details] is omitted.|

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
  "name": "string",
  "description": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|true|none|A short string representing the delivery method. This is what goes into the DeliveryMethod fields.|
|units|string|true|none|The units of the delivery method.|
|name|string|true|none|The full name of the delivery method.|
|description|string|false|none|A description of the delivery method.|

<h2 id="tocS_Dosage">Dosage</h2>
<!-- backwards compatibility -->
<a id="schemadosage"></a>
<a id="schema_Dosage"></a>
<a id="tocSdosage"></a>
<a id="tocsdosage"></a>

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
[
  {
    "id": 0,
    "deliveryMethod": "string",
    "dose": 0.1,
    "takenAt": "2019-08-24T14:15:22Z",
    "takenOffAt": "2019-08-24T14:15:22Z",
    "comment": "string"
  }
]

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|[[DosageObservation](#schemadosageobservation)]|false|none|none|

<h2 id="tocS_DosageHistoryCSV">DosageHistoryCSV</h2>
<!-- backwards compatibility -->
<a id="schemadosagehistorycsv"></a>
<a id="schema_DosageHistoryCSV"></a>
<a id="tocSdosagehistorycsv"></a>
<a id="tocsdosagehistorycsv"></a>

```json
"dose,method,takenAt,takenOffAt,comment\n100,patch tw,2020-01-01T00:00:00Z,,\n100,patch tw,2020-01-02T12:00:00Z,,\n100,patch tw,2020-01-04T00:00:00Z,,\n100,patch tw,2020-01-05T12:00:00Z,,"

```

The CSV format of the user's dosage history.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|The CSV format of the user's dosage history.|

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
  "takenOffAt": "2019-08-24T14:15:22Z",
  "comment": "string"
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
|comment|string|false|none|A comment about the dosage, if any.|

<h2 id="tocS_Notification">Notification</h2>
<!-- backwards compatibility -->
<a id="schemanotification"></a>
<a id="schema_Notification"></a>
<a id="tocSnotification"></a>
<a id="tocsnotification"></a>

```json
{
  "type": "welcome_message",
  "message": null,
  "username": "string"
}

```

### Properties

*None*

<h2 id="tocS_NotificationType">NotificationType</h2>
<!-- backwards compatibility -->
<a id="schemanotificationtype"></a>
<a id="schema_NotificationType"></a>
<a id="tocSnotificationtype"></a>
<a id="tocsnotificationtype"></a>

```json
"welcome_message"

```

The type of notification:

  - `welcome_message` is sent to welcome the user. Realistically, it is
    used as a test message.
  - `reminder_message` is sent to remind the user of their hormone dose.
  - `account_notice_message` is sent to notify the user that they need
    to check their account.
  - `web_push_expiring_message` is sent to notify the user that their
    web push subscription is expiring.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|The type of notification:<br><br>  - `welcome_message` is sent to welcome the user. Realistically, it is<br>    used as a test message.<br>  - `reminder_message` is sent to remind the user of their hormone dose.<br>  - `account_notice_message` is sent to notify the user that they need<br>    to check their account.<br>  - `web_push_expiring_message` is sent to notify the user that their<br>    web push subscription is expiring.|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|welcome_message|
|*anonymous*|reminder_message|
|*anonymous*|account_notice_message|
|*anonymous*|web_push_expiring_message|

<h2 id="tocS_NotificationMessage">NotificationMessage</h2>
<!-- backwards compatibility -->
<a id="schemanotificationmessage"></a>
<a id="schema_NotificationMessage"></a>
<a id="tocSnotificationmessage"></a>
<a id="tocsnotificationmessage"></a>

```json
{
  "title": "string",
  "message": "string"
}

```

The message of the notification. This is derived from the notification type but can be overridden by the user.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|title|string|true|none|The title of the notification.|
|message|string|true|none|The message of the notification.|

<h2 id="tocS_PushDeviceID">PushDeviceID</h2>
<!-- backwards compatibility -->
<a id="schemapushdeviceid"></a>
<a id="schema_PushDeviceID"></a>
<a id="tocSpushdeviceid"></a>
<a id="tocspushdeviceid"></a>

```json
"7996e974"

```

A short ID associated with the device that the push subscription is for This is used to identify the device when updating its push subscription later on.
Realistically, this will be handled as an opaque random string generated on the device side, so the server has no way to correlate  it with any fingerprinting.
The recommended way to generate this string in JavaScript is:
```js crypto.randomUUID().slice(0, 8) ```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|A short ID associated with the device that the push subscription is for This is used to identify the device when updating its push subscription later on.<br>Realistically, this will be handled as an opaque random string generated on the device side, so the server has no way to correlate  it with any fingerprinting.<br>The recommended way to generate this string in JavaScript is:<br>```js crypto.randomUUID().slice(0, 8) ```|

<h2 id="tocS_PushInfo">PushInfo</h2>
<!-- backwards compatibility -->
<a id="schemapushinfo"></a>
<a id="schema_PushInfo"></a>
<a id="tocSpushinfo"></a>
<a id="tocspushinfo"></a>

```json
{
  "applicationServerKey": "string"
}

```

This is returned by the server and contains information that the client would need to subscribe to push notifications.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|applicationServerKey|string|true|none|A Base64-encoded string or ArrayBuffer containing an ECDSA P-256 public key that the push server will use to authenticate your application server. If specified, all messages from your application server must use the VAPID authentication scheme, and include a JWT signed with the corresponding private key. This key IS NOT the same ECDH key that you use to encrypt the data. For more information, see "Using VAPID with WebPush".|

<h2 id="tocS_PushSubscription">PushSubscription</h2>
<!-- backwards compatibility -->
<a id="schemapushsubscription"></a>
<a id="schema_PushSubscription"></a>
<a id="tocSpushsubscription"></a>
<a id="tocspushsubscription"></a>

```json
{
  "deviceID": "7996e974",
  "endpoint": "string",
  "expirationTime": "2019-08-24T14:15:22Z",
  "keys": {
    "p256dh": "string",
    "auth": "string"
  }
}

```

The configuration for a push notification subscription.
This is the object that is returned by calling PushSubscription.toJSON(). More information can be found at: https://developer.mozilla.org/en-US/docs/Web/API/PushSubscription/toJSON

### Properties

*None*

<h2 id="tocS_ReturnedNotificationMethods">ReturnedNotificationMethods</h2>
<!-- backwards compatibility -->
<a id="schemareturnednotificationmethods"></a>
<a id="schema_ReturnedNotificationMethods"></a>
<a id="tocSreturnednotificationmethods"></a>
<a id="tocsreturnednotificationmethods"></a>

```json
{
  "webPush": [
    {
      "deviceID": "7996e974",
      "expirationTime": "2019-08-24T14:15:22Z",
      "keys": {
        "p256dh": "string"
      }
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|webPush|[[ReturnedPushSubscription](#schemareturnedpushsubscription)]|false|none|[Similar to a [PushSubscription], but specifically for returning to the user. This type contains no secrets.]|

<h2 id="tocS_ReturnedPushSubscription">ReturnedPushSubscription</h2>
<!-- backwards compatibility -->
<a id="schemareturnedpushsubscription"></a>
<a id="schema_ReturnedPushSubscription"></a>
<a id="tocSreturnedpushsubscription"></a>
<a id="tocsreturnedpushsubscription"></a>

```json
{
  "deviceID": "7996e974",
  "expirationTime": "2019-08-24T14:15:22Z",
  "keys": {
    "p256dh": "string"
  }
}

```

Similar to a [PushSubscription], but specifically for returning to the user. This type contains no secrets.

### Properties

*None*

<h2 id="tocS_CustomNotifications">CustomNotifications</h2>
<!-- backwards compatibility -->
<a id="schemacustomnotifications"></a>
<a id="schema_CustomNotifications"></a>
<a id="tocScustomnotifications"></a>
<a id="tocscustomnotifications"></a>

```json
{
  "property1": {
    "title": "string",
    "message": "string"
  },
  "property2": {
    "title": "string",
    "message": "string"
  }
}

```

Custom notifications that the user can override with. The object keys are the notification types.

### Properties

*None*

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

