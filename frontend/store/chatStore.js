import { create } from 'zustand';

export const useChatStore = create((set, get) => ({
  conversations: [],
  currentConversation: null,
  messages: [],
  ws: null,

  setConversations: (conversations) => set({ conversations }),

  setCurrentConversation: (conversation) => set({ 
    currentConversation: conversation,
    messages: [] 
  }),

  setMessages: (messages) => set({ messages }),

  addMessage: (message) => set((state) => ({
    messages: [...state.messages, message],
  })),

  updateMessage: (messageId, updates) => set((state) => ({
    messages: state.messages.map(msg => 
      msg.id === messageId ? { ...msg, ...updates } : msg
    ),
  })),

  setWebSocket: (ws) => set({ ws }),

  disconnectWebSocket: () => {
    const { ws } = get();
    if (ws) {
      ws.close();
      set({ ws: null });
    }
  },
}));
