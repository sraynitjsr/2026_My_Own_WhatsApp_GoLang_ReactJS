import React, { useEffect, useState } from 'react';
import { useAuthStore } from '../store/authStore';
import { useChatStore } from '../store/chatStore';
import { conversationAPI, messageAPI } from '../api';
import { connectWebSocket, disconnectWebSocket } from '../services/websocket';
import ConversationList from '../components/ConversationList';
import ChatWindow from '../components/ChatWindow';

function Chat() {
  const { user, token, logout } = useAuthStore();
  const { setConversations, setCurrentConversation } = useChatStore();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Connect WebSocket
    if (token) {
      connectWebSocket(token);
    }

    // Fetch conversations
    loadConversations();

    // Cleanup on unmount
    return () => {
      disconnectWebSocket();
    };
  }, [token]);

  const loadConversations = async () => {
    try {
      const conversations = await conversationAPI.getConversations();
      setConversations(conversations);
    } catch (error) {
      console.error('Error loading conversations:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    disconnectWebSocket();
    logout();
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p>Loading...</p>
      </div>
    );
  }

  return (
    <div className="h-screen flex bg-gray-100">
      {/* Sidebar */}
      <div className="w-1/3 bg-white border-r border-gray-300 flex flex-col">
        {/* Header */}
        <div className="bg-gray-100 p-4 flex items-center justify-between border-b border-gray-300">
          <div className="flex items-center">
            <div className="w-10 h-10 bg-blue-500 rounded-full flex items-center justify-center text-white font-bold">
              {user?.display_name?.[0]?.toUpperCase() || 'U'}
            </div>
            <span className="ml-3 font-semibold">{user?.display_name}</span>
          </div>
          <button
            onClick={handleLogout}
            className="text-gray-600 hover:text-gray-800 text-sm"
          >
            Logout
          </button>
        </div>

        {/* Conversation List */}
        <ConversationList />
      </div>

      {/* Chat Window */}
      <div className="flex-1 flex flex-col">
        <ChatWindow />
      </div>
    </div>
  );
}

export default Chat;
