-- ============================================
-- COMMON SQL QUERIES FOR WHATSAPP CLONE
-- Copy and paste these into your psql session
-- ============================================

-- ============================================
-- QUICK CHECKS
-- ============================================

-- List all tables
\dt

-- Count everything
SELECT 
    'Users' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'Messages', COUNT(*) FROM messages
UNION ALL
SELECT 'Conversations', COUNT(*) FROM conversations
UNION ALL
SELECT 'Conversation Members', COUNT(*) FROM conversation_members;


-- ============================================
-- USER QUERIES
-- ============================================

-- View all users (safe columns only)
SELECT id, username, email, display_name, is_online, created_at 
FROM users 
ORDER BY created_at DESC;

-- Find user by email
SELECT id, username, email, display_name 
FROM users 
WHERE email = 'test@example.com';

-- Find user by username
SELECT id, username, email, display_name 
FROM users 
WHERE username = 'alice';

-- See who's online
SELECT username, display_name, last_seen_at 
FROM users 
WHERE is_online = true;

-- User registration timeline
SELECT 
    DATE(created_at) as signup_date,
    COUNT(*) as new_users
FROM users
GROUP BY DATE(created_at)
ORDER BY signup_date DESC;


-- ============================================
-- MESSAGE QUERIES
-- ============================================

-- View all messages with sender info
SELECT 
    m.id,
    m.content,
    u.username as sender,
    m.message_type,
    m.is_read,
    m.created_at
FROM messages m
JOIN users u ON m.sender_id = u.id
ORDER BY m.created_at DESC
LIMIT 20;

-- Get messages from specific conversation
SELECT 
    m.content,
    u.display_name as sender,
    m.created_at,
    m.is_read
FROM messages m
JOIN users u ON m.sender_id = u.id
WHERE m.conversation_id = 1
ORDER BY m.created_at ASC;

-- Count messages per user
SELECT 
    u.username,
    u.display_name,
    COUNT(m.id) as total_messages
FROM users u
LEFT JOIN messages m ON u.id = m.sender_id
GROUP BY u.id, u.username, u.display_name
ORDER BY total_messages DESC;

-- Find unread messages
SELECT 
    m.id,
    m.content,
    u.username as sender,
    m.created_at
FROM messages m
JOIN users u ON m.sender_id = u.id
WHERE m.is_read = false
ORDER BY m.created_at DESC;

-- Messages sent today
SELECT 
    u.username,
    m.content,
    m.created_at
FROM messages m
JOIN users u ON m.sender_id = u.id
WHERE DATE(m.created_at) = CURRENT_DATE
ORDER BY m.created_at DESC;


-- ============================================
-- CONVERSATION QUERIES
-- ============================================

-- All conversations with details
SELECT 
    c.id,
    c.type,
    COALESCE(c.name, 'Direct Message') as name,
    COUNT(DISTINCT cm.user_id) as member_count,
    COUNT(DISTINCT m.id) as message_count,
    MAX(m.created_at) as last_message_time
FROM conversations c
LEFT JOIN conversation_members cm ON c.id = cm.conversation_id
LEFT JOIN messages m ON c.id = m.conversation_id
GROUP BY c.id, c.type, c.name
ORDER BY last_message_time DESC;

-- Who's in which conversation
SELECT 
    c.id as conv_id,
    COALESCE(c.name, 'Direct Chat') as conversation,
    u.username,
    u.display_name,
    cm.role
FROM conversations c
JOIN conversation_members cm ON c.id = cm.conversation_id
JOIN users u ON cm.user_id = u.id
ORDER BY c.id, cm.role DESC;

-- Conversations for a specific user (replace user_id)
SELECT 
    c.id,
    c.type,
    c.name,
    COUNT(m.id) as messages
FROM conversations c
JOIN conversation_members cm ON c.id = cm.conversation_id
LEFT JOIN messages m ON c.id = m.conversation_id
WHERE cm.user_id = 1  -- Change this to your user ID
GROUP BY c.id, c.type, c.name;


-- ============================================
-- ADVANCED QUERIES
-- ============================================

-- User's conversations with last message
SELECT 
    c.id,
    c.type,
    COALESCE(c.name, 'Direct Message') as conversation_name,
    last_msg.content as last_message,
    last_msg.created_at as last_message_time,
    sender.display_name as last_sender
FROM conversations c
JOIN conversation_members cm ON c.id = cm.conversation_id
LEFT JOIN LATERAL (
    SELECT m.*, u.display_name as sender_name
    FROM messages m
    JOIN users u ON m.sender_id = u.id
    WHERE m.conversation_id = c.id
    ORDER BY m.created_at DESC
    LIMIT 1
) last_msg ON true
LEFT JOIN users sender ON last_msg.sender_id = sender.id
WHERE cm.user_id = 1  -- Change to your user ID
ORDER BY last_msg.created_at DESC NULLS LAST;

-- Most active conversations
SELECT 
    c.id,
    COALESCE(c.name, 'Direct Chat') as conversation,
    COUNT(m.id) as total_messages,
    COUNT(DISTINCT m.sender_id) as active_users,
    MIN(m.created_at) as first_message,
    MAX(m.created_at) as last_message
FROM conversations c
LEFT JOIN messages m ON c.id = m.conversation_id
GROUP BY c.id, c.name
ORDER BY total_messages DESC;

-- Daily message activity
SELECT 
    DATE(created_at) as date,
    COUNT(*) as messages_sent,
    COUNT(DISTINCT sender_id) as active_users
FROM messages
GROUP BY DATE(created_at)
ORDER BY date DESC;


-- ============================================
-- DATA MODIFICATION (BE CAREFUL!)
-- ============================================

-- Mark messages as read (in a transaction for safety)
BEGIN;
UPDATE messages 
SET is_read = true, read_at = NOW()
WHERE conversation_id = 1 
  AND sender_id != 1  -- Not your own messages
  AND is_read = false;
-- Check before committing:
SELECT * FROM messages WHERE conversation_id = 1;
-- If OK: COMMIT; If not: ROLLBACK;
COMMIT;

-- Update user display name
BEGIN;
UPDATE users 
SET display_name = 'New Name', updated_at = NOW()
WHERE id = 1;
-- Check: SELECT * FROM users WHERE id = 1;
COMMIT;

-- Set user online/offline
UPDATE users 
SET is_online = true, updated_at = NOW()
WHERE id = 1;


-- ============================================
-- MAINTENANCE & DIAGNOSTICS
-- ============================================

-- Database size
SELECT pg_size_pretty(pg_database_size('whatsapp_clone')) as size;

-- Table sizes
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- Index information
SELECT 
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE schemaname = 'public'
ORDER BY tablename, indexname;

-- Active connections
SELECT 
    pid,
    usename,
    application_name,
    client_addr,
    state,
    query
FROM pg_stat_activity
WHERE datname = 'whatsapp_clone';


-- ============================================
-- TESTING & DEVELOPMENT
-- ============================================

-- Insert test user (password is 'password123' hashed with bcrypt)
INSERT INTO users (username, email, password, display_name, created_at, updated_at)
VALUES (
    'newuser',
    'newuser@example.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'New User',
    NOW(),
    NOW()
)
RETURNING *;

-- Insert test message
INSERT INTO messages (conversation_id, sender_id, content, message_type, created_at, updated_at)
VALUES (1, 1, 'Test message from SQL!', 'text', NOW(), NOW())
RETURNING *;

-- Create test direct conversation
INSERT INTO conversations (type, created_at, updated_at)
VALUES ('direct', NOW(), NOW())
RETURNING *;


-- ============================================
-- CLEANUP (USE WITH CAUTION!)
-- ============================================

-- Delete old test data
-- DELETE FROM messages WHERE content LIKE '%test%';
-- DELETE FROM users WHERE email LIKE '%test%';

-- Clear all data (DANGEROUS!)
-- TRUNCATE messages, conversation_members, conversations, users CASCADE;


-- ============================================
-- HELPFUL PSQL COMMANDS (not SQL)
-- ============================================

-- \dt              List all tables
-- \d table_name    Describe table structure
-- \l               List all databases
-- \du              List all users/roles
-- \x               Toggle expanded display
-- \q               Quit
-- \?               Help
-- \timing          Show query execution time
