# Responses from admin routes

Every response will be send as the `data` part of the generalised success response.

## Get all registered users

- URL: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;/api/admins/users
- Method: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;GET
- Response body:

```json
{
  "users": [
    {
      "id": "uid",
      "created_at": "2023-05-18-19:40:13",
      "updated_at": "2023-05-18-19:40:13",
      "e_mail": "",
      "role": "",
      "verified_mail": true,
      "register_ip": "",
      "login_ip": "",
      "last_login": "2023-05-18-19:40:13",
      "crosshairs_registered": 1,
      "avatar_url": ""
    },
    {}
  ]
}
```

## Get one user by their e-mail

- URL: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;/api/admins/users?email=
- Method: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;GET
- Response body:

```json
{
  "id": "uid",
  "created_at": "2023-05-18-19:40:13",
  "updated_at": "2023-05-18-19:40:13",
  "e_mail": "",
  "role": "",
  "verified_mail": true,
  "register_ip": "",
  "login_ip": "",
  "last_login": "2023-05-18-19:40:13",
  "crosshairs_registered": 1,
  "avatar_url": ""
}
```

## Get all crosshairs

- URL: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;/api/admins/crosshairs
- Method: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;GET
- Response body:

```json
{
  "crosshairs": [
    {
      "id": "uid",
      "added": "2023-05-18-19:40:13",
      "code": "",
      "note": ""
    },
    {}
  ]
}
```

## Get all crosshairs from a user by their e-mail

- URL: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;/api/admins/crosshairs?email=
- Method: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;GET
- Response body:

```json
{
  "crosshair": [
    {
      "id": "uid",
      "added": "2023-05-18-19:40:13",
      "code": "",
      "note": ""
    },
    {}
  ]
}
```

## Get all events

- URL: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;/api/admins/events
- Method: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;GET
- Response body:

```json
{
    [
        {
          "id": "uid",
	        "type": "",
	        "data": {
              "url": "",
	            "method": "",
	            "issuer": "",
	            "data": {}
            },
	        "timestamp": "2023-05-18-19:40:13"
        },
        {}
    ]
}
```

## Get all events with a limit

- URL: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;/api/admins/events?limit=
- Method: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;GET
- Response body:

```json
{
    [
        {
          "id": "uid",
	        "type": "",
	        "data": {
              "url": "",
	            "method": "",
	            "issuer": "",
	            "data": {}
            },
	        "timestamp": "2023-05-18-19:40:13"
        },
        {}
    ]
}
```

## Get all events by event type

- URL: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;/api/admins/events?type=
- Method: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;GET
- Response body:

```json
{
    [
        {
          "id": "uid",
	        "type": "",
	        "data": {
              "url": "",
	            "method": "",
	            "issuer": "",
	            "data": {}
            },
	        "timestamp": "2023-05-18-19:40:13"
        },
        {}
    ]
}
```

## Get all events by event type with a limit

- URL: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;/api/admins/events?type=&limit=
- Method: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;GET
- Response body:

```json
{
    [
        {
          "id": "uid",
	        "type": "",
	        "data": {
              "url": "",
	            "method": "",
	            "issuer": "",
	            "data": {}
            },
	        "timestamp": "2023-05-18-19:40:13"
        },
        {}
    ]
}
```

## Get total (overall) API stats

- URL: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;/api/admins/stats/total
- Method: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;GET
- Response body:

```json
{
  "registered_users": 0,
  "registered_crosshairs": 0
}
```

## Get daily API stats

- URL: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;/api/admins/stats/daily
- Method: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;GET
- Response body:

```json
{
  "users_registered": 0,
  "user_logins": 0,
  "api_requests": 0
}
```

## Get all logs sorted by date (default)

- URL: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;/api/admins/logs
- Method: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;GET
- Response body:

```json
{
  "code": 200,
  "data": {
    "logs": [
      {
        "level": "info",
        "ts": 1687635086.142514,
        "msg": "/api/admins/logs",
        "status": 200,
        "method": "GET",
        "path": "/api/admins/logs",
        "query": "",
        "ip": "172.29.0.1",
        "user-agent": "PostmanRuntime/7.32.3",
        "latency": 0.000875542,
        "time": "2023-06-24T19:31:26Z"
      },
      {}
    ]
  }
}
```
