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
        "created_at": 2023-05-18-19:40:13,
	    "updated_at": 2023-05-18-19:40:13,
	    "e_mail": "",
	    "role": "",
	    "verified_mail": true,
	    "register_ip": "",
	    "login_ip": "",
	    "last_login": 2023-05-18-19:40:13,
	    "crosshairs_registered": 1,
       },
       ...
    ]
}
```

## Get one user by their e-mail

- URL: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;/api/admins/users?email=
- Method: &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;GET
- Response body:

```json
{
    "created_at": 2023-05-18-19:40:13,
	"updated_at": 2023-05-18-19:40:13,
	"e_mail": "",
	"role": "",
	"verified_mail": true,
	"register_ip": "",
	"login_ip": "",
	"last_login": 2023-05-18-19:40:13,
	"crosshairs_registered": 1,
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
            "added": 2023-05-18-19:40:13,
	        "code": "",
	        "note": ""
        },
        ...
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
            "added": 2023-05-18-19:40:13,
	        "code": "",
	        "note": ""
        },
        ...
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
            "id": "",
	        "type": "",
	        "data": {
                "url": "",
	            "method": "",
	            "issuer": "",
	            "data": {}
            },
	        "timestamp": 2023-05-18-19:40:13
        },
        ...
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
            "id": "",
	        "type": "",
	        "data": {
                "url": "",
	            "method": "",
	            "issuer": "",
	            "data": {}
            },
	        "timestamp": 2023-05-18-19:40:13
        },
        ...
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
            "id": "",
	        "type": "",
	        "data": {
                "url": "",
	            "method": "",
	            "issuer": "",
	            "data": {}
            },
	        "timestamp": 2023-05-18-19:40:13
        },
        ...
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
            "id": "",
	        "type": "",
	        "data": {
                "url": "",
	            "method": "",
	            "issuer": "",
	            "data": {}
            },
	        "timestamp": 2023-05-18-19:40:13
        },
        ...
    ]
}
```