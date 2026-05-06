# 🚀 Getting Started Guide - For Complete Beginners

Welcome! This guide will walk you through setting up and running your WhatsApp clone app, even if you're new to programming.

## 📋 What You'll Need

### Required Software

1. **Go (GoLang)** - Backend programming language
   - Download: https://go.dev/dl/
   - Choose the installer for macOS
   - After installation, verify: Open Terminal and type `go version`

2. **Node.js** - JavaScript runtime for frontend
   - Download: https://nodejs.org/ (choose LTS version)
   - After installation, verify: `node --version` and `npm --version`

3. **PostgreSQL** - Database
   - Option A (Easier): Use Docker (see below)
   - Option B: Download from https://www.postgresql.org/download/macosx/
   
4. **Docker** (Recommended - makes everything easier!)
   - Download: https://www.docker.com/products/docker-desktop
   - Install Docker Desktop for Mac
   - After installation, verify: `docker --version`

5. **Code Editor**
   - VS Code (Recommended): https://code.visualstudio.com/
   - Comes with great extensions for Go and React

## 🎯 Step-by-Step Setup

### Option 1: The Easy Way (Using Docker) ⭐ RECOMMENDED

This is perfect for beginners because Docker handles everything for you!

#### Step 1: Open Terminal
- Press `Cmd + Space`, type "Terminal", press Enter

#### Step 2: Navigate to Project
```bash
cd /Users/sray/.njdk/2026_My_Own_WhatsApp_GoLang_ReactJS
```

#### Step 3: Create Environment Files
```bash
# Copy the example environment files
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
```

#### Step 4: Start Everything with One Command!
```bash
docker-compose up -d
```

This will:
- Download and setup PostgreSQL database
- Build and start your backend server
- Build and start your frontend app

Wait about 2-3 minutes for everything to start...

#### Step 5: Open Your App
- Open your web browser
- Go to: http://localhost:3000
- You should see the WhatsApp Clone login page! 🎉

#### Step 6: Create Your First Account
1. Click "Register here"
2. Fill in:
   - Username: yourname
   - Email: your@email.com
   - Display Name: Your Name
   - Password: password123
3. Click Register

You're in! Now you can start chatting.

#### Useful Docker Commands
```bash
# See what's running
docker-compose ps

# View logs (helpful for debugging)
docker-compose logs -f

# Stop everything
docker-compose down

# Restart everything
docker-compose restart
```

---

### Option 2: Running Locally (For Learning More)

This option helps you understand each piece of the application.

#### Part 1: Setup PostgreSQL Database

1. **Install PostgreSQL** (if not using Docker)
   ```bash
   # On macOS with Homebrew
   brew install postgresql@15
   brew services start postgresql@15
   ```

2. **Create the Database**
   ```bash
   # Open PostgreSQL
   psql postgres
   
   # In psql, create database (copy and paste this line)
   CREATE DATABASE whatsapp_clone;
   
   # Quit psql
   \q
   ```

#### Part 2: Setup Backend

1. **Open a new Terminal window**

2. **Navigate to backend folder**
   ```bash
   cd /Users/sray/.njdk/2026_My_Own_WhatsApp_GoLang_ReactJS/backend
   ```

3. **Create environment file**
   ```bash
   cp .env.example .env
   ```

4. **Edit the .env file** (open with any text editor)
   - Change database password if needed
   - Keep other settings as they are

5. **Install Go dependencies**
   ```bash
   go mod download
   ```
   This downloads all required packages. Takes about 1-2 minutes.

6. **Run the backend server**
   ```bash
   go run main.go
   ```
   
   You should see:
   ```
   Database connected and migrated successfully
   Server starting on port 8080
   ```

   ✅ Backend is running! Keep this terminal window open.

#### Part 3: Setup Frontend

1. **Open ANOTHER new Terminal window** (yes, you need 2 terminals!)

2. **Navigate to frontend folder**
   ```bash
   cd /Users/sray/.njdk/2026_My_Own_WhatsApp_GoLang_ReactJS/frontend
   ```

3. **Create environment file**
   ```bash
   cp .env.example .env
   ```

4. **Install Node dependencies**
   ```bash
   npm install
   ```
   This downloads all required packages. Takes about 3-5 minutes.
   You might see some warnings - that's okay!

5. **Start the frontend**
   ```bash
   npm start
   ```
   
   After a minute, your browser should automatically open to http://localhost:3000
   
   ✅ Frontend is running! Keep this terminal window open too.

Now you have:
- Terminal 1: Backend server running on port 8080
- Terminal 2: Frontend app running on port 3000
- PostgreSQL database running in the background

## 🎮 Using Your App

### Creating Your First Account

1. You'll see the login page
2. Click "Register here"
3. Fill in the form:
   ```
   Username: john
   Email: john@test.com
   Display Name: John Doe
   Password: test123
   ```
4. Click "Register"

### Creating a Second User (to test chatting)

To really test your app, you need two users:

1. In a **private/incognito browser window**, go to http://localhost:3000
2. Register a second account:
   ```
   Username: jane
   Email: jane@test.com
   Display Name: Jane Smith
   Password: test123
   ```

Now you can chat between John and Jane!

### Starting a Conversation

Currently, the basic version requires some manual steps. We'll make this easier later:

1. You can use the API directly to create a conversation
2. Or modify the UI to add a "New Chat" button (a great learning exercise!)

## 🐛 Common Issues and Solutions

### "Port 8080 is already in use"

Something else is using that port.

**Solution**:
```bash
# Find what's using the port
lsof -ti:8080

# Kill it
lsof -ti:8080 | xargs kill -9

# Or change the port in backend/.env
PORT=8081
```

### "Database connection failed"

PostgreSQL isn't running or credentials are wrong.

**Solution**:
```bash
# Check if PostgreSQL is running
brew services list

# Start it if needed
brew services start postgresql@15

# Verify credentials in backend/.env match your setup
```

### "Module not found" errors in frontend

Dependencies didn't install properly.

**Solution**:
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### "WebSocket connection failed"

Backend isn't running or frontend can't reach it.

**Solution**:
1. Check backend terminal - is it running?
2. Check backend/.env has correct PORT
3. Check frontend/.env has correct REACT_APP_API_URL

### Docker "Container already in use"

Previous containers are still running.

**Solution**:
```bash
docker-compose down
docker-compose up -d
```

## 📚 Understanding the Code (Beginner's Guide)

### Backend Structure (Go)

```
backend/
├── main.go              ← Starts the server
├── config/
│   └── database.go      ← Connects to database
├── models/              ← Database tables (User, Message, etc.)
├── controllers/         ← Handle requests (like login, send message)
├── routes/              ← Define URLs (/api/v1/auth/login)
├── middleware/          ← Security stuff (authentication)
└── websocket/           ← Real-time messaging magic
```

**Example - What happens when you login:**
1. You submit email + password from browser
2. Request goes to `routes/routes.go` → finds `/auth/login`
3. Calls `controllers/auth.go` → `Login()` function
4. Checks password, creates JWT token
5. Sends token back to browser

### Frontend Structure (React)

```
frontend/src/
├── App.js               ← Main app entry
├── pages/               ← Full page components
│   ├── Login.js         ← Login page
│   ├── Register.js      ← Signup page
│   └── Chat.js          ← Main chat interface
├── components/          ← Reusable UI pieces
│   ├── ConversationList.js  ← Left sidebar with chats
│   ├── ChatWindow.js        ← Message display area
│   └── MessageInput.js      ← Text box to type messages
├── store/               ← App state (what user is logged in, etc.)
└── api/                 ← Functions to talk to backend
```

**Example - What happens when you send a message:**
1. You type in `MessageInput.js` and press Send
2. Calls API function in `api/index.js`
3. Sends HTTP POST to backend
4. Backend saves to database
5. Backend broadcasts via WebSocket
6. Other users receive message instantly!

## 🎓 Learning Path

### Week 1: Get Comfortable
- [ ] Get the app running (you're here!)
- [ ] Create test accounts
- [ ] Send messages between users
- [ ] Explore the code - open files, read comments
- [ ] Make a small change: Update the app title in frontend

### Week 2: Simple Modifications
- [ ] Change color scheme (edit Tailwind classes)
- [ ] Add timestamp formatting to messages
- [ ] Add a "New Chat" button
- [ ] Modify user profile fields

### Week 3: Add Features
- [ ] Implement user search in UI
- [ ] Add conversation creation from frontend
- [ ] Show online/offline indicators
- [ ] Add loading spinners

### Week 4: More Complex Features
- [ ] File upload button and preview
- [ ] Image messages display
- [ ] Typing indicators
- [ ] Message reactions

## 🔍 Debugging Tips

### Reading Error Messages

**Backend (Go) errors appear in Terminal 1:**
```
Error: pq: password authentication failed for user "postgres"
```
This tells you the database password is wrong.

**Frontend (React) errors appear in:**
1. Terminal 2 (build errors)
2. Browser Console (press F12) (runtime errors)

### Using Console Logs

Add debug prints:

**In Go (backend):**
```go
log.Println("User ID:", userID)
log.Println("Message content:", message.Content)
```

**In React (frontend):**
```javascript
console.log('User:', user)
console.log('Messages:', messages)
```

### Testing the API Directly

Use `curl` to test backend without frontend:

```bash
# Test health check
curl http://localhost:8080/health

# Register a user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@test.com",
    "password": "password123"
  }'
```

## 📞 Next Steps

Now that you're set up:

1. **Explore the Code**
   - Open the project in VS Code
   - Read through the files
   - Try to understand what each part does

2. **Make Your First Change**
   - Try changing the app title
   - Modify a color
   - Add your own welcome message

3. **Follow the Roadmap**
   - Check README.md for the development roadmap
   - Pick a feature to implement
   - Start small, build up gradually

4. **Learn as You Go**
   - Don't worry if you don't understand everything
   - Google error messages
   - Read Go and React documentation
   - Experiment and break things (that's how you learn!)

## 💡 Pro Tips

1. **Always keep both terminal windows open** when running locally
2. **Save your work frequently** - use Git if you know how
3. **Test after each change** - don't make 10 changes at once
4. **Read error messages carefully** - they often tell you exactly what's wrong
5. **Google is your friend** - someone else has had your error before
6. **Take breaks** - your brain learns better when not exhausted

## 🎉 You're Ready!

You now have a working WhatsApp clone! This is just the beginning. As you learn more about Go and React, you can add amazing features like video calls, encryption, and more.

**Remember**: Every expert was once a beginner. Take it one step at a time, and don't be afraid to experiment!

Happy coding! 🚀

---

**Need More Help?**

- Go documentation: https://go.dev/doc/
- React documentation: https://react.dev/
- PostgreSQL tutorial: https://www.postgresqltutorial.com/
- Docker getting started: https://docs.docker.com/get-started/
