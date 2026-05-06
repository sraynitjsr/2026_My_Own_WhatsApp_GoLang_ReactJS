# 🎯 TablePlus Testing Guide - Run These Now!

## 🔌 Step 1: Connect TablePlus

**Click "Create a new connection" in TablePlus:**

```
Connection Name: WhatsApp Clone
Host: localhost
Port: 5432
User: postgres
Password: postgres123
Database: whatsapp_clone
```

Click **Test** → Then **Connect**

---

## 🚀 Step 2: Run These Queries (Copy & Paste!)

### Query 1: See Your Current Data
```sql
-- Quick overview of everything
SELECT 'Users' as table_name, COUNT(*) as records FROM users
UNION ALL
SELECT 'Messages', COUNT(*) FROM messages
UNION ALL
SELECT 'Conversations', COUNT(*) FROM conversations;
```

**Expected Result:**
```
Users          | 2
Messages       | 1
Conversations  | 1
```

---

### Query 2: View All Users
```sql
SELECT 
    id,
    username,
    email,
    display_name,
    is_online,
    created_at
FROM users
ORDER BY created_at DESC;
```

**Expected Result:** You'll see `testuser` and `alice`

---

### Query 3: View Messages with Sender Names
```sql
SELECT 
    m.id,
    u.username as sender,
    u.display_name as sender_name,
    m.content,
    m.message_type,
    m.is_read,
    m.created_at
FROM messages m
JOIN users u ON m.sender_id = u.id
ORDER BY m.created_at DESC;
```

**Expected Result:** "Hello! This is a test message!" from Alice

---

### Query 4: See Who's in Each Conversation
```sql
SELECT 
    c.id as conv_id,
    c.type,
    u1.username as member1,
    u2.username as member2
FROM conversations c
JOIN conversation_members cm1 ON c.id = cm1.conversation_id
JOIN users u1 ON cm1.user_id = u1.id
LEFT JOIN conversation_members cm2 ON c.id = cm2.conversation_id AND cm2.user_id != cm1.user_id
LEFT JOIN users u2 ON cm2.user_id = u2.id
WHERE cm1.user_id < COALESCE(cm2.user_id, 999)
   OR cm2.user_id IS NULL;
```

**Expected Result:** testuser ↔️ alice (conversation 1)

---

## 🧪 Step 3: Test Your App by Adding Data

### Test 1: Register a New User via API

**In Terminal, run:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "bob",
    "email": "bob@example.com",
    "password": "password123",
    "display_name": "Bob Johnson"
  }'
```

**Then in TablePlus, run:**
```sql
SELECT * FROM users ORDER BY created_at DESC LIMIT 1;
```

**You should see:** Bob Johnson just appeared! ✅

---

### Test 2: Create a Group Chat via API

**In Terminal:**
```bash
# First, get Alice's token (you created this earlier)
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFsaWNlQGV4YW1wbGUuY29tIiwiZXhwIjoxNzc4MTMwNzQxLCJ1c2VyX2lkIjoyfQ.Mkb7ExILYxr96qtcxynPgEajRiDkZHEglexWHw7MArw"

# Create a group chat
curl -X POST http://localhost:8080/api/v1/conversations \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "group",
    "name": "Team Chat",
    "description": "Our awesome team",
    "member_ids": [1, 3]
  }'
```

**Then in TablePlus:**
```sql
-- See all conversations
SELECT 
    c.id,
    c.type,
    c.name,
    COUNT(cm.user_id) as members
FROM conversations c
LEFT JOIN conversation_members cm ON c.id = cm.conversation_id
GROUP BY c.id;
```

**You should see:** "Team Chat" group with 3 members! ✅

---

### Test 3: Send Multiple Messages

**In Terminal:**
```bash
TOKEN="YOUR_TOKEN_HERE"

# Send message 1
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"conversation_id": 1, "content": "Hey there!", "message_type": "text"}'

# Send message 2
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"conversation_id": 1, "content": "How are you?", "message_type": "text"}'

# Send message 3
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"conversation_id": 1, "content": "Great weather today! 🌞", "message_type": "text"}'
```

**Then in TablePlus:**
```sql
-- See conversation timeline
SELECT 
    m.id,
    u.display_name as sender,
    m.content,
    m.created_at::time as time_sent
FROM messages m
JOIN users u ON m.sender_id = u.id
WHERE m.conversation_id = 1
ORDER BY m.created_at ASC;
```

**You should see:** Full conversation thread! ✅

---

## 🔥 Step 4: Advanced Testing Queries

### Test 4: Message Activity Analysis
```sql
-- Messages per user
SELECT 
    u.username,
    u.display_name,
    COUNT(m.id) as messages_sent,
    MAX(m.created_at) as last_message_time
FROM users u
LEFT JOIN messages m ON u.id = m.sender_id
GROUP BY u.id, u.username, u.display_name
ORDER BY messages_sent DESC;
```

**What to look for:** Who's the most active user?

---

### Test 5: Conversation Activity
```sql
-- Most active conversations
SELECT 
    c.id,
    COALESCE(c.name, 'Direct Message') as conversation,
    c.type,
    COUNT(DISTINCT m.id) as total_messages,
    COUNT(DISTINCT m.sender_id) as active_users,
    MAX(m.created_at) as last_activity
FROM conversations c
LEFT JOIN messages m ON c.id = m.conversation_id
GROUP BY c.id, c.name, c.type
ORDER BY total_messages DESC;
```

**What to look for:** Which chat has the most activity?

---

### Test 6: User's Full Conversation List (Like WhatsApp UI)
```sql
-- Show all conversations for user ID 1 with last message
SELECT 
    c.id,
    c.type,
    CASE 
        WHEN c.type = 'group' THEN c.name
        ELSE other_user.display_name
    END as conversation_name,
    last_msg.content as last_message,
    last_msg.sender_name,
    last_msg.created_at as last_message_time,
    unread_count.count as unread_messages
FROM conversations c
JOIN conversation_members cm ON c.id = cm.conversation_id
-- Get the other user in direct chats
LEFT JOIN LATERAL (
    SELECT u.display_name
    FROM conversation_members cm2
    JOIN users u ON cm2.user_id = u.id
    WHERE cm2.conversation_id = c.id 
      AND cm2.user_id != 1
      AND c.type = 'direct'
    LIMIT 1
) other_user ON true
-- Get last message
LEFT JOIN LATERAL (
    SELECT 
        m.content,
        m.created_at,
        u.display_name as sender_name
    FROM messages m
    JOIN users u ON m.sender_id = u.id
    WHERE m.conversation_id = c.id
    ORDER BY m.created_at DESC
    LIMIT 1
) last_msg ON true
-- Count unread messages
LEFT JOIN LATERAL (
    SELECT COUNT(*) as count
    FROM messages m
    WHERE m.conversation_id = c.id
      AND m.sender_id != 1
      AND m.is_read = false
) unread_count ON true
WHERE cm.user_id = 1
ORDER BY last_msg.created_at DESC NULLS LAST;
```

**What to look for:** This looks exactly like WhatsApp's chat list! 📱

---

## 🎮 Step 5: Interactive Testing Scenarios

### Scenario 1: Simulate a Real Conversation

**Run these in sequence to see data flow:**

```sql
-- 1. Check current messages
SELECT COUNT(*) as message_count FROM messages WHERE conversation_id = 1;

-- 2. Send message via API (Terminal)
-- curl -X POST http://localhost:8080/api/v1/messages ...

-- 3. Verify message was saved
SELECT * FROM messages ORDER BY created_at DESC LIMIT 1;

-- 4. Check unread count
SELECT COUNT(*) FROM messages 
WHERE conversation_id = 1 AND is_read = false;

-- 5. Mark as read
UPDATE messages 
SET is_read = true, read_at = NOW()
WHERE conversation_id = 1 AND sender_id != 1;

-- 6. Verify all read
SELECT COUNT(*) FROM messages 
WHERE conversation_id = 1 AND is_read = false;
```

---

### Scenario 2: Test User Online Status

```sql
-- 1. Check current online status
SELECT username, is_online, last_seen_at FROM users;

-- 2. Set user online
UPDATE users 
SET is_online = true, updated_at = NOW()
WHERE id = 1;

-- 3. Wait a moment, then set offline
UPDATE users 
SET is_online = false, last_seen_at = NOW(), updated_at = NOW()
WHERE id = 1;

-- 4. View updated status
SELECT username, is_online, last_seen_at FROM users;
```

---

### Scenario 3: Group Chat Management

```sql
-- 1. Create a group (or use existing)
SELECT * FROM conversations WHERE type = 'group';

-- 2. See all members
SELECT 
    c.name as group_name,
    u.display_name as member,
    cm.role,
    cm.joined_at
FROM conversations c
JOIN conversation_members cm ON c.id = cm.conversation_id
JOIN users u ON cm.user_id = u.id
WHERE c.type = 'group'
ORDER BY c.id, cm.role DESC;

-- 3. Count messages per member in group
SELECT 
    u.display_name,
    COUNT(m.id) as messages_in_group
FROM conversation_members cm
JOIN users u ON cm.user_id = u.id
LEFT JOIN messages m ON m.sender_id = u.id AND m.conversation_id = cm.conversation_id
WHERE cm.conversation_id = 2  -- Change to your group ID
GROUP BY u.display_name
ORDER BY messages_in_group DESC;
```

---

## 📊 Step 6: Data Validation Queries

### Validate Data Integrity

```sql
-- Check for orphaned messages (no sender)
SELECT COUNT(*) as orphaned_messages
FROM messages m
LEFT JOIN users u ON m.sender_id = u.id
WHERE u.id IS NULL;
-- Should be: 0

-- Check for orphaned conversation members
SELECT COUNT(*) as orphaned_members
FROM conversation_members cm
LEFT JOIN users u ON cm.user_id = u.id
WHERE u.id IS NULL;
-- Should be: 0

-- Check for conversations without members
SELECT c.id, c.type, c.name
FROM conversations c
LEFT JOIN conversation_members cm ON c.id = cm.conversation_id
WHERE cm.id IS NULL;
-- Should be: empty

-- Check for users without conversations
SELECT u.id, u.username, COUNT(cm.id) as conversation_count
FROM users u
LEFT JOIN conversation_members cm ON u.id = cm.user_id
GROUP BY u.id, u.username
HAVING COUNT(cm.id) = 0;
-- These are users who haven't joined any chat yet
```

---

## 🎯 Quick Reference: Testing Checklist

Run these in order to fully test your app:

```sql
-- ✅ Step 1: Verify base data exists
SELECT COUNT(*) FROM users;          -- Should be >= 2
SELECT COUNT(*) FROM conversations;  -- Should be >= 1
SELECT COUNT(*) FROM messages;       -- Should be >= 1

-- ✅ Step 2: Test authentication worked
SELECT username, email, created_at FROM users ORDER BY created_at DESC;

-- ✅ Step 3: Test messaging system
SELECT 
    u.username as sender,
    m.content,
    m.created_at
FROM messages m
JOIN users u ON m.sender_id = u.id
ORDER BY m.created_at DESC
LIMIT 10;

-- ✅ Step 4: Test conversation structure
SELECT 
    c.id,
    c.type,
    COUNT(cm.user_id) as members,
    COUNT(m.id) as messages
FROM conversations c
LEFT JOIN conversation_members cm ON c.id = cm.conversation_id
LEFT JOIN messages m ON c.id = m.conversation_id
GROUP BY c.id, c.type;

-- ✅ Step 5: Test relationships work
SELECT 
    'Users with no conversations' as test,
    COUNT(*) as count
FROM users u
LEFT JOIN conversation_members cm ON u.id = cm.user_id
WHERE cm.id IS NULL
UNION ALL
SELECT 
    'Messages with no sender',
    COUNT(*)
FROM messages m
LEFT JOIN users u ON m.sender_id = u.id
WHERE u.id IS NULL
UNION ALL
SELECT 
    'Conversations with no messages',
    COUNT(*)
FROM conversations c
LEFT JOIN messages m ON c.id = m.conversation_id
WHERE m.id IS NULL;
```

---

## 🔥 Pro Tips for TablePlus

### Tip 1: Save Your Favorite Queries
- Click "⭐" in TablePlus to save frequently used queries
- Create a "Testing" folder for these queries

### Tip 2: Use Multiple Tabs
- Open several query tabs
- Keep one for viewing data, one for testing

### Tip 3: Export Results
- Right-click query results → Export → CSV/JSON
- Great for sharing with team or analysis

### Tip 4: View Table Structure
- Click on table name in left sidebar
- See all columns, indexes, relationships

### Tip 5: Quick Filters
- Click column header in results
- Type to filter
- Much faster than writing WHERE clauses!

---

## 🚨 Safety Tips

**Always use transactions for updates:**
```sql
BEGIN;
-- Your UPDATE/DELETE query here
SELECT * FROM table_name WHERE...; -- Check before committing
-- If good: COMMIT;
-- If bad: ROLLBACK;
```

**Create backups before major changes:**
```bash
docker exec whatsapp_db pg_dump -U postgres whatsapp_clone > backup_before_test.sql
```

---

## 🎉 You're Ready!

**Right now, do this:**

1. **Open TablePlus**
2. **Create connection** (details at top)
3. **Click "Connect"**
4. **Copy the first query** (Step 2, Query 1)
5. **Press Cmd+Enter** (Mac) or **Ctrl+Enter** (Windows)
6. **See your data!** 🎊

**Then try:** Adding a new user via API and immediately checking TablePlus to see it appear!

---

## 📚 Full Test Flow Example

```bash
# Terminal 1: Send API request
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"charlie","email":"charlie@test.com","password":"pass123","display_name":"Charlie"}'
```

```sql
-- TablePlus: Verify it saved
SELECT * FROM users WHERE username = 'charlie';

-- See the token was created (it's in the API response)
-- Use it to send a message
```

```bash
# Terminal: Send message with Charlie's token
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Authorization: Bearer CHARLIE_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{"conversation_id":1,"content":"Hi from Charlie!","message_type":"text"}'
```

```sql
-- TablePlus: See message appear instantly
SELECT * FROM messages ORDER BY created_at DESC LIMIT 1;
```

**You just tested the full flow: API → Database → Verification!** ✅
