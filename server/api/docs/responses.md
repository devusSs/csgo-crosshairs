# JSON responses structures

## Generalised JSON responses

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

## Route specific responses

- to be implemented...
