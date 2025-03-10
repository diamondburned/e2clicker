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

`GET /delivery-methods`

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
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully recorded dosage.|[Dose](#schemadose)|

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
|doseTimes|query|array[string]|true|none|

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

## Update a dosage in the user's history

<a id="opIdeditDose"></a>

`PUT /dosage/dose/{doseTime}`

> Body parameter

```json
{
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
|doseTime|path|string(date-time)|true|none|
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
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Successfully updated the dosage.|None|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Delete a dosage from the user's history

<a id="opIdforgetDose"></a>

`DELETE /dosage/dose/{doseTime}`

This operation is broken in the backend due to a parsing error and
should not be used. Instead, prefer using [forgetDoses].

<h3 id="delete-a-dosage-from-the-user's-history-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|doseTime|path|string(date-time)|true|none|

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

<h3 id="delete-a-dosage-from-the-user's-history-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Successfully deleted the dosage.|None|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Export the user's dosage history

<a id="opIdexportDoses"></a>

`GET /dosage/export-doses`

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
|200|Content-Disposition|string||none|
|429|Retry-After|integer|int32|If the client should retry the request after a certain amount of time (in seconds), this header will be set. Often times, this will be set if the request is being rate limmmited.|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Import a CSV file of dosage history

<a id="opIdimportDoses"></a>

`POST /dosage/import-doses`

> Body parameter

```json
"deliveryMethod,dose,takenAt,takenOffAt,comment\npatch tw,100,2020-01-01T00:00:00Z,,\npatch tw,100,2020-01-02T12:00:00Z,,\npatch tw,100,2020-01-04T00:00:00Z,,\npatch tw,100,2020-01-05T12:00:00Z,,"
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

`GET /push-info`

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

## Get the server's supported notification methods

<a id="opIdsupportedNotificationMethods"></a>

`GET /notifications/methods`

> Example responses

> 200 Response

```json
[
  "webPush"
]
```

<h3 id="get-the-server's-supported-notification-methods-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully retrieved the server's notification methods.|[NotificationMethodSupports](#schemanotificationmethodsupports)|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<aside class="success">
This operation does not require authentication
</aside>

## Send a test notification

<a id="opIdsendTestNotification"></a>

`POST /notifications/test`

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

<h3 id="send-a-test-notification-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Successfully sent the test notification.|None|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Get the user's notification preferences

<a id="opIduserNotificationPreferences"></a>

`GET /notifications/preferences`

> Example responses

> 200 Response

```json
{
  "notificationConfigs": {
    "webPush": [
      {
        "deviceID": "7996e974",
        "endpoint": "string",
        "expirationTime": "2019-08-24T14:15:22Z",
        "keys": {
          "p256dh": "string",
          "auth": "string"
        }
      }
    ],
    "email": [
      {
        "address": "user@example.com",
        "name": "string"
      }
    ]
  },
  "customNotifications": {
    "property1": {
      "title": "string",
      "message": "string"
    },
    "property2": {
      "title": "string",
      "message": "string"
    }
  }
}
```

<h3 id="get-the-user's-notification-preferences-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|Successfully retrieved the user's notification preferences.|[NotificationPreferences](#schemanotificationpreferences)|
|default|Default|The request is invalid.|[Error](#schemaerror)|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
bearerAuth
</aside>

## Update the user's notification preferences

<a id="opIduserUpdateNotificationPreferences"></a>

`PUT /notifications/preferences`

> Body parameter

```json
null
```

<h3 id="update-the-user's-notification-preferences-parameters">Parameters</h3>

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

<h3 id="update-the-user's-notification-preferences-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|204|[No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5)|Successfully updated the user's notification methods.|None|
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
|*anonymous*|[[Dose](#schemadose)]|false|none|[A dose of medication in time.]|

<h2 id="tocS_DosageHistoryCSV">DosageHistoryCSV</h2>
<!-- backwards compatibility -->
<a id="schemadosagehistorycsv"></a>
<a id="schema_DosageHistoryCSV"></a>
<a id="tocSdosagehistorycsv"></a>
<a id="tocsdosagehistorycsv"></a>

```json
"deliveryMethod,dose,takenAt,takenOffAt,comment\npatch tw,100,2020-01-01T00:00:00Z,,\npatch tw,100,2020-01-02T12:00:00Z,,\npatch tw,100,2020-01-04T00:00:00Z,,\npatch tw,100,2020-01-05T12:00:00Z,,"

```

The CSV format of the user's dosage history.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|The CSV format of the user's dosage history.|

<h2 id="tocS_Dose">Dose</h2>
<!-- backwards compatibility -->
<a id="schemadose"></a>
<a id="schema_Dose"></a>
<a id="tocSdose"></a>
<a id="tocsdose"></a>

```json
{
  "deliveryMethod": "string",
  "dose": 0.1,
  "takenAt": "2019-08-24T14:15:22Z",
  "takenOffAt": "2019-08-24T14:15:22Z",
  "comment": "string"
}

```

A dose of medication in time.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
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
  - `test_message` is sent to test your notification settings.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|The type of notification:<br><br>  - `welcome_message` is sent to welcome the user. Realistically, it is<br>    used as a test message.<br>  - `reminder_message` is sent to remind the user of their hormone dose.<br>  - `account_notice_message` is sent to notify the user that they need<br>    to check their account.<br>  - `web_push_expiring_message` is sent to notify the user that their<br>    web push subscription is expiring.<br>  - `test_message` is sent to test your notification settings.|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|welcome_message|
|*anonymous*|reminder_message|
|*anonymous*|account_notice_message|
|*anonymous*|web_push_expiring_message|
|*anonymous*|test_message|

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

<h2 id="tocS_EmailSubscription">EmailSubscription</h2>
<!-- backwards compatibility -->
<a id="schemaemailsubscription"></a>
<a id="schema_EmailSubscription"></a>
<a id="tocSemailsubscription"></a>
<a id="tocsemailsubscription"></a>

```json
{
  "address": "user@example.com",
  "name": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|address|string(email)|true|none|The email address to send the notification to. This email address will appear in the `To` field of the email.|
|name|string|false|none|The name of the user to send the email to. This name will be used with the email address in the `To` field.|

<h2 id="tocS_NotificationPreferences">NotificationPreferences</h2>
<!-- backwards compatibility -->
<a id="schemanotificationpreferences"></a>
<a id="schema_NotificationPreferences"></a>
<a id="tocSnotificationpreferences"></a>
<a id="tocsnotificationpreferences"></a>

```json
{
  "notificationConfigs": {
    "webPush": [
      {
        "deviceID": "7996e974",
        "endpoint": "string",
        "expirationTime": "2019-08-24T14:15:22Z",
        "keys": {
          "p256dh": "string",
          "auth": "string"
        }
      }
    ],
    "email": [
      {
        "address": "user@example.com",
        "name": "string"
      }
    ]
  },
  "customNotifications": {
    "property1": {
      "title": "string",
      "message": "string"
    },
    "property2": {
      "title": "string",
      "message": "string"
    }
  }
}

```

The user's notification preferences.
Each key is a notification type and the value is the notification configuration for that type. It may be nil if the server does not support a particular notification type.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|notificationConfigs|object|true|none|none|
|» webPush|[[PushSubscription](#schemapushsubscription)]|false|none|[The configuration for a push notification subscription.<br>This is the object that is returned by calling PushSubscription.toJSON(). More information can be found at: https://developer.mozilla.org/en-US/docs/Web/API/PushSubscription/toJSON]|

<h2 id="tocS_NotificationMethodSupports">NotificationMethodSupports</h2>
<!-- backwards compatibility -->
<a id="schemanotificationmethodsupports"></a>
<a id="schema_NotificationMethodSupports"></a>
<a id="tocSnotificationmethodsupports"></a>
<a id="tocsnotificationmethodsupports"></a>

```json
[
  "webPush"
]

```

A list of notification methods that the server supports.

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
  "locale": "string"
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

