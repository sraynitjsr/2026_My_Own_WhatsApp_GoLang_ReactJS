# WhatsApp Clone - Real-time Messaging Application (Not An Android or Mobile App)

A full-featured WhatsApp-like messaging application built with **GoLang** (backend) and **ReactJS** (frontend).

## 🚀 Features

### Core Features (Implemented)
- ✅ **User Authentication** - Secure signup/login with JWT
- ✅ **Real-time Messaging** - WebSocket-based instant messaging
- ✅ **Direct Conversations** - One-on-one chats
- ✅ **Group Chats** - Create and manage group conversations
- ✅ **File/Image Sharing** - Upload and share media files
- ✅ **Online/Offline Status** - See who's online
- ✅ **Read Receipts** - Track message read status
- ✅ **Message History** - Persistent message storage
- ✅ **User Search** - Find users to chat with

### Advanced Features (To Be Implemented)
- ⏳ **Video/Voice Calling** - WebRTC integration
- ⏳ **Message Encryption** - End-to-end encryption
- ⏳ **Typing Indicators** - See when someone is typing
- ⏳ **Push Notifications** - Real-time notifications
- ⏳ **Message Reactions** - React to messages with emojis
- ⏳ **Voice Messages** - Record and send audio
- ⏳ **Stories/Status** - Share temporary status updates

## 🏗️ Tech Stack

### Backend
- **Go 1.21+** - Main programming language
- **Gin** - Web framework
- **GORM** - ORM for database operations
- **PostgreSQL** - Primary database
- **WebSocket (Gorilla)** - Real-time communication
- **JWT** - Authentication
- **bcrypt** - Password hashing

### Frontend
- **React 18** - UI library
- **React Router** - Navigation
- **Zustand** - State management
- **Axios** - HTTP client
- **WebSocket** - Real-time updates
- **Tailwind CSS** - Styling

### Infrastructure
- **Docker** - Containerization
- **Docker Compose** - Multi-container orchestration

## 📁 Project Structure

```
.
├── backend/
│   ├── config/          # Database and app configuration
│   ├── controllers/     # Request handlers
│   ├── middleware/      # Authentication, CORS, etc.
│   ├── models/          # Database models
│   ├── routes/          # API routes
│   ├── websocket/       # WebSocket logic
│   ├── main.go          # Application entry point
│   ├── go.mod           # Go dependencies
│   └── .env.example     # Environment variables template
│
├── frontend/
│   ├── public/          # Static files
│   ├── src/
│   │   ├── api/         # API client
│   │   ├── components/  # Reusable components
│   │   ├── pages/       # Page components
│   │   ├── services/    # WebSocket and other services
│   │   ├── store/       # State management
│   │   └── App.js       # Main app component
│   └── package.json     # Node dependencies
│
├── docker-compose.yml   # Docker orchestration
├── Makefile            # Build and run commands
└── README.md           # This file
```

## 🚦 Getting Started

### Prerequisites
- **Go 1.21+** installed
- **Node.js 18+** and npm installed
- **PostgreSQL** installed (or use Docker)
- **Docker & Docker Compose** (optional, for containerized setup)

### Option 1: Run with Docker (Recommended for Beginners)

1. **Clone the repository** (if using git):
   ```bash
   cd /Users/sray/.njdk/2026_My_Own_WhatsApp_GoLang_ReactJS
   ```

2. **Start all services**:
   ```bash
   make docker-up
   ```
   This will:
   - Start PostgreSQL database
   - Start backend server on http://localhost:8080
   - Start frontend on http://localhost:3000

3. **View logs**:
   ```bash
   make docker-logs
   ```

4. **Stop services**:
   ```bash
   make docker-down
   ```

### Option 2: Run Locally (For Development)

#### Step 1: Setup Database
Install and start PostgreSQL, then create the database:
```sql
CREATE DATABASE whatsapp_clone;
```

#### Step 2: Setup Backend
```bash
cd backend

# Create environment file
cp .env.example .env
# Edit .env with your database credentials

# Install dependencies
go mod download

# Run the server
go run main.go
```

Backend will run on http://localhost:8080

#### Step 3: Setup Frontend
Open a new terminal:
```bash
cd frontend

# Create environment file
cp .env.example .env

# Install dependencies
npm install

# Start development server
npm start
```

Frontend will run on http://localhost:3000

## 📚 API Documentation

### Authentication Endpoints

#### Register
```
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "password123",
  "display_name": "John Doe"
}
```

#### Login
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

### User Endpoints (Require Authentication)

All authenticated endpoints require the `Authorization` header:
```
Authorization: Bearer <your_jwt_token>
```

#### Get Current User
```
GET /api/v1/users/me
```

#### Search Users
```
GET /api/v1/users/search?q=john
```

### Conversation Endpoints

#### Get All Conversations
```
GET /api/v1/conversations
```

#### Create Conversation
```
POST /api/v1/conversations
Content-Type: application/json

{
  "type": "direct",  // or "group"
  "name": "My Group",  // for group chats
  "member_ids": [2, 3, 4]
}
```

### Message Endpoints

#### Get Messages
```
GET /api/v1/messages/conversation/:conversationId
```

#### Send Message
```
POST /api/v1/messages
Content-Type: application/json

{
  "conversation_id": 1,
  "content": "Hello!",
  "message_type": "text"
}
```

### WebSocket Connection
```
WS /api/v1/ws?token=<your_jwt_token>
```

## 🗺️ Development Roadmap

### Phase 1: Foundation ✅ (COMPLETED)
- [x] Project setup and structure
- [x] Database schema design
- [x] User authentication (signup/login)
- [x] Basic REST API
- [x] Frontend routing and auth flow

### Phase 2: Core Messaging ✅ (COMPLETED)
- [x] WebSocket implementation
- [x] Real-time message delivery
- [x] Conversation management
- [x] Message history
- [x] User interface for chat

### Phase 3: Enhanced Features 🔄 (NEXT)
- [ ] File upload and storage (S3/local)
- [ ] Image preview and optimization
- [ ] Group chat management UI
- [ ] User profile settings
- [ ] Message search functionality

### Phase 4: Advanced Features
- [ ] Typing indicators
- [ ] Message reactions (emojis)
- [ ] Voice messages
- [ ] Video/Voice calling (WebRTC)
- [ ] End-to-end encryption
- [ ] Push notifications

### Phase 5: Polish & Production
- [ ] Error handling improvements
- [ ] Loading states and skeletons
- [ ] Responsive design for mobile
- [ ] Performance optimization
- [ ] Security audit
- [ ] Deployment setup (AWS/GCP/Heroku)

## 🛠️ Development Tips

### Backend Development

1. **Adding a New Model**:
   - Create model in `backend/models/`
   - Add to auto-migrate in `config/database.go`

2. **Adding a New API Endpoint**:
   - Create controller in `backend/controllers/`
   - Add route in `backend/routes/routes.go`

3. **Testing the API**:
   Use curl or Postman:
   ```bash
   # Health check
   curl http://localhost:8080/health

   # Register
   curl -X POST http://localhost:8080/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{"username":"test","email":"test@test.com","password":"test123"}'
   ```

### Frontend Development

1. **Adding a New Component**:
   - Create component in `frontend/src/components/`
   - Import and use in pages

2. **State Management**:
   - Global state: Use Zustand stores in `src/store/`
   - Local state: Use React useState

3. **Making API Calls**:
   ```javascript
   import { messageAPI } from '../api';
   
   const messages = await messageAPI.getMessages(conversationId);
   ```

## 🐛 Troubleshooting

### Backend Issues

**Database connection failed**:
- Check PostgreSQL is running: `pg_isctl status`
- Verify credentials in `.env`
- Ensure database exists: `psql -l`

**Port already in use**:
- Change `PORT` in backend `.env`
- Or kill process: `lsof -ti:8080 | xargs kill`

### Frontend Issues

**Module not found**:
- Run `npm install` in frontend directory

**WebSocket connection failed**:
- Ensure backend is running
- Check `REACT_APP_WS_URL` in frontend `.env`

**CORS errors**:
- Verify `FRONTEND_URL` in backend `.env` matches your frontend URL

## 📖 Learning Resources

### For Beginners

**Go Learning**:
- [Go Tour](https://go.dev/tour/) - Official interactive tutorial
- [Go by Example](https://gobyexample.com/) - Practical examples

**React Learning**:
- [React Docs](https://react.dev/) - Official documentation
- [React Tutorial](https://react.dev/learn) - Step-by-step guide

**WebSocket**:
- [WebSocket API](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)
- [Gorilla WebSocket](https://github.com/gorilla/websocket)

**PostgreSQL**:
- [PostgreSQL Tutorial](https://www.postgresqltutorial.com/)

### Recommended Next Steps

1. **Week 1-2**: Get the basic app running, understand the code structure
2. **Week 3-4**: Add file upload functionality
3. **Week 5-6**: Implement typing indicators and online status
4. **Week 7-8**: Add group chat management UI
5. **Week 9-10**: Work on voice/video calling
6. **Week 11-12**: Polish UI, add encryption, deploy

## 🤝 Contributing

Since you're building this project to learn:
1. Try to implement features yourself first
2. Read error messages carefully
3. Use console.log() and print statements for debugging
4. Break down large features into small steps
5. Test frequently as you build

## 📝 License

This project is for educational purposes. Feel free to use and modify as needed.

## 🙋 Need Help?

- **Backend errors**: Check `backend/` terminal for Go errors
- **Frontend errors**: Check browser console (F12)
- **Database issues**: Check PostgreSQL logs
- **Docker issues**: Run `docker-compose logs -f`

---

## 🎯 Quick Start Commands

```bash
# Setup environment files
make setup-env

# Run with Docker
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down

# Run locally (in separate terminals)
make run-backend
make run-frontend
```

**Happy Coding! 🚀**
