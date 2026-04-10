# Webhooks

Endpoints for receiving incoming webhook events from external services.

---

## POST /api/polka/webhooks

Receives webhook events from the Polka payment service. Currently handles the `user.upgraded` event, which upgrades a user to **Chirpy Red** membership.

All other event types are acknowledged with `204 No Content` and ignored.

**Auth required:** Yes — API key via `Authorization: ApiKey <polka_api_key>`

> This endpoint uses API key authentication rather than JWT. The key must match the `POLKA_KEY` value set in the server's environment variables.

### Request Headers

```
Authorization: ApiKey <your_polka_api_key>
Content-Type: application/json
```

### Request Body

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "e1b1b1b1-1b1b-1b1b-1b1b-1b1b1b1b1b1b"
  }
}
```

| Field | Type | Description |
|---|---|---|
| `event` | string | The event type. Only `user.upgraded` triggers an action; all others are silently acknowledged. |
| `data.user_id` | UUID string | The ID of the user the event applies to |

### Response `204 No Content`

Returned on success, or when the event type is unrecognised. No response body.

### Error Responses

| Status | Meaning |
|---|---|
| `400 Bad Request` | Malformed request body or invalid `user_id` format |
| `401 Unauthorized` | Missing, malformed, or incorrect API key |
| `404 Not Found` | No user found with the given `user_id` (only returned for `user.upgraded` events) |
