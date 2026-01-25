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

Once started, access the UI at [http://localhost:3000](http://localhost:3000) for direct Svelte access, or [http://localhost](http://localhost) to test through the Nginx production proxy.

## Production Deployment

For production (e.g., AWS EC2), use the Nginx proxy and set your domain:

1. Copy `.env.example` to `.env` and set your domain:
   ```bash
   cp .env.example .env
   # Edit .env and set DOMAIN=http://yourdomain.com
   ```
2. Update `nginx.conf` if using HTTPS/SSL.
3. Start the stack:
   ```bash
   docker compose up --build
   ```

The app will be accessible on port 80 via Nginx, which routes traffic to the frontend and backend.

## Configuration

Environment variables can be configured in `docker-compose.yml`:

- `BASE_URL`: The base URL returned by the backend for shortened links.
- `PUBLIC_API_URL`: The URL the frontend uses to contact the backend.

## API

- `POST /shorten`: Shortens a URL.
- `GET /:slug`: Redirects to the original URL.
