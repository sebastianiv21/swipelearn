import { QueryClient } from '@tanstack/react-query';

// Create a client with default config
export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      // Retry failed requests 1 time
      retry: 1,
      // Cache for 5 minutes
      staleTime: 1000 * 60 * 5,
      // Refetch on window focus
      refetchOnWindowFocus: true,
      // Don't refetch on reconnect by default
      refetchOnReconnect: false,
    },
    mutations: {
      // Retry failed mutations once
      retry: 1,
    },
  },
});