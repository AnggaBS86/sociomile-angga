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
cd "/sociomile-angga/sociomile-fe"
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

## User Interface Result

**Login Page** 
<img width="1842" height="1000" alt="Screenshot from 2026-03-09 01-39-58" src="https://github.com/user-attachments/assets/6e827db4-90c5-41ac-b22a-5e628b1271d1" />


**Dashboard Page**
<img width="1842" height="1000" alt="Screenshot from 2026-03-09 01-40-06" src="https://github.com/user-attachments/assets/10dc1c5c-b59d-4ee3-947f-da57ae3df2c8" />

