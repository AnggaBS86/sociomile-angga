# sociomile-fe

React frontend for the Sociomile backend.

## Features

- Login with backend JWT auth
- Conversation list and detail view
- Agent reply flow
- Agent escalation form
- Ticket list
- Admin ticket status update
- Vite dev proxy to `http://127.0.0.1:8080`

## Run

```bash
cd "/media/user/New Volume/go/sociomile-angga/sociomile-fe"
npm install
npm run dev
```

## Backend requirement

Run the backend first at `http://127.0.0.1:8080`.

## Optional env

Create `.env` in this project if you want to bypass the Vite proxy:

```bash
VITE_API_BASE_URL=http://127.0.0.1:8080
```

## Seed login

```text
email: angga@email.com
password: 123456
```
