| Method | Route                            | Description                                             | Status | Auth needed |
| ------ | -------------------------------- | ------------------------------------------------------- | ------ | ----------- |
| GET    | /api/users/me                    | gets information about the logged in user               | ❌     | ❌          |
| POST   | /api/users/register              | registers a new user                                    | ❌     | ❌          |
| POST   | /api/users/login                 | login for existing users                                | ❌     | ❌          |
|        |                                  |                                                         |        |             |
| GET    | /api/crosshairs                  | gets all saved crosshairs from a specific user          | ❌     | ✅ (user)   |
| GET    | /api/crosshairs?code=            | gets a specific crosshair by it's code                  | ❌     | ✅ (admin)  |
| GET    | /api/crosshairs?date=            | gets crosshairs specified by a date range or single dat | ❌     | ✅ (admin)  |
| DELETE | /api/crosshairs                  | deletes all saved crosshairs from a specific user       | ❌     | ✅ (admin)  |
| DELETE | /api/crosshairs/:code            | deletes a specific crosshair by it's code               | ❌     | ✅ (admin)  |
|        |                                  |                                                         |        |             |
| GET    | /api/admins/users                | gets all users registered                               | ❌     | ✅ (admin)  |
| GET    | /api/admins/users?email=         | gets a user by their email                              | ❌     | ✅ (admin)  |
| GET    | /api/admins/crosshairs/          | gets all saved crosshairs                               | ❌     | ✅ (admin)  |
| GET    | /api/admins/crosshairs?code=     | gets a crosshair by it's code                           | ❌     | ✅ (admin)  |
| GET    | /api/admins/events?limit=        | gets X (limit) most recent events                       | ❌     | ✅ (admin)  |
| GET    | /api/admins/events`?type=&limit= | gets X (limit) most recent events by a specific type    | ❌     | ✅ (admin)  |
