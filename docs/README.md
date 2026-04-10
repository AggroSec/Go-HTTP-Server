# Chirpy API Documentation

Base URL: `http://localhost:8080`

All request and response bodies use JSON. Endpoints that require authentication expect a JWT access token passed as a Bearer token in the `Authorization` header unless otherwise noted.

## Sections

- [Users](./users.md) — Register and manage user accounts
- [Authentication](./authentication.md) — Login, token refresh, and token revocation
- [Chirps](./chirps.md) — Create, retrieve, and delete chirps
- [Admin](./admin.md) — Server metrics and dev tooling
- [Webhooks](./webhooks.md) — Incoming webhook events from Polka

## Authentication Overview

Chirpy uses a two-token auth system:

- **Access token** — A short-lived JWT (1 hour) passed as `Authorization: Bearer <token>` on protected requests
- **Refresh token** — A long-lived token used to obtain a new access token without re-entering credentials. Pass it the same way via `Authorization: Bearer <refresh_token>` when calling `POST /api/refresh`

Tokens are issued on login. See [Authentication](./authentication.md) for full details.
