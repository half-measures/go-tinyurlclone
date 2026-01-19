## Intro

Uses Secure Random String algo of base-62 as usual with a char set of a-z, A-Z, 0-9.
Expects collusions, so it will retry with a longer string if a collision is found.
56 Billion possible strings of length 10.

## DB

Uses MariaDB as the database.

## Configuration

- `BASE_URL`: The base URL for the shortened links (default: `http://localhost:8080`).

PORT=8081 BASE_URL=https://my-tinyurl.com go run main.go

## API

POST /shorten - Shortens a URL

## Frontend

Uses SvelteKit for the frontend.
