# Authentication

Endpoints for logging in and managing access and refresh tokens.

---

## POST /api/login

Authenticates a user and returns a JWT access token and a refresh token.

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
| `email` | string | Yes | The user's registered email |
| `password` | string | Yes | The user's password |

### Response `200 OK`

```json
{
  "id": "e1b1b1b1-1b1b-1b1b-1b1b-1b1b1b1b1b1b",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false,
  "token": "<jwt_access_token>",
  "refresh_token": "<refresh_token>"
}
```

The `token` is a JWT valid for **1 hour**. The `refresh_token` is a long-lived token that can be used to obtain a new access token without logging in again.

### Error Responses

| Status | Meaning |
|---|---|
| `400 Bad Request` | Malformed request body |
| `401 Unauthorized` | Incorrect email or password |

---

## POST /api/refresh

Issues a new JWT access token using a valid, non-expired, non-revoked refresh token.

**Auth required:** Yes — `Authorization: Bearer <refresh_token>`

### Request Body

None.

### Response `200 OK`

```json
{
  "token": "<new_jwt_access_token>"
}
```

### Error Responses

| Status | Meaning |
|---|---|
| `400 Bad Request` | Malformed or missing Authorization header |
| `401 Unauthorized` | Refresh token is invalid, expired, or has been revoked |

---

## POST /api/revoke

Revokes a refresh token, preventing it from being used to issue new access tokens. Use this on logout.

**Auth required:** Yes — `Authorization: Bearer <refresh_token>`

### Request Body

None.

### Response `204 No Content`

No response body.

### Error Responses

| Status | Meaning |
|---|---|
| `400 Bad Request` | Malformed or missing Authorization header |
| `500 Internal Server Error` | Database error during revocation |
