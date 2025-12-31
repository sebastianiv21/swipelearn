import axios from 'axios';
import { authClient } from './auth-client';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:5050/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000,
});

// Request interceptor - Add auth token
api.interceptors.request.use(async (config) => {
  try {
        const session = await authClient.getSession();
    if (session?.data?.session?.token) {
      config.headers.Authorization = `Bearer ${session.data.session.token}`;
    }
  } catch {
    // Failed to get session for request
  }
  return config;
});

// Response interceptor - Handle token refresh and errors
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    // Handle 401 Unauthorized
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      
      try {
        // Try to refresh token
        await authClient.refreshToken();
        
        // Retry the original request with new token
    const session = await authClient.getSession();
        if (session?.data?.session?.token) {
          originalRequest.headers.Authorization = `Bearer ${session.data.session.token}`;
        }
        
        return api(originalRequest);
      } catch (refreshError) {
        // Refresh failed, sign out user
        await authClient.signOut();
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }

    // Handle other errors
    // Error status: ${error.response?.status}

    return Promise.reject(error);
  }
);

export default api;