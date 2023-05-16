# Routes, requests and response structure

## Routes

| Method | Route                           | Description                                             | Status | Auth needed |
| ------ | ------------------------------- | ------------------------------------------------------- | ------ | ----------- |
| GET    | /api/users/me                   | gets information about the logged in user               | ❌     | ✅ (user)   |
| POST   | /api/users/register             | registers a new user                                    | ❌     | ❌          |
| GET    | /api/users/verifyMail/:code     | verified the user's email on registration               | ❌     | ❌          |
| POST   | /api/users/login                | login for existing users                                | ❌     | ❌          |
| GET    | /api/users/logout               | logout for logged in user                               | ❌     | ✅ (user)   |
|        |                                 |                                                         |        |             |
| GET    | /api/crosshairs                 | gets all saved crosshairs from a specific user          | ❌     | ✅ (user)   |
| GET    | /api/crosshairs?code=           | gets a specific crosshair by it's code                  | ❌     | ✅ (user)   |
| GET    | /api/crosshairs?date=           | gets crosshairs specified by a date range or single dat | ❌     | ✅ (user)   |
| POST   | /api/crosshairs                 | saves a new crosshair from a specific user              | ❌     | ✅ (user)   |
| DELETE | /api/crosshairs                 | deletes all saved crosshairs from a specific user       | ❌     | ✅ (user)   |
| DELETE | /api/crosshairs/:code           | deletes a specific crosshair by it's code               | ❌     | ✅ (user)   |
|        |                                 |                                                         |        |             |
| GET    | /api/admins/users               | gets all users registered                               | ❌     | ✅ (admin)  |
| GET    | /api/admins/users?email=        | gets a user by their email                              | ❌     | ✅ (admin)  |
| GET    | /api/admins/crosshairs/         | gets all saved crosshairs                               | ❌     | ✅ (admin)  |
| GET    | /api/admins/crosshairs?code=    | gets a crosshair by it's code                           | ❌     | ✅ (admin)  |
| GET    | /api/admins/events?limit=       | gets X (limit) most recent events                       | ❌     | ✅ (admin)  |
| GET    | /api/admins/events?type=&limit= | gets X (limit) most recent events by a specific type    | ❌     | ✅ (admin)  |

## Requests structure

### User routes

- GET /api/users/me (needs user auth)

  - does not need a JSON request body
  - authorization will be fetched from sessions cookie

- POST /api/users/register

  ```json
  {
    "e_mail": "",
    "password": "",
    "admin_token": "blank_or_valid_token"
  }
  ```

- POST /api/users/login

  ```json
  {
    "e_mail": "",
    "password": ""
  }
  ```

### Crosshair routes

- POST /api/crosshairs (needs user auth)

  ```json
  {
    "code": "the encoded crosshair code"
  }
  ```

- every other route does not need a request body
- authorization will be fetched from sessions cookie

### Admin routes

- no route needs a request body
- authorization will be fetched from sessions cookie
- make sure the user has admin role

## Response structure

### Generalised JSON responses

- every route returns the same base (!) JSON structure

### For success responses

```json
{
    "code": 20X,
    "data": {}
}
```

### For fail responses

```json
{
    "code": 40X / 500,
    "error": {
        "error_code": "e.g. not_found",
        "error_message": "This is a formal written description of the error code."
    }
}
```

### Route specific responses

- to be implemented...
