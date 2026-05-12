import React, { useEffect, useState } from 'react';
import { useChatStore } from '../store/chatStore';
import { messageAPI } from '../api';
import MessageInput from './MessageInput';

function ChatWindow() {
  const { currentConversation, messages, setMessages } = useChatStore();
  const [loading, setLoading] = useState(false);
  const currentUser = JSON.parse(localStorage.getItem('user'));

  useEffect(() => {
    if (currentConversation) {
      loadMessages();
    }
  }, [currentConversation]);

  const loadMessages = async () => {
    if (!currentConversation) return;
    
    setLoading(true);
    try {
      const msgs = await messageAPI.getMessages(currentConversation.id);
      setMessages(msgs);
    } catch (error) {
      console.error('Error loading messages:', error);
    } finally {
      setLoading(false);
    }
  };

  if (!currentConversation) {
    return (
      <div className="flex-1 flex items-center justify-center bg-gray-50">
        <p className="text-gray-500">Select a conversation to start chatting</p>
      </div>
    );
  }

  return (
    <div className="flex-1 flex flex-col">
      {/* Chat Header */}
      <div className="bg-gray-100 p-4 border-b border-gray-300 flex items-center">
        <div className="w-10 h-10 bg-blue-500 rounded-full flex items-center justify-center text-white font-bold">
          {currentConversation.name?.[0]?.toUpperCase() || 'C'}
        </div>
        <span className="ml-3 font-semibold">
          {currentConversation.name || 'Chat'}
        </span>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4 bg-gray-50">
        {loading ? (
          <p className="text-center text-gray-500">Loading messages...</p>
        ) : messages.length === 0 ? (
          <p className="text-center text-gray-500">No messages yet. Start the conversation!</p>
        ) : (
          messages.map((message) => {
            const isOwnMessage = message.sender_id === currentUser.id;
            return (
              <div
                key={message.id}
                className={`mb-4 flex ${isOwnMessage ? 'justify-end' : 'justify-start'}`}
              >
                <div
                  className={`max-w-xs lg:max-w-md px-4 py-2 rounded-lg ${
                    isOwnMessage
                      ? 'bg-blue-500 text-white'
                      : 'bg-white text-gray-800 border border-gray-300'
                  }`}
                >
                  {!isOwnMessage && (
                    <p className="text-xs font-semibold mb-1">
                      {message.sender?.display_name}
                    </p>
                  )}
                  <p>{message.content}</p>
                  <p className={`text-xs mt-1 ${isOwnMessage ? 'text-blue-100' : 'text-gray-500'}`}>
                    {new Date(message.created_at).toLocaleTimeString([], {
                      hour: '2-digit',
                      minute: '2-digit'
                    })}
                  </p>
                </div>
              </div>
            );
          })
        )}
      </div>

      {/* Message Input */}
      <MessageInput />
    </div>
  );
}

export default ChatWindow;
