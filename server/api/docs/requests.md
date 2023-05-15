# JSON requests structures

## User routes

- GET /api/users/me (needs user auth)

  - does not need a JSON request body
  - only needs the JWT as "Authorization: Bearer <token>" header

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

## Crosshair routes

- POST /api/crosshairs (needs user auth)

  ```json
  {
    "code": "the encoded crosshair code"
  }
  ```

- every other route does not need a request body
- only needs the JWT as "Authorization: Bearer <token>" header

## Admin routes

- no route needs a request body
- only needs the JWT as "Authorization: Bearer <token>" header
- make sure the JWT of the user has admin role
