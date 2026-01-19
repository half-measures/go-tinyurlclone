## Intro

Uses Secure Random String algo of base-62 as usual with a char set of a-z, A-Z, 0-9.
Expects collusions, so it will retry with a longer string if a collision is found.
56 Billion possible strings of length 10.

## DB

Uses MariaDB as the database.

## API

POST /shorten - Shortens a URL

## Frontend

Uses SvelteKit for the frontend.
