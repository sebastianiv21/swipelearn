import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from '../lib/api';
import type { Deck, CreateDeckRequest, UpdateDeckRequest } from '../types/api';

// Get all decks for the authenticated user
export function useDecks() {
  return useQuery({
    queryKey: ['decks'],
    queryFn: async (): Promise<Deck[]> => {
      const response = await api.get('/decks');
      return response.data.data; // Backend wraps decks in "data" object
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

// Get a single deck by ID
export function useDeck(id: string) {
  return useQuery({
    queryKey: ['decks', id],
    queryFn: async (): Promise<Deck> => {
      const response = await api.get(`/decks/${id}`);
      return response.data;
    },
    enabled: !!id,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

// Create a new deck
export function useCreateDeck() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async (data: CreateDeckRequest): Promise<Deck> => {
      const response = await api.post('/decks', data);
      return response.data; // Backend returns deck directly
    },
    onSuccess: (newDeck) => {
      // Optimistically update the cache with the new deck
      queryClient.setQueryData(['decks'], (oldDecks: Deck[] | undefined) => {
        return oldDecks ? [...oldDecks, newDeck] : [newDeck];
      });
      // Also invalidate to ensure consistency
      queryClient.invalidateQueries({ queryKey: ['decks'] });
    },
  });
}

// Update an existing deck
export function useUpdateDeck() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async ({ id, data }: { id: string; data: UpdateDeckRequest }): Promise<Deck> => {
      const response = await api.put(`/decks/${id}`, data);
      return response.data;
    },
    onSuccess: (_, { id }) => {
      // Invalidate both the deck list and the specific deck
      queryClient.invalidateQueries({ queryKey: ['decks'] });
      queryClient.invalidateQueries({ queryKey: ['decks', id] });
    },
  });
}

// Delete a deck
export function useDeleteDeck() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async (id: string): Promise<void> => {
      await api.delete(`/decks/${id}`);
    },
    onSuccess: () => {
      // Invalidate the decks list
      queryClient.invalidateQueries({ queryKey: ['decks'] });
    },
  });
}