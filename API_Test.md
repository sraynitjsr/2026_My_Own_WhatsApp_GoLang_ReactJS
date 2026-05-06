# API Testing Guide

This guide shows you how to test the backend API using curl or any HTTP client.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication Flow

### 1. Register a New User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "password123",
    "display_name": "John Doe"
  }'
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "display_name": "John Doe",
    "is_online": false
  }
}
```

Save the `token` - you'll need it for authenticated requests!

### 2. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

## User Endpoints

### Get Current User Profile

```bash
curl http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Update User Profile

```bash
curl -X PUT http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "display_name": "John Smith",
    "bio": "Software Developer"
  }'
```

### Search Users

```bash
curl "http://localhost:8080/api/v1/users/search?q=john" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Update Online Status

```bash
curl -X PUT http://localhost:8080/api/v1/users/status \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "is_online": true
  }'
```

## Conversation Endpoints

### Create a Direct Conversation

```bash
curl -X POST http://localhost:8080/api/v1/conversations \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "direct",
    "member_ids": [2]
  }'
```

### Create a Group Conversation

```bash
curl -X POST http://localhost:8080/api/v1/conversations \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "group",
    "name": "Team Chat",
    "description": "Our awesome team group",
    "member_ids": [2, 3, 4]
  }'
```

### Get All Conversations

```bash
curl http://localhost:8080/api/v1/conversations \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Get Specific Conversation

```bash
curl http://localhost:8080/api/v1/conversations/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Update Conversation (Admins only)

```bash
curl -X PUT http://localhost:8080/api/v1/conversations/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Team Name",
    "description": "New description"
  }'
```

## Message Endpoints

### Send a Text Message

```bash
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "conversation_id": 1,
    "content": "Hello, World!",
    "message_type": "text"
  }'
```

### Get Messages from a Conversation

```bash
curl http://localhost:8080/api/v1/messages/conversation/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Mark Message as Read

```bash
curl -X PUT http://localhost:8080/api/v1/messages/1/read \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Delete a Message

```bash
curl -X DELETE http://localhost:8080/api/v1/messages/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## File Upload

### Upload a File

```bash
curl -X POST http://localhost:8080/api/v1/upload \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -F "file=@/path/to/your/file.jpg"
```

**Response:**
```json
{
  "url": "/uploads/1234567890_file.jpg",
  "filename": "file.jpg",
  "size": 102400
}
```

### Send a File Message

After uploading, send a message with the file URL:

```bash
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "conversation_id": 1,
    "content": "Check out this image!",
    "message_type": "image",
    "file_url": "/uploads/1234567890_file.jpg"
  }'
```

## WebSocket Connection

To connect to WebSocket for real-time messages:

```javascript
// In JavaScript
const token = "YOUR_JWT_TOKEN";
const ws = new WebSocket(`ws://localhost:8080/api/v1/ws?token=${token}`);

ws.onopen = () => {
  console.log('Connected to WebSocket');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('New message:', message);
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('Disconnected from WebSocket');
};
```

## Complete Example Workflow

Here's a complete test scenario:

```bash
# 1. Register User A
USER_A=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "password": "password123",
    "display_name": "Alice"
  }')

TOKEN_A=$(echo $USER_A | jq -r '.token')
echo "Alice Token: $TOKEN_A"

# 2. Register User B
USER_B=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "bob",
    "email": "bob@example.com",
    "password": "password123",
    "display_name": "Bob"
  }')

TOKEN_B=$(echo $USER_B | jq -r '.token')
echo "Bob Token: $TOKEN_B"

# 3. Alice creates a conversation with Bob
CONV=$(curl -s -X POST http://localhost:8080/api/v1/conversations \
  -H "Authorization: Bearer $TOKEN_A" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "direct",
    "member_ids": [2]
  }')

CONV_ID=$(echo $CONV | jq -r '.id')
echo "Conversation ID: $CONV_ID"

# 4. Alice sends a message
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Authorization: Bearer $TOKEN_A" \
  -H "Content-Type: application/json" \
  -d "{
    \"conversation_id\": $CONV_ID,
    \"content\": \"Hi Bob!\",
    \"message_type\": \"text\"
  }"

# 5. Bob reads the messages
curl http://localhost:8080/api/v1/messages/conversation/$CONV_ID \
  -H "Authorization: Bearer $TOKEN_B"

# 6. Bob replies
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Authorization: Bearer $TOKEN_B" \
  -H "Content-Type: application/json" \
  -d "{
    \"conversation_id\": $CONV_ID,
    \"content\": \"Hi Alice! How are you?\",
    \"message_type\": \"text\"
  }"
```

## Common HTTP Status Codes

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Missing or invalid token
- `403 Forbidden` - You don't have permission
- `404 Not Found` - Resource doesn't exist
- `500 Internal Server Error` - Server error

## Tips for Testing

1. **Save your tokens**: Export them as environment variables:
   ```bash
   export TOKEN="your_jwt_token_here"
   curl -H "Authorization: Bearer $TOKEN" ...
   ```

2. **Use jq for pretty JSON**: Install jq and pipe responses:
   ```bash
   curl ... | jq '.'
   ```

3. **Use Postman or Insomnia**: GUI tools make API testing easier
   - Import these curl commands
   - Save tokens automatically
   - Better visualization

4. **Check backend logs**: Watch the backend terminal for errors

5. **Test in order**: Create users → Create conversations → Send messages
