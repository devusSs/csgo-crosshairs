# Routes, requests and response structure

## Routes

| Method | Route                              | Description                                             | Status | Auth needed                                   |
| ------ | ---------------------------------- | ------------------------------------------------------- | ------ | --------------------------------------------- |
| GET    | /api/users/me                      | gets information about the logged in user               | ✅     | ✅ (user)                                     |
| POST   | /api/users/register                | registers a new user                                    | ✅     | ❌                                            |
| POST   | /api/users/register?action=resend  | resends verify email                                    | ✅     | ❌                                            |
| GET    | /api/users/verifyMail?code=        | verified the user's email on registration               | ✅     | ❌                                            |
| POST   | /api/users/login                   | login for existing users                                | ✅     | ❌                                            |
| GET    | /api/users/logout                  | logout for logged in user                               | ✅     | ✅ (user)                                     |
| POST   | /api/users/resetPass               | reset password for registered user                      | ✅     | ❌                                            |
| GET    | /api/users/resetPass?email=&code=  | check reset password code from email                    | ✅     | ❌                                            |
| PATCH  | /api/users/resetPass?email=&code=  | performs the actual password reset                      | ✅     | ❌                                            |
| PATCH  | /api/users/newPass                 | performs password reset for logged in user              | ✅     | ✅ (user)                                     |
| POST   | /api/users/avatar                  | updates the user avatar                                 | ✅     | ✅ (user)                                     |
| DELETE | /api/users/avatar                  | deletes the current avatar of a user                    | ✅     | ✅ (user)                                     |
| GET    | /api/integration/twitch/login      | makes Twitch integration possible for user              | ✅     | ✅ (user)                                     |
| GET    | /api/integration/twitch/disconnect | removes Twitch integration for user                     | ✅     | ✅ (user)                                     |
|        |                                    |                                                         |        |
| GET    | /api/crosshairs                    | gets all saved crosshairs from a specific user          | ✅     | ✅ (user)                                     |
| GET    | /api/crosshairs?code=              | gets a specific crosshair by it's code                  | ✅     | ✅ (user)                                     |
| GET    | /api/crosshairs?start=&end=        | gets crosshairs specified by a date range or single dat | ✅     | ✅ (user)                                     |
| POST   | /api/crosshairs/add                | saves a new crosshair from a specific user              | ✅     | ✅ (user)                                     |
| DELETE | /api/crosshairs                    | deletes all saved crosshairs from a specific user       | ✅     | ✅ (user)                                     |
| DELETE | /api/crosshairs?code=              | deletes a specific crosshair by it's code               | ✅     | ✅ (user)                                     |
|        |                                    |                                                         |        |                                               |
| GET    | /api/admins/users                  | gets all users registered                               | ✅     | ✅ (admin)                                    |
| GET    | /api/admins/users?email=           | gets a user by their email                              | ✅     | ✅ (admin)                                    |
| GET    | /api/admins/crosshairs             | gets all saved crosshairs                               | ✅     | ✅ (admin)                                    |
| GET    | /api/admins/crosshairs?email=      | gets all saved crosshairs from a specific user          | ✅     | ✅ (admin)                                    |
| GET    | /api/admins/logs                   | gets all logs sorted by timestamp                       | ✅     | ✅ (admin)                                    |
| GET    | /api/admins/events                 | gets all events                                         | ✅     | ✅ (admin)                                    |
| GET    | /api/admins/events?limit=          | gets X (limit) most recent events                       | ✅     | ✅ (admin)                                    |
| GET    | /api/admins/events?type=           | gets all events by a specific type                      | ✅     | ✅ (admin)                                    |
| GET    | /api/admins/events?type=&limit=    | gets X (limit) most recent events by a specific type    | ✅     | ✅ (admin)                                    |
| GET    | /api/admins/stats/total            | gets overall API stats                                  | ✅     | ✅ (admin)                                    |
| GET    | /api/admins/stats/daily            | gets API stats for day                                  | ✅     | ✅ (admin)                                    |
| GET    | /api/admins/stats/system           | gets system stats for API host                          | ✅     | ✅ (engineer, see below for more information) |

Note:

Regarding engineer authorization:

```bash
Authorization: Bearer <token>
```

The token will ONLY be accessable via Postgres commands on the actual backend server.<br/>

Events supported so far:

- "user_registered"
- "user_password_change"
- "user_uploaded_avatar"

## Response structure

### Generalised JSON responses

- every route returns the same base (!) JSON structure

### For success responses

Code is subject to change depending on the result of the query.

```json
{
  "code": 200,
  "data": {}
}
```

### For fail responses

Code is subject to change depending on the result of the query.

```json
{
  "code": 400,
  "error": {
    "error_code": "e.g. not_found",
    "error_message": "This is a formal written description of the error code."
  }
}
```

### Route specific responses

Please head to the [response docs](responses) and choose the corresponding route type.

## Request structure

Every route and subroute takes in different request bodies or may not need a request body at all.<br/>

Please head to the [request docs](requests) and choose the corresponding route type.<br/>

Note: The admin routes do not need a request body hence there is no documentation on them.
