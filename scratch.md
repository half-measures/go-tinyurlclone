1. Hardcoded api call in +page.svelte - API hardcoded to localhost:8080/shorten, needs to be a relative path
2. both BASE_URL and MARIADB_URI are hardcoded in main.go and have local defaults.
3. CORS is broad, should be restricted to frontend domain
4. Dockerization - not actually done in the app
