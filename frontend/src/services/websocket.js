import { useChatStore } from '../store/chatStore';

const WS_URL = process.env.REACT_APP_WS_URL || 'ws://localhost:8080/api/v1/ws';

export const connectWebSocket = (token) => {
  const ws = new WebSocket(`${WS_URL}?token=${token}`);
  const { addMessage, setWebSocket } = useChatStore.getState();

  ws.onopen = () => {
    console.log('WebSocket connected');
    setWebSocket(ws);
  };

  ws.onmessage = (event) => {
    try {
      const message = JSON.parse(event.data);
      addMessage(message);
    } catch (error) {
      console.error('Error parsing WebSocket message:', error);
    }
  };

  ws.onerror = (error) => {
    console.error('WebSocket error:', error);
  };

  ws.onclose = () => {
    console.log('WebSocket disconnected');
    setWebSocket(null);
  };

  return ws;
};

export const disconnectWebSocket = () => {
  const { disconnectWebSocket } = useChatStore.getState();
  disconnectWebSocket();
};
