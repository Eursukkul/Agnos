# API Spec

## `POST /staff/create`

Create hospital staff account.

Request:
```json
{
  "username": "string",
  "password": "string",
  "hospital": "string"
}
```

Response `201`:
```json
{
  "id": 1,
  "username": "alice",
  "hospital": "hospital-a"
}
```

## `POST /staff/login`

Login staff and receive JWT token.

Request:
```json
{
  "username": "string",
  "password": "string",
  "hospital": "string"
}
```

Response `200`:
```json
{
  "token": "<jwt>"
}
```

## `POST /patient/search`

Search patients in the same hospital as authenticated staff.

Auth header:
```text
Authorization: Bearer <jwt>
```

Request (all optional):
```json
{
  "national_id": "string",
  "passport_id": "string",
  "first_name": "string",
  "middle_name": "string",
  "last_name": "string",
  "date_of_birth": "YYYY-MM-DD",
  "phone_number": "string",
  "email": "string"
}
```

Response `200`:
```json
{
  "patients": [
    {
      "id": 1,
      "hospital": "hospital-a",
      "first_name_th": "สมชาย",
      "last_name_th": "ใจดี",
      "first_name_en": "Somchai",
      "last_name_en": "Jaidee",
      "date_of_birth": "1990-01-01T00:00:00Z",
      "patient_hn": "HN001",
      "national_id": "1234567890123",
      "passport_id": "AA123456",
      "phone_number": "0812345678",
      "email": "x@example.com",
      "gender": "M"
    }
  ]
}
```

Error codes:
- `400`: invalid body
- `401`: missing/invalid token or login failure
- `500`: internal search failure
