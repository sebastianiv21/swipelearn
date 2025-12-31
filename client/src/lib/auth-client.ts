import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import api from './api';

export interface AuthUser {
  id: string;
  email: string;
  name: string;
  created_at: string;
  updated_at: string;
}

interface AuthSession {
  access_token: string;
  refresh_token: string;
}

interface AuthState {
  user: AuthUser | null;
  session: AuthSession | null;
  isPending: boolean;
  error: string | null;
  setAuth: (user: AuthUser, session: AuthSession) => void;
  clearAuth: () => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
}

const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      session: null,
      isPending: false,
      error: null,
      setAuth: (user, session) => set({ user, session, error: null }),
      clearAuth: () => set({ user: null, session: null, error: null }),
      setLoading: (isPending) => set({ isPending }),
      setError: (error) => set({ error }),
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({ user: state.user, session: state.session }),
    }
  )
);

// Auth client implementation
export const authClient = {
  signIn: {
    email: async (credentials: { email: string; password: string }) => {
      const { setLoading, setError, setAuth } = useAuthStore.getState();
      setLoading(true);
      setError(null);

      try {
        const response = await api.post('/auth/login', credentials);
        const { user, access_token, refresh_token } = response.data;
        
        setAuth(user, { access_token, refresh_token });
        
        return {
          data: {
            user,
            session: {
              id: 'session',
              token: access_token,
              expiresAt: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString(),
            }
          }
        };
      } catch (error: unknown) {
        const axiosError = error as { response?: { data?: { details?: string; error?: string } } };
        const errorMessage = axiosError.response?.data?.details || axiosError.response?.data?.error || 'Login failed';
        setError(errorMessage);
        throw error;
      } finally {
        setLoading(false);
      }
    }
  },
  signUp: {
    email: async (credentials: { name: string; email: string; password: string }) => {
      const { setLoading, setError } = useAuthStore.getState();
      setLoading(true);
      setError(null);

      try {
        const response = await api.post('/auth/register', {
          email: credentials.email,
          name: credentials.name,
          password: credentials.password,
          confirm_password: credentials.password, // Backend expects confirm_password
        });
        
        return {
          data: {
            user: response.data.user
          }
        };
      } catch (error: unknown) {
        const axiosError = error as { response?: { data?: { details?: string; error?: string } } };
        const errorMessage = axiosError.response?.data?.details || axiosError.response?.data?.error || 'Registration failed';
        setError(errorMessage);
        throw error;
      } finally {
        setLoading(false);
      }
    }
  },
  signOut: async () => {
    const { session, clearAuth } = useAuthStore.getState();

    try {
      if (session?.access_token) {
        await api.post('/auth/logout', {}, {
          headers: {
            Authorization: `Bearer ${session.access_token}`
          }
        });
      }
    } catch {
      // Logout request failed
    } finally {
      clearAuth();
    }

    return {};
  },
  useSession: () => {
    const authState = useAuthStore();
    
    return {
      user: authState.user,
      session: authState.session ? {
        id: 'session',
        token: authState.session.access_token,
        expiresAt: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString(),
      } : null,
      isPending: authState.isPending,
      error: authState.error,
    };
  },
  getSession: async () => {
    const { user, session } = useAuthStore.getState();
    
    return {
      data: {
        user,
        session: session ? {
          id: 'session',
          token: session.access_token,
          expiresAt: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString(),
        } : null
      }
    };
  },
  refreshToken: async () => {
    const { session, setAuth, clearAuth } = useAuthStore.getState();

    if (!session?.refresh_token) {
      throw new Error('No refresh token available');
    }

    try {
      const response = await api.post('/auth/refresh', {
        refresh_token: session.refresh_token
      });

      const { access_token, refresh_token } = response.data;
      const { user } = useAuthStore.getState();
      
      if (user) {
        setAuth(user, { access_token, refresh_token });
      }

      return { access_token, refresh_token };
    } catch (error) {
      clearAuth();
      throw error;
    }
  }
};

export const {
  signIn,
  signUp,
  signOut,
  useSession,
} = authClient;