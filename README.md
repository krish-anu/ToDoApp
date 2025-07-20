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
cd client
npm install
npm run dev
```

### Backend Setup

In root directory
```bash
air
```

### Set up MongoDB_URI 
- Make sure MongoDB is running and the connection string is set in .env
```bash
MONGODB_URI=mongodb+srv://<username>:<password>@cluster0.mongodb.net/todo_db
```
