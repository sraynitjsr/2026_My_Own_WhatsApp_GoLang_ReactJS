import React from 'react';
import { useChatStore } from '../store/chatStore';

function ConversationList() {
  const { conversations, currentConversation, setCurrentConversation } = useChatStore();

  const getConversationName = (conversation) => {
    if (conversation.type === 'group') {
      return conversation.name || 'Group Chat';
    }
    // For direct messages, show the other user's name
    const otherMember = conversation.members?.find(
      member => member.user_id !== JSON.parse(localStorage.getItem('user')).id
    );
    return otherMember?.user?.display_name || 'Unknown User';
  };

  const getLastMessage = (conversation) => {
    if (conversation.messages && conversation.messages.length > 0) {
      return conversation.messages[0].content.substring(0, 50);
    }
    return 'No messages yet';
  };

  return (
    <div className="flex-1 overflow-y-auto">
      {conversations.length === 0 ? (
        <div className="p-4 text-center text-gray-500">
          No conversations yet. Start a new chat!
        </div>
      ) : (
        conversations.map((conversation) => (
          <div
            key={conversation.id}
            onClick={() => setCurrentConversation(conversation)}
            className={`p-4 border-b border-gray-200 cursor-pointer hover:bg-gray-50 ${
              currentConversation?.id === conversation.id ? 'bg-gray-100' : ''
            }`}
          >
            <div className="flex items-center">
              <div className="w-12 h-12 bg-blue-500 rounded-full flex items-center justify-center text-white font-bold">
                {getConversationName(conversation)[0]?.toUpperCase()}
              </div>
              <div className="ml-3 flex-1">
                <div className="flex justify-between items-center">
                  <span className="font-semibold">{getConversationName(conversation)}</span>
                  {conversation.messages && conversation.messages[0] && (
                    <span className="text-xs text-gray-500">
                      {new Date(conversation.messages[0].created_at).toLocaleTimeString([], {
                        hour: '2-digit',
                        minute: '2-digit'
                      })}
                    </span>
                  )}
                </div>
                <p className="text-sm text-gray-600 truncate">{getLastMessage(conversation)}</p>
              </div>
            </div>
          </div>
        ))
      )}
    </div>
  );
}

export default ConversationList;
