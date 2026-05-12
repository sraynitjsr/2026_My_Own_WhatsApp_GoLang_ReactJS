import api from './axios';

export const authAPI = {
  register: async (data) => {
    const response = await api.post('/auth/register', data);
    return response.data;
  },

  login: async (data) => {
    const response = await api.post('/auth/login', data);
    return response.data;
  },
};

export const userAPI = {
  getCurrentUser: async () => {
    const response = await api.get('/users/me');
    return response.data;
  },

  updateCurrentUser: async (data) => {
    const response = await api.put('/users/me', data);
    return response.data;
  },

  searchUsers: async (query) => {
    const response = await api.get(`/users/search?q=${query}`);
    return response.data;
  },

  updateOnlineStatus: async (isOnline) => {
    const response = await api.put('/users/status', { is_online: isOnline });
    return response.data;
  },
};

export const conversationAPI = {
  getConversations: async () => {
    const response = await api.get('/conversations');
    return response.data;
  },

  getConversation: async (id) => {
    const response = await api.get(`/conversations/${id}`);
    return response.data;
  },

  createConversation: async (data) => {
    const response = await api.post('/conversations', data);
    return response.data;
  },

  updateConversation: async (id, data) => {
    const response = await api.put(`/conversations/${id}`, data);
    return response.data;
  },

  deleteConversation: async (id) => {
    const response = await api.delete(`/conversations/${id}`);
    return response.data;
  },
};

export const messageAPI = {
  getMessages: async (conversationId) => {
    const response = await api.get(`/messages/conversation/${conversationId}`);
    return response.data;
  },

  sendMessage: async (data) => {
    const response = await api.post('/messages', data);
    return response.data;
  },

  markAsRead: async (messageId) => {
    const response = await api.put(`/messages/${messageId}/read`);
    return response.data;
  },

  deleteMessage: async (messageId) => {
    const response = await api.delete(`/messages/${messageId}`);
    return response.data;
  },
};

export const fileAPI = {
  uploadFile: async (file) => {
    const formData = new FormData();
    formData.append('file', file);
    const response = await api.post('/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },
};
