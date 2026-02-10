## Intro

Uses Secure Random String algo of base-62 as usual with a char set of a-z, A-Z, 0-9.
Expects collisions, so it will retry with a longer string if a collision is found.
56 Billion possible strings of length 10.

## Running with Docker

The entire stack can be started with a single command:

```bash
docker compose up --build
```

This will spin up:

- **MariaDB**: Database (Port 3306)
- **Go Backend**: API (Port 8080)
- **Svelte Frontend**: UI (Port 3000)

Once started, access the UI at [http://localhost](http://localhost) to test through the Nginx production proxy.

## Production Deployment

For production (e.g., AWS EC2 with domain `t.michaeldev.org`):

1. **Environment Config**: Create a `.env` file in the root directory:
   ```env
   DOMAIN=http://t.michaeldev.org
   ```
2. **Nginx**: The included `nginx.conf` is pre-configured to handle any incoming hostname (`server_name _`).
3. **Start**:
   ```bash
   docker compose down
   docker compose up -d --build
   ```

The app will be accessible on port 80. Nginx handles the routing for both the frontend and the `/shorten` API.

## Troubleshooting

- **Backend Unreachable**: Ensure your EC2 Security Group allows inbound traffic on port 80.
- **localhost links**: Ensure the `DOMAIN` variable is set in `.env` before running `docker compose up`.

## Configuration

Environment variables can be configured in `docker-compose.yml`:

- `BASE_URL`: The base URL returned by the backend for shortened links.
- `PUBLIC_API_URL`: The URL the frontend uses to contact the backend.

## API

- `POST /shorten`: Shortens a URL.
- `GET /:slug`: Redirects to the original URL.

## Once Deployed...

Any changes to svelte needs a rebuild, as it just reuses what it has.
docker compose up -d --build will force that
can also just do docker compose up -d --build frontend
