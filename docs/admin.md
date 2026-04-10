# Admin

Endpoints for server monitoring and development utilities. The reset endpoint is restricted to the `dev` platform environment and should never be exposed in production.

---

## GET /api/healthz

Returns the health status of the server. Useful for checking that the server is up and reachable.

**Auth required:** No

### Response `200 OK`

```
Content-Type: text/plain; charset=utf-8

OK
```

---

## GET /admin/metrics

Returns an HTML page showing how many times the `/app` file server has been visited since the server started (or since the last reset).

**Auth required:** No

### Response `200 OK`

```
Content-Type: text/html; charset=utf-8
```

```html
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited 42 times!</p>
  </body>
</html>
```

---

## POST /admin/reset

Resets the file server hit counter to zero and clears all users from the database.

**Auth required:** No

> ⚠️ **This endpoint is restricted to the `dev` environment.** It will return `403 Forbidden` if the server is not running with `PLATFORM=dev`. It exists purely for development and testing purposes and should never be exposed in a production deployment.

### Response `200 OK`

```
Content-Type: text/plain; charset=utf-8

Metrics reset
```

### Error Responses

| Status | Meaning |
|---|---|
| `403 Forbidden` | Server is not running in `dev` platform mode |
| `500 Internal Server Error` | Database error during reset |
