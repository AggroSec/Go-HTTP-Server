# Users

Endpoints for creating and updating user accounts.

---

## POST /api/users

Creates a new user account.

**Auth required:** No

### Request Body

```json
{
  "email": "user@example.com",
  "password": "yourpassword"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `email` | string | Yes | The user's email address |
| `password` | string | Yes | The user's password (stored as a bcrypt hash) |

### Response `201 Created`

```json
{
  "id": "e1b1b1b1-1b1b-1b1b-1b1b-1b1b1b1b1b1b",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

### Error Responses

| Status | Meaning |
|---|---|
| `400 Bad Request` | Missing or malformed email/password |

---

## PUT /api/users

Updates the email and/or password for the currently authenticated user.

**Auth required:** Yes — `Authorization: Bearer <access_token>`

### Request Body

```json
{
  "email": "newemail@example.com",
  "password": "newpassword"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `email` | string | Yes | New email address |
| `password` | string | Yes | New password |

### Response `200 OK`

```json
{
  "id": "e1b1b1b1-1b1b-1b1b-1b1b-1b1b1b1b1b1b",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-15T00:00:00Z",
  "email": "newemail@example.com",
  "is_chirpy_red": false
}
```

### Error Responses

| Status | Meaning |
|---|---|
| `400 Bad Request` | Missing or malformed request body |
| `401 Unauthorized` | Missing or invalid access token |
