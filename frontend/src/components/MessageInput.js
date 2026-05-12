import React, { useState } from 'react';
import { useChatStore } from '../store/chatStore';
import { messageAPI } from '../api';

function MessageInput() {
  const { currentConversation } = useChatStore();
  const [message, setMessage] = useState('');
  const [sending, setSending] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!message.trim() || !currentConversation || sending) return;

    setSending(true);
    try {
      await messageAPI.sendMessage({
        conversation_id: currentConversation.id,
        content: message,
        message_type: 'text',
      });
      setMessage('');
    } catch (error) {
      console.error('Error sending message:', error);
      alert('Failed to send message. Please try again.');
    } finally {
      setSending(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="bg-white p-4 border-t border-gray-300">
      <div className="flex items-center">
        <input
          type="text"
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          placeholder="Type a message..."
          className="flex-1 px-4 py-2 border border-gray-300 rounded-full focus:outline-none focus:ring-2 focus:ring-blue-500"
          disabled={sending}
        />
        <button
          type="submit"
          disabled={!message.trim() || sending}
          className="ml-3 bg-blue-500 text-white px-6 py-2 rounded-full hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
        >
          {sending ? 'Sending...' : 'Send'}
        </button>
      </div>
    </form>
  );
}

export default MessageInput;
