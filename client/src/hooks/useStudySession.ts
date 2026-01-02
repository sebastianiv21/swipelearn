import { useMutation, useQueryClient } from '@tanstack/react-query';
import api from '../lib/api';
import type { Flashcard, ReviewFlashcardRequest } from '../types/api';

// Review a flashcard using SM-2 algorithm
export function useReviewFlashcard() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async ({ 
      id, 
      quality 
    }: { 
      id: string; 
      quality: number; // 0-5 rating for SM-2 algorithm
    }): Promise<{
      flashcard: Flashcard;
      next_review: string;
      ease_factor: number;
      interval: number;
    }> => {
      const data: ReviewFlashcardRequest = { quality };
      const response = await api.post(`/flashcards/${id}/review`, data);
      return response.data;
    },
    onSuccess: () => {
      // Invalidate all flashcard-related queries since reviews affect scheduling
      queryClient.invalidateQueries({ queryKey: ['flashcards'] });
      queryClient.invalidateQueries({ queryKey: ['flashcards', 'due'] });
      queryClient.invalidateQueries({ queryKey: ['decks'] }); // For deck stats
    },
  });
}