<script lang="ts">
  import { PUBLIC_API_URL } from "$env/static/public";
  let url = $state("");
  let isLoading = $state(false);
  let shortUrl = $state("");
  let error = $state("");

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (!url) return;
    isLoading = true;
    shortUrl = "";
    error = "";

    try {
      const response = await fetch(`${PUBLIC_API_URL}/shorten`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ long_url: url }),
      });

      const data = await response.json();
      if (data.short_url) {
        shortUrl = data.short_url;
      } else if (data.error) {
        error = data.error;
      }
    } catch (err) {
      error = "Failed to connect to backend.";
      console.error(err);
    } finally {
      isLoading = false;
    }
  }
</script>

<main class="hero-section">
  <div class="container">
    <div class="content">
      <h1 class="font-heading main-title">
        Shorten Your <span class="gradient-text">Links</span>,
        <br />
        Amplify Your <span class="gradient-text">Reach</span>
      </h1>
      <p class="subtitle">
        A minimalist and high-speed URL shortener built with SvelteKit and Go.
      </p>

      <form class="glass input-card transition-all" onsubmit={handleSubmit}>
        <div class="input-wrapper">
          <input
            type="url"
            placeholder="Paste your long link here..."
            bind:value={url}
            class="url-input"
            required
          />
          <button
            type="submit"
            class="submit-btn transition-all"
            disabled={isLoading}
          >
            {#if isLoading}
              <div class="loader"></div>
            {:else}
              Shorten
            {/if}
          </button>
        </div>
      </form>

      {#if shortUrl}
        <div class="glass result-card transition-all">
          <p>Your shortened link is ready!</p>
          <div class="short-url-wrapper">
            <a href={shortUrl} target="_blank" class="short-url">{shortUrl}</a>
          </div>
        </div>
      {/if}

      {#if error}
        <div class="error-message">
          {error}
        </div>
      {/if}

      <div class="features">
        <div class="feature-tag">✨ Lightning Fast</div>
        <div class="feature-tag">🔒 End-to-End Secure</div>
        <div class="feature-tag">📊 It Just Works</div>
      </div>
    </div>
  </div>
</main>

<style>
  .hero-section {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    padding: 2rem;
    text-align: center;
  }

  .container {
    width: 100%;
    max-width: 800px;
  }

  .main-title {
    font-size: 4rem;
    font-weight: 800;
    line-height: 1.1;
    margin-bottom: 1.5rem;
    color: var(--text-main);
  }

  .gradient-text {
    background: linear-gradient(135deg, #60a5fa 0%, #a855f7 100%);
    background-clip: text;
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }

  .subtitle {
    font-size: 1.25rem;
    color: var(--text-muted);
    margin-bottom: 3rem;
    max-width: 600px;
    margin-left: auto;
    margin-right: auto;
  }

  .input-card {
    padding: 0.75rem;
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
    margin-bottom: 2rem;
  }

  .result-card {
    padding: 1.5rem;
    margin-bottom: 2rem;
    text-align: center;
  }

  .short-url-wrapper {
    margin-top: 1rem;
    padding: 1rem;
    background: rgba(255, 255, 255, 0.05);
    border-radius: 12px;
  }

  .short-url {
    font-size: 1.25rem;
    font-weight: 700;
    color: #60a5fa;
    text-decoration: none;
  }

  .error-message {
    color: #ef4444;
    margin-bottom: 2rem;
  }

  .input-wrapper {
    display: flex;
    gap: 0.5rem;
  }

  .url-input {
    flex: 1;
    background: transparent;
    border: none;
    padding: 1rem 1.5rem;
    font-size: 1.1rem;
    color: var(--text-main);
    outline: none;
    width: 100%;
  }

  .url-input::placeholder {
    color: rgba(255, 255, 255, 0.3);
  }

  .submit-btn {
    background: var(--text-main);
    color: var(--bg-color);
    border: none;
    padding: 1rem 2rem;
    border-radius: 16px;
    font-weight: 700;
    font-size: 1rem;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    min-width: 120px;
  }

  .submit-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 10px 20px -5px rgba(255, 255, 255, 0.2);
  }

  .submit-btn:active {
    transform: translateY(0);
  }

  .submit-btn:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }

  .features {
    display: flex;
    justify-content: center;
    gap: 1.5rem;
    flex-wrap: wrap;
  }

  .feature-tag {
    font-size: 0.875rem;
    color: var(--text-muted);
    background: rgba(255, 255, 255, 0.05);
    padding: 0.5rem 1rem;
    border-radius: 99px;
    border: 1px solid rgba(255, 255, 255, 0.05);
  }

  .loader {
    width: 20px;
    height: 20px;
    border: 3px solid rgba(0, 0, 0, 0.1);
    border-top: 3px solid var(--bg-color);
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }

  @media (max-width: 640px) {
    .main-title {
      font-size: 2.5rem;
    }
    .input-wrapper {
      flex-direction: column;
    }
    .submit-btn {
      width: 100%;
    }
  }
</style>
