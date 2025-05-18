# 🎬 MyTop100Movies

A full-stack web app where users can create and manage their personal **Top 100 Favorite Movies** list. Powered by Go (backend), PostgreSQL (database), and Next.js (frontend). Integrates with the TMDb API for fetching movie metadata and posters.

---

## 🚀 Features

- 🔐 User-specific movie lists
- 📥 Add a movie to your Top 100 (with rank + notes)
- 🔁 Update rank or notes
- ❌ Remove a movie from your list
- 🧹 **Clear** your entire Top 100 list with one click
- 📦 Persistent storage using PostgreSQL
- 🎨 Responsive UI built with Next.js

- ---

## 🛠️ Tech Stack

### Backend
- **Language:** Go
- **Database:** PostgreSQL
- **Router:** net/http
- **Data Format:** JSON (REST API)
- **3rd Party API:** [TMDb API](https://developer.themoviedb.org/)

 ### Frontend
- **Framework:** Next.js (App Router with Client Components)
- **Styling:** TailwindCSS
- **Image Handling:** Next/Image

---

## ⚙️ How to Run Locally

### 🧩 Prerequisites

- Go 1.20+
- PostgreSQL
- Node.js 18+
- TMDb API Key (optional for fetching metadata)

---

### 🗃️ Database Setup

```sql
CREATE TABLE movies (
  id SERIAL PRIMARY KEY,
  tmdb_id INT UNIQUE,
  title TEXT,
  overview TEXT,
  poster_url TEXT
);

CREATE TABLE user_movies (
  id SERIAL PRIMARY KEY,
  user_id INT,
  movie_id INT REFERENCES movies(id),
  rank INT,
  notes TEXT
);
```

🔙 Backend Setup

```cd backend
go mod tidy
go run main.go
```
> Runs on localhost:8080

> Example endpoint: GET /top100/get?user_id=1


🌐 Frontend Setup

```cd frontend
npm install
npm run dev
```
> Runs on localhost:3000

> Visit /mytoplist to view the Top 100 page

📡 API Endpoints

| Method | Endpoint         | Description                     |
| ------ | ---------------- | ------------------------------- |
| GET    | `/top100/get`    | Fetch top 100 list for user     |
| POST   | `/top100/add`    | Add a new movie to list         |
| PUT    | `/top100/update` | Update rank or notes            |
| DELETE | `/top100/remove` | Delete a specific movie         |
| DELETE | `/top100/clear`  | Clear entire Top 100 for a user |


🖼️ Screenshots
![image](https://github.com/user-attachments/assets/f645622b-ae48-4d45-b4a8-43ebfe6e08f9)

![image](https://github.com/user-attachments/assets/06b09e9f-7b81-4959-9717-fd3130181234)

👨‍💻 Developer Notes
Basic validation included on both backend and frontend.

TMDb API calls are handled once and stored — avoids hitting rate limits.

Designed to be modular for easy future upgrades (e.g., auth, drag-and-drop reorder, user login).

💡 Future Features
✅ Swap movie ranks

⏳ User authentication

⏳ Drag-and-drop sorting

⏳ Movie search using TMDb

⏳ Public sharing of lists





