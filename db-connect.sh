#!/bin/bash
# Quick Database Connection Script
# Save this as: db-connect.sh
# Make executable: chmod +x db-connect.sh
# Run: ./db-connect.sh

echo "🗄️  WhatsApp Clone - Database Manager"
echo "======================================"
echo ""
echo "Connection Details:"
echo "  Database: whatsapp_clone"
echo "  User: postgres"
echo "  Container: whatsapp_db"
echo ""
echo "What would you like to do?"
echo ""
echo "1) Connect to database (interactive)"
echo "2) View all users"
echo "3) View all messages"
echo "4) View all conversations"
echo "5) Show database statistics"
echo "6) Backup database"
echo "7) View table structure"
echo "8) Run custom query"
echo "9) Exit"
echo ""
read -p "Enter your choice [1-9]: " choice

case $choice in
    1)
        echo "Connecting to database..."
        echo "Type '\dt' to list tables, '\q' to exit"
        docker exec -it whatsapp_db psql -U postgres -d whatsapp_clone
        ;;
    2)
        echo "📊 All Users:"
        docker exec whatsapp_db psql -U postgres -d whatsapp_clone -c "SELECT id, username, email, display_name, is_online, created_at FROM users;"
        ;;
    3)
        echo "💬 All Messages:"
        docker exec whatsapp_db psql -U postgres -d whatsapp_clone -c "SELECT m.id, u.username as sender, m.content, m.message_type, m.created_at FROM messages m JOIN users u ON m.sender_id = u.id ORDER BY m.created_at DESC;"
        ;;
    4)
        echo "💭 All Conversations:"
        docker exec whatsapp_db psql -U postgres -d whatsapp_clone -c "SELECT c.id, c.type, c.name, COUNT(cm.user_id) as members FROM conversations c LEFT JOIN conversation_members cm ON c.id = cm.conversation_id GROUP BY c.id;"
        ;;
    5)
        echo "📈 Database Statistics:"
        docker exec whatsapp_db psql -U postgres -d whatsapp_clone -c "SELECT 'Users' as table_name, COUNT(*) as count FROM users UNION ALL SELECT 'Messages', COUNT(*) FROM messages UNION ALL SELECT 'Conversations', COUNT(*) FROM conversations;"
        ;;
    6)
        BACKUP_FILE="backup_$(date +%Y%m%d_%H%M%S).sql"
        echo "Creating backup: $BACKUP_FILE"
        docker exec whatsapp_db pg_dump -U postgres whatsapp_clone > "$BACKUP_FILE"
        echo "✅ Backup created: $BACKUP_FILE"
        ;;
    7)
        echo "Which table? (users/messages/conversations/conversation_members)"
        read -p "Table name: " table_name
        docker exec whatsapp_db psql -U postgres -d whatsapp_clone -c "\d $table_name"
        ;;
    8)
        echo "Enter your SQL query:"
        read -p "Query: " query
        docker exec whatsapp_db psql -U postgres -d whatsapp_clone -c "$query"
        ;;
    9)
        echo "Goodbye! 👋"
        exit 0
        ;;
    *)
        echo "Invalid option. Please try again."
        ;;
esac
