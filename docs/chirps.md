# Chirps

Endpoints for creating, retrieving, and deleting chirps.

A chirp is a short message (max 140 characters) posted by an authenticated user. Certain words are automatically filtered and replaced with `****` — specifically: `kerfuffle`, `sharbert`, and `fornax`.

---

## POST /api/chirps

Creates a new chirp for the authenticated user.

**Auth required:** Yes — `Authorization: Bearer <access_token>`

### Request Body

```json
{
  "body": "Hello, Chirpy world!"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `body` | string | Yes | The content of the chirp. Max 140 characters. Profanity-filtered automatically. |

### Response `201 Created`

```json
{
  "id": "a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z",
  "body": "Hello, Chirpy world!",
  "user_id": "e1b1b1b1-1b1b-1b1b-1b1b-1b1b1b1b1b1b"
}
```

### Error Responses

| Status | Meaning |
|---|---|
| `400 Bad Request` | Missing Authorization header or malformed body |
| `401 Unauthorized` | Invalid or expired access token |
| `400 Bad Request` | Chirp body exceeds 140 characters |

---

## GET /api/chirps

Returns a list of all chirps. Supports optional filtering by author and sorting by creation date.

**Auth required:** No

### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `author_id` | UUID string | No | Filter chirps to only those posted by this user ID |
| `sort` | string | No | Sort order by creation date. Accepts `asc` or `desc`. Defaults to unsorted if omitted. |

### Example Requests

```
GET /api/chirps
GET /api/chirps?sort=desc
GET /api/chirps?author_id=e1b1b1b1-1b1b-1b1b-1b1b-1b1b1b1b1b1b
GET /api/chirps?author_id=e1b1b1b1-1b1b-1b1b-1b1b-1b1b1b1b1b1b&sort=asc
```

### Response `200 OK`

```json
[
  {
    "id": "a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1",
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z",
    "body": "Hello, Chirpy world!",
    "user_id": "e1b1b1b1-1b1b-1b1b-1b1b-1b1b1b1b1b1b"
  }
]
```

Returns an empty array `[]` if no chirps exist.

### Error Responses

| Status | Meaning |
|---|---|
| `400 Bad Request` | Invalid `sort` parameter value or malformed `author_id` |
| `500 Internal Server Error` | Database error |

---

## GET /api/chirps/{chirpID}

Returns a single chirp by its ID.

**Auth required:** No

### Path Parameters

| Parameter | Type | Description |
|---|---|---|
| `chirpID` | UUID string | The ID of the chirp to retrieve |

### Response `200 OK`

```json
{
  "id": "a1a1a1a1-a1a1-a1a1-a1a1-a1a1a1a1a1a1",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z",
  "body": "Hello, Chirpy world!",
  "user_id": "e1b1b1b1-1b1b-1b1b-1b1b-1b1b1b1b1b1b"
}
```

### Error Responses

| Status | Meaning |
|---|---|
| `404 Not Found` | No chirp found with the given ID |

---

## DELETE /api/chirps/{chirpID}

Deletes a chirp. Users may only delete their own chirps.

**Auth required:** Yes — `Authorization: Bearer <access_token>`

### Path Parameters

| Parameter | Type | Description |
|---|---|---|
| `chirpID` | UUID string | The ID of the chirp to delete |

### Response `204 No Content`

No response body.

### Error Responses

| Status | Meaning |
|---|---|
| `401 Unauthorized` | Missing or invalid access token |
| `403 Forbidden` | The authenticated user does not own this chirp |
| `404 Not Found` | No chirp found with the given ID |
| `500 Internal Server Error` | Database error during deletion |
