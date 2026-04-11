# Full-Stack ToDo Application

A modern, responsive Todo app built with **React** and **Tailwind CSS** on the frontend, and a powerful **Go (Fiber)** backend with **MongoDB** for data storage.

---

## Features

- Add, edit, delete tasks
- Mark tasks as completed
- Filter active and completed tasks
- Responsive UI (mobile & desktop)
- RESTful API with Go (Fiber)
- MongoDB integration for persistent storage

---

## Tech Stack

### Frontend:

- React
- Tailwind CSS
- Axios (for API calls)

### Backend:

- Go (Fiber framework)
- MongoDB
- Gorilla Mux (if applicable)
- MongoDB Go Driver

---

## Getting Started

### Frontend Setup

```bash
cd frontend
npm install
npm run dev
```

### Backend Setup

In backend directory

```bash
cd backend
air
```

### Set up MongoDB_URI

- Make sure MongoDB is running and the connection string is set in `backend/.env`

```bash
MONGODB_URI=mongodb+srv://<username>:<password>@cluster0.mongodb.net/todo_db
```

### AI Workflow Setup (Task From Text)

Add these values to `backend/.env` to enable the workflow endpoint:

```bash
OPENAI_API_KEY=<your_openai_api_key>
OPENAI_MODEL=gpt-4.1-mini
OPENAI_BASE_URL=https://api.openai.com
WORKFLOW_DEFAULT_USER_ID=local-user
WORKFLOW_DEFAULT_TIMEZONE=UTC
```

You can create a task from natural language using:

```http
POST /api/workflows/task-from-text
Content-Type: application/json

{
	"message": "Pay electricity bill tomorrow at 6 pm and remind me 1 hour before",
	"timezone": "Asia/Kolkata"
}
```

The workflow extracts task fields using OpenAI, creates a todo in MongoDB, and schedules a reminder record when requested.
