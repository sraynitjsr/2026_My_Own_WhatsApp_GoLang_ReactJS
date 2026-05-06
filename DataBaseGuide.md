# 🗄️ Database Connection Guide

## 📋 Database Connection Details

Your PostgreSQL database is running in Docker with these credentials:

```
Host: localhost
Port: 5432
Database: whatsapp_clone
Username: postgres
Password: postgres123
```

---

## 🚀 Method 1: Command Line (Easiest & Fastest)

### Connect Using Docker Exec

```bash
# Connect to PostgreSQL inside the Docker container
docker exec -it whatsapp_db psql -U postgres -d whatsapp_clone
```

Once connected, you'll see:
```
whatsapp_clone=#
```

### Basic SQL Commands

```sql
-- List all tables
\dt

-- Describe a table structure
\d users
\d messages
\d conversations

-- View table with more details
\d+ users

-- List all databases
\l

-- Quit
\q
```

---

## 📊 Method 2: Run Queries from Terminal (Quick Queries)

### Check if database exists
```bash
docker exec whatsapp_db psql -U postgres -c "\l"
```

### View all users
```bash
docker exec whatsapp_db psql -U postgres -d whatsapp_clone -c "SELECT * FROM users;"
```

### Count messages
```bash
docker exec whatsapp_db psql -U postgres -d whatsapp_clone -c "SELECT COUNT(*) FROM messages;"
```

### See recent messages
```bash
docker exec whatsapp_db psql -U postgres -d whatsapp_clone -c "SELECT id, sender_id, content, created_at FROM messages ORDER BY created_at DESC LIMIT 10;"
```

---

## 💻 Method 3: GUI Tools (Visual Interface)

### Option A: pgAdmin (Free, Popular)

**Download:** https://www.pgadmin.org/download/

### Password Can't Be Hardcoded in Production Code, FYI....
**Setup:**
1. Install pgAdmin
2. Right-click "Servers" → Create → Server
3. Enter:
   - Name: `WhatsApp Clone`
   - Host: `localhost`
   - Port: `5432`
   - Database: `whatsapp_clone`
   - Username: `postgres`
   - Password: `postgres123`

### Option B: DBeaver (Free, Multi-platform)

**Download:** https://dbeaver.io/download/

**Setup:**
1. Install DBeaver
2. Click "New Database Connection"
3. Select "PostgreSQL"
4. Enter the connection details above

### Option C: TablePlus (Mac - Pretty UI)

**Download:** https://tableplus.com/

**Setup:**
1. Click "Create a new connection"
2. Select PostgreSQL
3. Enter connection details

---

## 🔍 Useful SQL Queries

### Users Table

```sql
-- See all users
SELECT * FROM users;

-- See specific user details
SELECT id, username, email, display_name, is_online, created_at 
FROM users;

-- Find a user by email
SELECT * FROM users WHERE email = 'test@example.com';

-- Count total users
SELECT COUNT(*) as total_users FROM users;

-- See online users
SELECT username, display_name, last_seen_at 
FROM users 
WHERE is_online = true;
```

### Messages Table

```sql
-- See all messages
SELECT * FROM messages;

-- See messages with sender info
SELECT 
    m.id,
    m.content,
    m.message_type,
    m.created_at,
    u.username as sender_username,
    u.display_name as sender_name
FROM messages m
JOIN users u ON m.sender_id = u.id
ORDER BY m.created_at DESC;

-- Count messages per user
SELECT 
    u.username,
    COUNT(m.id) as message_count
FROM users u
LEFT JOIN messages m ON u.id = m.sender_id
GROUP BY u.username
ORDER BY message_count DESC;

-- Get messages from a specific conversation
SELECT 
    m.content,
    u.display_name as sender,
    m.created_at
FROM messages m
JOIN users u ON m.sender_id = u.id
WHERE m.conversation_id = 1
ORDER BY m.created_at ASC;

-- See unread messages
SELECT * FROM messages WHERE is_read = false;
```

### Conversations Table

```sql
-- See all conversations
SELECT * FROM conversations;

-- See conversations with member count
SELECT 
    c.id,
    c.type,
    c.name,
    COUNT(cm.id) as member_count
FROM conversations c
LEFT JOIN conversation_members cm ON c.id = cm.conversation_id
GROUP BY c.id;

-- See who's in a conversation
SELECT 
    c.id as conversation_id,
    c.name as conversation_name,
    u.username,
    u.display_name,
    cm.role
FROM conversations c
JOIN conversation_members cm ON c.id = cm.conversation_id
JOIN users u ON cm.user_id = u.id
WHERE c.id = 1;
```

### Complex Queries

```sql
-- Get user's conversations with last message
SELECT 
    c.id,
    c.type,
    c.name,
    m.content as last_message,
    m.created_at as last_message_time,
    sender.display_name as last_sender
FROM conversations c
JOIN conversation_members cm ON c.id = cm.conversation_id
LEFT JOIN LATERAL (
    SELECT * FROM messages 
    WHERE conversation_id = c.id 
    ORDER BY created_at DESC 
    LIMIT 1
) m ON true
LEFT JOIN users sender ON m.sender_id = sender.id
WHERE cm.user_id = 1  -- Replace with your user ID
ORDER BY m.created_at DESC;

-- Message statistics
SELECT 
    DATE(created_at) as date,
    COUNT(*) as message_count,
    COUNT(DISTINCT sender_id) as active_users
FROM messages
GROUP BY DATE(created_at)
ORDER BY date DESC;
```

---

## 🛠️ Maintenance Queries

### Check Database Size

```sql
SELECT 
    pg_size_pretty(pg_database_size('whatsapp_clone')) as database_size;
```

### Check Table Sizes

```sql
SELECT 
    table_name,
    pg_size_pretty(pg_total_relation_size(quote_ident(table_name))) as size
FROM information_schema.tables
WHERE table_schema = 'public'
ORDER BY pg_total_relation_size(quote_ident(table_name)) DESC;
```

### View Active Connections

```sql
SELECT 
    datname,
    usename,
    application_name,
    client_addr,
    state
FROM pg_stat_activity
WHERE datname = 'whatsapp_clone';
```

---

## 🧪 Testing Data Queries

### Insert Test Data

```sql
-- Insert a test user
INSERT INTO users (username, email, password, display_name, created_at, updated_at)
VALUES ('testuser2', 'test2@example.com', '$2a$10$HASH...', 'Test User 2', NOW(), NOW())
RETURNING *;

-- Insert a test message
INSERT INTO messages (conversation_id, sender_id, content, message_type, created_at, updated_at)
VALUES (1, 1, 'Test message from SQL!', 'text', NOW(), NOW())
RETURNING *;
```

### Update Data

```sql
-- Update user display name
UPDATE users 
SET display_name = 'New Display Name', updated_at = NOW()
WHERE id = 1;

-- Mark all messages as read
UPDATE messages 
SET is_read = true, read_at = NOW()
WHERE conversation_id = 1 AND sender_id != 1;
```

### Delete Data (Be Careful!)

```sql
-- Delete a specific message
DELETE FROM messages WHERE id = 5;

-- Delete all messages from a conversation
DELETE FROM messages WHERE conversation_id = 3;

-- Delete a user (careful - this might fail due to foreign keys)
DELETE FROM users WHERE id = 10;
```

---

## 📁 Export Data

### Export to CSV

```bash
# Export users to CSV
docker exec whatsapp_db psql -U postgres -d whatsapp_clone -c "\COPY users TO '/tmp/users.csv' CSV HEADER"
docker cp whatsapp_db:/tmp/users.csv ./users.csv

# Export messages to CSV
docker exec whatsapp_db psql -U postgres -d whatsapp_clone -c "\COPY messages TO '/tmp/messages.csv' CSV HEADER"
docker cp whatsapp_db:/tmp/messages.csv ./messages.csv
```

### Backup Entire Database

```bash
# Create backup
docker exec whatsapp_db pg_dump -U postgres whatsapp_clone > backup.sql

# Restore from backup
docker exec -i whatsapp_db psql -U postgres whatsapp_clone < backup.sql
```

---

## 🔧 Database Schema

### View Your Current Schema

```sql
-- See all tables and their columns
SELECT 
    table_name,
    column_name,
    data_type,
    is_nullable
FROM information_schema.columns
WHERE table_schema = 'public'
ORDER BY table_name, ordinal_position;
```

### Tables in Your Database:

**1. users**
- id, username, email, password
- display_name, avatar, bio
- is_online, last_seen_at
- created_at, updated_at, deleted_at

**2. messages**
- id, conversation_id, sender_id
- content, message_type, file_url
- is_read, read_at
- created_at, updated_at, deleted_at

**3. conversations**
- id, type (direct/group)
- name, avatar, description
- created_at, updated_at, deleted_at

**4. conversation_members**
- id, conversation_id, user_id
- role (admin/member)
- joined_at
- created_at, updated_at, deleted_at

---

## 🎯 Quick Start Commands

### Connect to Database
```bash
docker exec -it whatsapp_db psql -U postgres -d whatsapp_clone
```

### Once Connected - Essential Commands

```sql
-- See all tables
\dt

-- See users
SELECT * FROM users;

-- See messages
SELECT * FROM messages;

-- Count everything
SELECT 
    (SELECT COUNT(*) FROM users) as total_users,
    (SELECT COUNT(*) FROM messages) as total_messages,
    (SELECT COUNT(*) FROM conversations) as total_conversations;

-- Exit
\q
```

---

## 🐛 Troubleshooting

### Can't Connect?

```bash
# Check if container is running
docker ps | grep whatsapp_db

# Check container logs
docker logs whatsapp_db

# Restart database container
docker restart whatsapp_db
```

### Permission Denied?

Make sure you're using the correct credentials from docker-compose.yml:
- Username: `postgres`
- Password: `postgres123`

### Database Doesn't Exist?

```bash
# Create it manually
docker exec -it whatsapp_db psql -U postgres -c "CREATE DATABASE whatsapp_clone;"
```

---

## 📚 Learning Resources

**PostgreSQL Tutorial:**
- https://www.postgresqltutorial.com/

**SQL Practice:**
- https://sqlbolt.com/
- https://www.w3schools.com/sql/

**Understanding Your Schema:**
```sql
-- See table relationships
SELECT
    tc.table_name, 
    kcu.column_name,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name 
FROM information_schema.table_constraints AS tc 
JOIN information_schema.key_column_usage AS kcu
    ON tc.constraint_name = kcu.constraint_name
JOIN information_schema.constraint_column_usage AS ccu
    ON ccu.constraint_name = tc.constraint_name
WHERE tc.constraint_type = 'FOREIGN KEY' 
    AND tc.table_schema = 'public';
```

---

## 🎓 Pro Tips

1. **Use transactions for safety:**
```sql
BEGIN;
UPDATE users SET display_name = 'New Name' WHERE id = 1;
-- Check if it looks good
SELECT * FROM users WHERE id = 1;
-- If good: COMMIT; If bad: ROLLBACK;
COMMIT;
```

2. **Always backup before major changes:**
```bash
docker exec whatsapp_db pg_dump -U postgres whatsapp_clone > backup_$(date +%Y%m%d).sql
```

3. **Use LIMIT when exploring:**
```sql
SELECT * FROM messages LIMIT 10;  -- Safer than SELECT * FROM messages;
```

4. **Format your output:**
```sql
\x  -- Toggle expanded display (great for wide tables)
SELECT * FROM users LIMIT 1;
```

---

## 🚀 Next Steps

1. **Connect right now:**
   ```bash
   docker exec -it whatsapp_db psql -U postgres -d whatsapp_clone
   ```

2. **Run your first query:**
   ```sql
   SELECT COUNT(*) FROM users;
   ```

3. **Explore your data:**
   ```sql
   SELECT * FROM users;
   SELECT * FROM messages;
   ```

4. **Install a GUI tool** (optional but recommended)
   - pgAdmin, DBeaver, or TablePlus

Happy querying! 🎉
