# sociomile-angga

Sociomile is an omnichannel customer support platform designed with a service-oriented approach and scalable architecture.
Workspace for the Sociomile Project implementation.

## Projects

### Backend

Path: `sociomile-be`

Documentation:
- [Backend README](./sociomile-be/README.md)

Covers:
- Echo-based API
- MySQL persistence
- Redis rate limiting and cache
- Asynq worker for async event processing
- Postman collection and API documentation

### Frontend

Path: `sociomile-fe`

Documentation:
- [Frontend README](./sociomile-fe/README.md)

Covers:
- React + Vite client
- Login flow
- Conversation inbox UI
- Reply and escalation flow
- Ticket list and admin ticket update flow

## Suggested Run Order

1. Start infrastructure used by backend:
   - MySQL
   - Redis
2. Run database migration in `sociomile-be`
3. Start backend API in `sociomile-be`
4. Start async worker in `sociomile-be`
5. Start frontend in `sociomile-fe`

## Quick Links

- [Backend README](./sociomile-be/README.md)
- [Frontend README](./sociomile-fe/README.md)

## Directory Layout

```text
sociomile-angga/
├── README.md
├── sociomile-be/
└── sociomile-fe/
```

## User Interface Result

**Login Page** 
```text
username : angga@email.com
password : 123456
```

<img width="1842" height="1000" alt="Screenshot from 2026-03-09 01-39-58" src="https://github.com/user-attachments/assets/6e827db4-90c5-41ac-b22a-5e628b1271d1" />


**Dashboard Page**
<img width="1842" height="1000" alt="Screenshot from 2026-03-09 01-40-06" src="https://github.com/user-attachments/assets/10dc1c5c-b59d-4ee3-947f-da57ae3df2c8" />

