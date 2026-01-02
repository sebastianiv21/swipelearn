import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from '../lib/api';
import type { 
  Flashcard, 
  CreateFlashcardRequest, 
  UpdateFlashcardRequest 
} from '../types/api';

// Get all flashcards for the authenticated user
export function useFlashcards(deckId?: string) {
  return useQuery({
    queryKey: ['flashcards', deckId],
    queryFn: async (): Promise<Flashcard[]> => {
      const params = deckId ? { deck_id: deckId } : {};
      const response = await api.get('/flashcards', { params });
      return response.data.data; // Backend wraps flashcards in "data" object
    },
    enabled: true, // Can be called with or without deckId
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}



// Get flashcards that are due for review
export function useDueFlashcards(deckId?: string) {
  return useQuery({
    queryKey: ['flashcards', 'due', deckId],
    queryFn: async (): Promise<Flashcard[]> => {
      const params = deckId ? { deck_id: deckId } : {};
      const response = await api.get('/flashcards/due', { params });
      return response.data.data; // Backend wraps flashcards in "data" object
    },
    enabled: true, // Always enabled since we need due cards for study sessions
    staleTime: 60 * 1000, // 1 minute - more frequent updates for study sessions
  });
}

// Create a new flashcard
export function useCreateFlashcard() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async (data: CreateFlashcardRequest): Promise<Flashcard> => {
      const response = await api.post('/flashcards', data);
      return response.data; // Backend returns flashcard directly
    },
    onSuccess: (newFlashcard, { deck_id }) => {
      // Optimistically update the cache with the new flashcard
      queryClient.setQueryData(['flashcards', deck_id], (oldFlashcards: Flashcard[] | undefined) => {
        return oldFlashcards ? [...oldFlashcards, newFlashcard] : [newFlashcard];
      });
      
      // Invalidate other queries that might be affected
      queryClient.invalidateQueries({ queryKey: ['flashcards'] });
      queryClient.invalidateQueries({ queryKey: ['flashcards', 'due'] });
      queryClient.invalidateQueries({ queryKey: ['decks'] }); // For deck stats
    },
  });
}

// Update an existing flashcard
export function useUpdateFlashcard() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async ({ 
      id, 
      data 
    }: { 
      id: string; 
      data: UpdateFlashcardRequest; 
    }): Promise<Flashcard> => {
      const response = await api.put(`/flashcards/${id}`, data);
      return response.data;
    },
    onSuccess: (_, { id }) => {
      // Invalidate all flashcard-related queries since updates affect scheduling
      queryClient.invalidateQueries({ queryKey: ['flashcards'] });
      queryClient.invalidateQueries({ queryKey: ['flashcards', 'due'] });
      queryClient.invalidateQueries({ queryKey: ['flashcards', id] });
      queryClient.invalidateQueries({ queryKey: ['decks'] }); // For deck stats
    },
  });
}

// Delete a flashcard
export function useDeleteFlashcard() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async (id: string): Promise<void> => {
      await api.delete(`/flashcards/${id}`);
    },
    onSuccess: () => {
      // Invalidate all flashcard-related queries
      queryClient.invalidateQueries({ queryKey: ['flashcards'] });
      queryClient.invalidateQueries({ queryKey: ['flashcards', 'due'] });
      queryClient.invalidateQueries({ queryKey: ['decks'] }); // For deck stats
    },
  });
}