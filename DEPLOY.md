# Deploy: GitHub Pages + Render

Architecture: React/Vite on GitHub Pages talks to a Golang/Gin API on Render.

```
GitHub  →  GitHub Pages  →  React + Vite
                                │ HTTPS
                                ▼
                             Render
                                │
                          Golang + Gin
                       /api/profile  /api/projects
```

Local first: from the repo root, run the API (`cd backend && go run ./cmd/server`) and the site (`cd frontend && npm install && npm run dev`). Default admin password is `changeme` unless `ADMIN_PASSWORD` is set. Open `http://localhost:5173`.

---

## 1. Render (API)

1. Push this repo to GitHub.
2. In [Render](https://render.com), create a **PostgreSQL** database (optional but recommended; the free web disk is ephemeral).
3. Create a **Web Service** from the same repo:
   - Root directory: `backend`
   - Runtime: Docker (uses `backend/Dockerfile`)
   - Or native Go: build `go build -o server ./cmd/server`, start `./server`
4. Set environment variables:

| Key | Value |
|---|---|
| `ADMIN_PASSWORD` | Password for `/#/admin` |
| `JWT_SECRET` | Long random string |
| `ENCRYPTION_KEY` | Long random string (encrypts LLM API keys) |
| `CORS_ORIGINS` | Your Pages origin, e.g. `https://<user>.github.io` |
| `DATABASE_URL` | Render Postgres URL (leave empty only for local SQLite) |
| `MAX_UPLOAD_MB` | `5` |

5. Deploy. Confirm `https://<service>.onrender.com/health` returns `{"ok":true}` and `/api/profile` returns seeded data from `porto.md`.
6. Copy the service URL. You will use it as `VITE_API_URL`.

On first boot the API parses `porto.md` into profile, skills, projects, and sections, and stores the full markdown as chatbot context.

If you do not attach Postgres, SQLite is used. On a free Render instance that file is lost on every deploy/restart and the seed runs again. Uploaded photos will also disappear. Use Postgres for a real site.

---

## 2. GitHub Pages (frontend)

1. In the GitHub repo: **Settings → Pages → Build and deployment → GitHub Actions**.
2. Add a repository secret:
   - `VITE_API_URL` = `https://<your-render-service>.onrender.com` (no trailing slash)
3. If this is a **project** site (`https://<user>.github.io/<repo>/`), add a repository variable:
   - `VITE_BASE` = `/<repo>/`
4. Push to `main` (or run the **Deploy GitHub Pages** workflow). The workflow in `.github/workflows/pages.yml` builds `frontend/` and publishes `dist`.
5. After the first deploy, add the Pages origin to Render `CORS_ORIGINS` if you have not already:
   - User site: `https://<user>.github.io`
   - Project site: `https://<user>.github.io` is enough (origin has no repo path)

Admin: `https://<user>.github.io/<repo>/#/admin`

---

## 3. After deploy

1. Open the admin panel and sign in with `ADMIN_PASSWORD`.
2. Upload a profile photo and project screenshots.
3. Open **llm**, paste keys for Gemini / GPT / Grok / Claude, pick a default model, save.
4. On the public site, use **Ask recruiter bot**. Switch the provider in the chat header. Off-topic questions are refused and redirected back to the portfolio.

Never commit API keys or `.env` files. Keys live in the admin UI (encrypted at rest) or in Render env vars for the app itself.
